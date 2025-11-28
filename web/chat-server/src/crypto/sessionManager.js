/**
 * 会话管理器
 * 管理加密会话的建立、消息加密/解密
 */

import nacl from 'tweetnacl';
import { put, get, remove, STORES } from './cryptoStore';
import { aliceX3DH, bobX3DH, deriveInitialKeys } from './x3dh';
import {
  sendRatchet,
  receiveRatchet,
  generateRatchetKeyPair,
  performDHRatchetSend,
  performDHRatchetReceive,
} from './doubleRatchet';
import { arrayBufferToBase64, base64ToArrayBuffer } from './keyGeneration';
import {
  getIdentityPrivateKey,
  getSignedPrePrivateKey,
  getOneTimePrePrivateKey,
} from './keyManager';

/**
 * 创建新会话（作为发起方 Alice）
 * @param {Uint8Array} masterKey - 主密钥
 * @param {string} contactId - 对方用户ID
 * @param {Object} remotePublicKeyBundle - 对方的公钥束
 * @returns {Promise<Object>} 会话初始化数据
 */
export async function createSession(masterKey, contactId, remotePublicKeyBundle) {
  console.log('创建新会话，contactId:', contactId);

  // 1. 获取本地密钥
  const identityPrivateKey = await getIdentityPrivateKey(masterKey);

  // 2. 生成临时密钥对
  const ephemeralKeyPair = generateRatchetKeyPair();

  // 3. 执行 X3DH 协商
  const x3dhResult = await aliceX3DH(
    {
      identityPrivateKey: identityPrivateKey,
      ephemeralPrivateKey: ephemeralKeyPair.privateKey,
    },
    remotePublicKeyBundle
  );

  // 4. 派生初始会话密钥（Alice 是发送方）
  const initialKeys = await deriveInitialKeys(x3dhResult.sharedSecret, true);
  
  console.log('🔑 Alice 派生初始密钥:', {
    sharedSecret_length: x3dhResult.sharedSecret.length,
    sharedSecret_preview: arrayBufferToBase64(x3dhResult.sharedSecret).substring(0, 20),
    rootKey_length: initialKeys.rootKey.length,
    rootKey_preview: arrayBufferToBase64(initialKeys.rootKey).substring(0, 20),
    sendingChainKey_length: initialKeys.sendingChainKey?.length || 0,
    sendingChainKey_preview: initialKeys.sendingChainKey ? arrayBufferToBase64(initialKeys.sendingChainKey).substring(0, 20) : null,
    receivingChainKey_length: initialKeys.receivingChainKey?.length || 0,
    receivingChainKey_preview: initialKeys.receivingChainKey ? arrayBufferToBase64(initialKeys.receivingChainKey).substring(0, 20) : null,
  });
  
  // 添加验证提示
  console.log('🔍 Alice 密钥验证提示:', {
    alice_sendingChainKey_preview: initialKeys.sendingChainKey ? arrayBufferToBase64(initialKeys.sendingChainKey).substring(0, 20) : null,
    note: 'Alice 的 sendingChainKey 应该等于 Bob 的 receivingChainKey',
  });

  // 5. 生成初始 DH 棘轮密钥对
  const initialRatchetKeyPair = generateRatchetKeyPair();

  // 6. 构造会话状态
  const sessionState = {
    contact_id: contactId,
    root_key: initialKeys.rootKey,
    sending_chain_key: initialKeys.sendingChainKey,  // Alice 使用这个发送
    sending_ratchet_key_private: initialRatchetKeyPair.privateKey,
    sending_ratchet_key_public: initialRatchetKeyPair.publicKey,
    receiving_chain_key: initialKeys.receivingChainKey || null,  // 等待 Bob 发送第一条消息后设置
    receiving_ratchet_key_public: null, // 等待对方发送
    send_counter: 0,
    receive_counter: 0,
    prev_counter: 0,
    remote_identity_key: base64ToArrayBuffer(remotePublicKeyBundle.identity_key),
    created_at: Date.now(),
    updated_at: Date.now(),
    last_message_at: Date.now(),
  };

  // 7. 保存会话状态
  await put(STORES.SESSIONS, sessionState);

  console.log('会话创建成功');

  // 从 Ed25519 身份私钥生成 Curve25519 密钥对（用于 ECDH）
  const identityCurve25519KeyPair = nacl.box.keyPair.fromSecretKey(identityPrivateKey.slice(0, 32));
  
  // 返回初始化数据（需要发送给对方）
  const initData = {
    identity_key: arrayBufferToBase64(identityPrivateKey.slice(32, 64)), // Ed25519 公钥（用于验证签名）
    identity_key_curve25519: arrayBufferToBase64(identityCurve25519KeyPair.publicKey), // Curve25519 公钥（用于 ECDH）
    ephemeral_key: arrayBufferToBase64(x3dhResult.ephemeralPublicKey),
    ratchet_key: arrayBufferToBase64(initialRatchetKeyPair.publicKey),
    used_one_time_pre_key_id: x3dhResult.usedOneTimePreKeyId, // 可能是 null 或数字
  };
  
  console.log('📤 Alice 返回初始化数据:', {
    has_identity_key: !!initData.identity_key,
    has_identity_key_curve25519: !!initData.identity_key_curve25519,
    has_ephemeral_key: !!initData.ephemeral_key,
    has_ratchet_key: !!initData.ratchet_key,
    used_one_time_pre_key_id: initData.used_one_time_pre_key_id,
    used_one_time_pre_key_id_type: typeof initData.used_one_time_pre_key_id,
  });
  
  return initData;
}

/**
 * 接受会话（作为接收方 Bob）
 * @param {Uint8Array} masterKey - 主密钥
 * @param {string} contactId - 对方用户ID
 * @param {Object} aliceInitData - Alice 的初始化数据
 * @returns {Promise<void>}
 */
export async function acceptSession(masterKey, contactId, aliceInitData) {
  console.log('接受会话，contactId:', contactId);
  console.log('📥 Bob 收到的 Alice 初始化数据:', {
    has_identity_key: !!aliceInitData.identity_key,
    has_identity_key_curve25519: !!aliceInitData.identity_key_curve25519,
    has_ephemeral_key: !!aliceInitData.ephemeral_key,
    has_ratchet_key: !!aliceInitData.ratchet_key,
    used_one_time_pre_key_id: aliceInitData.used_one_time_pre_key_id,
    identity_key_preview: aliceInitData.identity_key?.substring(0, 20),
    ephemeral_key_preview: aliceInitData.ephemeral_key?.substring(0, 20),
    ratchet_key_preview: aliceInitData.ratchet_key?.substring(0, 20),
  });

  // 1. 获取本地密钥
  const identityPrivateKey = await getIdentityPrivateKey(masterKey);
  const signedPreKeyPrivate = await getSignedPrePrivateKey(masterKey);

  let oneTimePreKeyPrivate = null;
  if (aliceInitData.used_one_time_pre_key_id) {
    console.log('🔑 Bob 尝试获取一次性预密钥，key_id:', aliceInitData.used_one_time_pre_key_id);
    oneTimePreKeyPrivate = await getOneTimePrePrivateKey(
      masterKey,
      aliceInitData.used_one_time_pre_key_id
    );
    if (!oneTimePreKeyPrivate) {
      console.warn('⚠️ Bob 无法获取一次性预密钥，key_id:', aliceInitData.used_one_time_pre_key_id);
      throw new Error(`无法获取一次性预密钥 key_id: ${aliceInitData.used_one_time_pre_key_id}，可能已被使用或不存在`);
    } else {
      console.log('✅ Bob 成功获取一次性预密钥');
    }
  } else {
    console.log('ℹ️ Alice 未使用一次性预密钥');
  }

  // 2. 解析 Alice 的公钥
  // 如果提供了 Curve25519 格式的身份公钥，使用它；否则使用 Ed25519 公钥（向后兼容）
  const aliceIdentityKeyCurve25519 = aliceInitData.identity_key_curve25519
    ? base64ToArrayBuffer(aliceInitData.identity_key_curve25519)
    : null;
  
  const alicePublicKeys = {
    identityKey: base64ToArrayBuffer(aliceInitData.identity_key), // Ed25519 公钥（用于验证签名）
    identityKeyCurve25519: aliceIdentityKeyCurve25519, // Curve25519 公钥（用于 ECDH）
    ephemeralKey: base64ToArrayBuffer(aliceInitData.ephemeral_key),
  };

  // 3. 执行 X3DH 协商
  const sharedSecret = await bobX3DH(
    {
      identityPrivateKey: identityPrivateKey,
      signedPreKeyPrivate: signedPreKeyPrivate,
      oneTimePreKeyPrivate: oneTimePreKeyPrivate,
    },
    alicePublicKeys
  );

  // 4. 派生初始会话密钥（Bob 是接收方）
  const initialKeys = await deriveInitialKeys(sharedSecret, false);
  
  console.log('🔑 Bob 派生初始密钥:', {
    sharedSecret_length: sharedSecret.length,
    sharedSecret_preview: arrayBufferToBase64(sharedSecret).substring(0, 20),
    rootKey_length: initialKeys.rootKey.length,
    rootKey_preview: arrayBufferToBase64(initialKeys.rootKey).substring(0, 20),
    sendingChainKey_length: initialKeys.sendingChainKey?.length || 0,
    sendingChainKey_preview: initialKeys.sendingChainKey ? arrayBufferToBase64(initialKeys.sendingChainKey).substring(0, 20) : null,
    receivingChainKey_length: initialKeys.receivingChainKey?.length || 0,
    receivingChainKey_preview: initialKeys.receivingChainKey ? arrayBufferToBase64(initialKeys.receivingChainKey).substring(0, 20) : null,
  });
  
  // 验证：记录密钥信息以便与 Alice 的日志对比
  console.log('🔍 Bob 密钥验证提示:', {
    bob_receivingChainKey_preview: initialKeys.receivingChainKey ? arrayBufferToBase64(initialKeys.receivingChainKey).substring(0, 20) : null,
    bob_sendingChainKey_preview: initialKeys.sendingChainKey ? arrayBufferToBase64(initialKeys.sendingChainKey).substring(0, 20) : null,
    bob_sharedSecret_preview: arrayBufferToBase64(sharedSecret).substring(0, 20),
    note: 'Alice 的 sendingChainKey 应该等于 Bob 的 receivingChainKey',
    note2: '如果链密钥不匹配，说明 sharedSecret 不一致，需要检查 X3DH 协商过程',
  });

  // 5. 生成初始 DH 棘轮密钥对
  const initialRatchetKeyPair = generateRatchetKeyPair();

  // 6. 构造会话状态
  const receivingRatchetKeyPublic = base64ToArrayBuffer(aliceInitData.ratchet_key);
  console.log('🔑 Bob 设置 receiving_ratchet_key_public:', {
    ratchet_key_base64_preview: aliceInitData.ratchet_key?.substring(0, 20),
    ratchet_key_bytes_length: receivingRatchetKeyPublic.length,
    receiving_chain_key_preview: arrayBufferToBase64(initialKeys.receivingChainKey).substring(0, 20),
    sending_chain_key_preview: arrayBufferToBase64(initialKeys.sendingChainKey).substring(0, 20),
  });
  
  const sessionState = {
    contact_id: contactId,
    root_key: initialKeys.rootKey,
    sending_chain_key: initialKeys.sendingChainKey || null,  // Bob 还没有发送消息，所以 Bob 还没有 sendingChainKey
    sending_ratchet_key_private: initialRatchetKeyPair.privateKey,
    sending_ratchet_key_public: initialRatchetKeyPair.publicKey,
    receiving_chain_key: initialKeys.receivingChainKey,  // Bob 使用这个接收 Alice 的消息
    receiving_ratchet_key_public: receivingRatchetKeyPublic,
    send_counter: 0,
    receive_counter: 0,
    prev_counter: 0,
    remote_identity_key: alicePublicKeys.identityKey,
    created_at: Date.now(),
    updated_at: Date.now(),
    last_message_at: Date.now(),
  };

  // 7. 保存会话状态
  await put(STORES.SESSIONS, sessionState);
  
  // 验证保存是否成功
  const savedSession = await get(STORES.SESSIONS, contactId);
  console.log('✅ Bob 会话保存验证:', {
    has_receiving_ratchet_key_public: !!savedSession?.receiving_ratchet_key_public,
    receiving_ratchet_key_public_type: savedSession?.receiving_ratchet_key_public?.constructor?.name,
    receiving_ratchet_key_public_length: savedSession?.receiving_ratchet_key_public?.length,
    has_receiving_chain_key: !!savedSession?.receiving_chain_key,
    receiving_chain_key_type: savedSession?.receiving_chain_key?.constructor?.name,
    receiving_chain_key_length: savedSession?.receiving_chain_key?.length,
  });

  console.log('会话接受成功');
}

/**
 * 加密并发送消息
 * @param {string} contactId - 对方用户ID
 * @param {string} plaintext - 明文消息
 * @returns {Promise<Object>} 加密消息数据
 */
export async function encryptAndSendMessage(contactId, plaintext) {
  console.log(`🔐 [sessionManager] 开始加密消息，contactId: ${contactId}`);
  
  // 1. 读取会话状态
  const session = await get(STORES.SESSIONS, contactId);
  if (!session) {
    console.error(`❌ [sessionManager] 会话不存在: ${contactId}`);
    throw new Error('会话不存在，请先建立会话');
  }
  
  console.log(`🔐 [sessionManager] 读取到会话状态:`, {
    contactId: contactId,
    has_root_key: !!session.root_key,
    has_sending_chain_key: !!session.sending_chain_key,
    has_receiving_chain_key: !!session.receiving_chain_key,
    send_counter: session.send_counter,
    receive_counter: session.receive_counter,
    has_sending_ratchet_key: !!session.sending_ratchet_key_private,
    has_receiving_ratchet_key: !!session.receiving_ratchet_key_public,
  });

  // 确保 Uint8Array 字段是正确的类型（IndexedDB 可能将其序列化为普通对象）
  if (session.root_key && !(session.root_key instanceof Uint8Array)) {
    session.root_key = new Uint8Array(session.root_key);
  }
  if (session.sending_chain_key && !(session.sending_chain_key instanceof Uint8Array)) {
    session.sending_chain_key = new Uint8Array(session.sending_chain_key);
  }
  if (session.receiving_chain_key && !(session.receiving_chain_key instanceof Uint8Array)) {
    session.receiving_chain_key = new Uint8Array(session.receiving_chain_key);
  }
  if (session.sending_ratchet_key_private && !(session.sending_ratchet_key_private instanceof Uint8Array)) {
    session.sending_ratchet_key_private = new Uint8Array(session.sending_ratchet_key_private);
  }
  if (session.sending_ratchet_key_public && !(session.sending_ratchet_key_public instanceof Uint8Array)) {
    session.sending_ratchet_key_public = new Uint8Array(session.sending_ratchet_key_public);
  }
  if (session.receiving_ratchet_key_public && !(session.receiving_ratchet_key_public instanceof Uint8Array)) {
    session.receiving_ratchet_key_public = new Uint8Array(session.receiving_ratchet_key_public);
  }
  if (session.remote_identity_key && !(session.remote_identity_key instanceof Uint8Array)) {
    session.remote_identity_key = new Uint8Array(session.remote_identity_key);
  }

  // 2. 检查是否需要 DH 棘轮更新
  // 2.1. 如果 sending_chain_key 为 null（Bob 第一次发送消息），需要先执行 DH ratchet
  if (!session.sending_chain_key && session.receiving_ratchet_key_public) {
    console.log('🔑 首次发送消息，执行 DH ratchet 生成 sending_chain_key...');
    console.log('🔑 DH ratchet 参数检查:', {
      root_key_preview: arrayBufferToBase64(session.root_key).substring(0, 20),
      receiving_ratchet_key_public_preview: arrayBufferToBase64(session.receiving_ratchet_key_public).substring(0, 20),
      receiving_ratchet_key_public_length: session.receiving_ratchet_key_public?.length,
    });
    
    const ratchetResult = await performDHRatchetSend(
      session.root_key,
      session.receiving_ratchet_key_public
    );

    session.root_key = ratchetResult.newRootKey;
    session.sending_chain_key = ratchetResult.newSendingChainKey;
    session.sending_ratchet_key_private = ratchetResult.newRatchetKeyPair.privateKey;
    session.sending_ratchet_key_public = ratchetResult.newRatchetKeyPair.publicKey;
    
    console.log('✅ 首次 DH ratchet 完成:', {
      new_root_key_preview: arrayBufferToBase64(ratchetResult.newRootKey).substring(0, 20),
      new_sending_chain_key_preview: arrayBufferToBase64(ratchetResult.newSendingChainKey).substring(0, 20),
      new_sending_ratchet_key_public_preview: arrayBufferToBase64(ratchetResult.newRatchetKeyPair.publicKey).substring(0, 20),
    });
  }
  // 2.2. 定期 DH 棘轮更新（每 100 条消息）
  else if (session.send_counter % 100 === 0 && session.send_counter > 0) {
    console.log('执行定期 DH 棘轮更新...');
    const ratchetResult = await performDHRatchetSend(
      session.root_key,
      session.receiving_ratchet_key_public
    );

    session.root_key = ratchetResult.newRootKey;
    session.sending_chain_key = ratchetResult.newSendingChainKey;
    session.sending_ratchet_key_private = ratchetResult.newRatchetKeyPair.privateKey;
    session.sending_ratchet_key_public = ratchetResult.newRatchetKeyPair.publicKey;
  }

  // 3. 检查 sending_chain_key 是否存在
  if (!session.sending_chain_key) {
    throw new Error('发送链密钥不存在，无法加密消息。可能需要先接收对方的消息以建立接收链密钥。');
  }

  // 4. 对称密钥棘轮：加密消息
  console.log('🔐 加密消息:', {
    plaintext_length: plaintext.length,
    sending_chain_key_length: session.sending_chain_key?.length,
    sending_chain_key_preview: session.sending_chain_key ? arrayBufferToBase64(session.sending_chain_key).substring(0, 20) : null,
    send_counter: session.send_counter,
  });
  
  const result = await sendRatchet(session.sending_chain_key, plaintext);
  
  console.log('✅ 加密完成:', {
    ciphertext_bytes_length: result.encryptedMessage.ciphertext.length,
    iv_bytes_length: result.encryptedMessage.iv.length,
    auth_tag_bytes_length: result.encryptedMessage.authTag.length,
    ciphertext_base64: arrayBufferToBase64(result.encryptedMessage.ciphertext).substring(0, 20),
  });

  // 4. 更新会话状态
  const oldChainKey = session.sending_chain_key;
  session.sending_chain_key = result.newChainKey;
  session.send_counter += 1;
  session.updated_at = Date.now();
  session.last_message_at = Date.now();

  console.log('🔄 Alice 更新发送链密钥:', {
    old_chain_key_preview: arrayBufferToBase64(oldChainKey).substring(0, 20),
    new_chain_key_preview: arrayBufferToBase64(result.newChainKey).substring(0, 20),
    send_counter: session.send_counter,
  });

  // 5. 保存会话状态
  await put(STORES.SESSIONS, session);

  // 6. 构造加密消息
  const encryptedData = {
    ratchet_key: arrayBufferToBase64(session.sending_ratchet_key_public),
    counter: session.send_counter - 1, // 当前消息的计数器
    prev_counter: session.receive_counter,
    ciphertext: arrayBufferToBase64(result.encryptedMessage.ciphertext),
    iv: arrayBufferToBase64(result.encryptedMessage.iv),
    auth_tag: arrayBufferToBase64(result.encryptedMessage.authTag),
  };
  
  console.log('📤 Alice 发送数据:', {
    ratchet_key_preview: encryptedData.ratchet_key.substring(0, 20),
    counter: encryptedData.counter,
    prev_counter: encryptedData.prev_counter,
    ciphertext_base64_length: encryptedData.ciphertext.length,
    ciphertext_base64_preview: encryptedData.ciphertext.substring(0, 20),
    iv_preview: encryptedData.iv.substring(0, 20),
    auth_tag_preview: encryptedData.auth_tag.substring(0, 20),
  });
  
  return encryptedData;
}

/**
 * 接收并解密消息
 * @param {string} contactId - 对方用户ID
 * @param {Object} encryptedMessage - 加密消息数据
 * @returns {Promise<string>} 明文消息
 */
export async function receiveAndDecryptMessage(contactId, encryptedMessage) {
  // 1. 读取会话状态
  const session = await get(STORES.SESSIONS, contactId);
  if (!session) {
    throw new Error('会话不存在');
  }

  // 确保 Uint8Array 字段是正确的类型（IndexedDB 可能将其序列化为普通对象）
  if (session.root_key && !(session.root_key instanceof Uint8Array)) {
    session.root_key = new Uint8Array(session.root_key);
  }
  if (session.sending_chain_key && !(session.sending_chain_key instanceof Uint8Array)) {
    session.sending_chain_key = new Uint8Array(session.sending_chain_key);
  }
  if (session.receiving_chain_key && !(session.receiving_chain_key instanceof Uint8Array)) {
    session.receiving_chain_key = new Uint8Array(session.receiving_chain_key);
  }
  if (session.sending_ratchet_key_private && !(session.sending_ratchet_key_private instanceof Uint8Array)) {
    session.sending_ratchet_key_private = new Uint8Array(session.sending_ratchet_key_private);
  }
  if (session.sending_ratchet_key_public && !(session.sending_ratchet_key_public instanceof Uint8Array)) {
    session.sending_ratchet_key_public = new Uint8Array(session.sending_ratchet_key_public);
  }
  if (session.receiving_ratchet_key_public && !(session.receiving_ratchet_key_public instanceof Uint8Array)) {
    session.receiving_ratchet_key_public = new Uint8Array(session.receiving_ratchet_key_public);
  }
  if (session.remote_identity_key && !(session.remote_identity_key instanceof Uint8Array)) {
    session.remote_identity_key = new Uint8Array(session.remote_identity_key);
  }

  console.log('🔓 开始解密消息:', {
    contactId,
    counter: encryptedMessage.counter,
    receive_counter: session.receive_counter,
    has_receiving_ratchet_key: !!session.receiving_ratchet_key_public,
  });

  // 2. 检查是否需要 DH 棘轮更新
  const remoteRatchetKey = base64ToArrayBuffer(encryptedMessage.ratchet_key);
  const hasReceivingRatchetKey = !!session.receiving_ratchet_key_public;
  const isNewRatchetKey = hasReceivingRatchetKey && !arraysEqual(session.receiving_ratchet_key_public, remoteRatchetKey);

  console.log('🔑 棘轮密钥检查:', {
    isNewRatchetKey,
    has_receiving_ratchet_key: hasReceivingRatchetKey,
    receiving_chain_key_length: session.receiving_chain_key?.length,
    receiving_chain_key_exists: !!session.receiving_chain_key,
    message_type: encryptedMessage.message_type,
    is_prekey_message: encryptedMessage.message_type === 'PreKeyMessage',
    receive_counter: session.receive_counter,
  });

  // 3. 特殊处理：PreKeyMessage 且 receive_counter === 0 时
  // 这表示这是第一条消息，应该使用初始的 receivingChainKey 解密，不需要 DH 棘轮更新
  // 注意：即使 receiving_ratchet_key_public 已经存在（在 acceptSession 中设置），
  // 对于第一条 PreKeyMessage，也应该使用初始链密钥，不执行 DH 棘轮更新
  if (encryptedMessage.message_type === 'PreKeyMessage' && session.receive_counter === 0) {
    console.log('📌 PreKeyMessage 且 receive_counter === 0，使用初始链密钥解密（不执行 DH 棘轮更新）');
    // 检查 receiving_chain_key 是否存在
    if (!session.receiving_chain_key) {
      throw new Error('接收链密钥不存在：PreKeyMessage 需要初始接收链密钥');
    }
    
    // 添加调试日志：输出链密钥信息以便对比
    console.log('🔍 PreKeyMessage 链密钥调试:', {
      receiving_chain_key_preview: arrayBufferToBase64(session.receiving_chain_key).substring(0, 20),
      receiving_chain_key_length: session.receiving_chain_key.length,
      root_key_preview: arrayBufferToBase64(session.root_key).substring(0, 20),
      note: '这个 receiving_chain_key 应该等于 Alice 的 sending_chain_key',
    });
    
    // 如果 receiving_ratchet_key_public 还没有设置，则设置它（但不执行 DH 棘轮更新）
    if (!hasReceivingRatchetKey) {
      session.receiving_ratchet_key_public = remoteRatchetKey;
      console.log('📌 设置 receiving_ratchet_key_public（但不执行 DH 棘轮更新）');
    } else {
      // 如果已经存在，验证它是否匹配
      if (!arraysEqual(session.receiving_ratchet_key_public, remoteRatchetKey)) {
        console.warn('⚠️ PreKeyMessage 的 ratchet_key 与已存储的不匹配，但使用初始链密钥解密');
        console.warn('⚠️ 已存储的 ratchet_key:', arrayBufferToBase64(session.receiving_ratchet_key_public).substring(0, 20));
        console.warn('⚠️ 消息中的 ratchet_key:', arrayBufferToBase64(remoteRatchetKey).substring(0, 20));
        // 更新为新的 ratchet_key，但不执行 DH 棘轮更新
        session.receiving_ratchet_key_public = remoteRatchetKey;
      }
    }
  } 
  // 4. SignalMessage 且 receiving_ratchet_key_public 为 null 时（Bob 第一次发送消息给 Alice）
  // 这表示对方第一次发送消息，需要执行 DH ratchet 更新来生成 receiving_chain_key
  else if (encryptedMessage.message_type === 'SignalMessage' && !hasReceivingRatchetKey && session.sending_ratchet_key_private) {
    console.log('📌 SignalMessage 且 receiving_ratchet_key_public 为 null，执行 DH ratchet 更新...');
    console.log('📌 DH ratchet 参数检查:', {
      root_key_preview: arrayBufferToBase64(session.root_key).substring(0, 20),
      sending_ratchet_key_private_exists: !!session.sending_ratchet_key_private,
      sending_ratchet_key_private_length: session.sending_ratchet_key_private?.length,
      remote_ratchet_key_preview: arrayBufferToBase64(remoteRatchetKey).substring(0, 20),
      remote_ratchet_key_length: remoteRatchetKey.length,
    });
    
    // 执行 DH ratchet 更新
    const ratchetResult = await performDHRatchetReceive(
      session.root_key,
      session.sending_ratchet_key_private,
      remoteRatchetKey
    );

    session.root_key = ratchetResult.newRootKey;
    session.receiving_chain_key = ratchetResult.newReceivingChainKey;
    session.receiving_ratchet_key_public = remoteRatchetKey;
    session.prev_counter = session.receive_counter;
    session.receive_counter = 0;
    session.updated_at = Date.now();
    
    // 立即保存会话状态（在执行 DH ratchet 更新后）
    await put(STORES.SESSIONS, session);
    
    console.log('✅ 首次 DH ratchet 更新完成:', {
      new_root_key_preview: arrayBufferToBase64(ratchetResult.newRootKey).substring(0, 20),
      new_receiving_chain_key_preview: arrayBufferToBase64(ratchetResult.newReceivingChainKey).substring(0, 20),
      receive_counter_reset: true,
    });
  }
  // 5. 检测到新的棘轮密钥，执行 DH 棘轮更新（仅对 SignalMessage）
  else if (isNewRatchetKey && encryptedMessage.message_type === 'SignalMessage') {
    console.log('🔄 检测到新的棘轮密钥，执行 DH 棘轮更新...');
    const oldRootKey = session.root_key;
    const oldReceivingChainKey = session.receiving_chain_key;
    
    const ratchetResult = await performDHRatchetReceive(
      session.root_key,
      session.sending_ratchet_key_private,
      remoteRatchetKey
    );

    session.root_key = ratchetResult.newRootKey;
    session.receiving_chain_key = ratchetResult.newReceivingChainKey;
    session.receiving_ratchet_key_public = remoteRatchetKey;
    session.prev_counter = session.receive_counter;
    session.receive_counter = 0;
    
    console.log('✅ DH 棘轮更新完成:', {
      old_root_key_preview: arrayBufferToBase64(oldRootKey).substring(0, 20),
      new_root_key_preview: arrayBufferToBase64(ratchetResult.newRootKey).substring(0, 20),
      old_receiving_chain_key_preview: arrayBufferToBase64(oldReceivingChainKey).substring(0, 20),
      new_receiving_chain_key_preview: arrayBufferToBase64(ratchetResult.newReceivingChainKey).substring(0, 20),
      receive_counter_reset: true,
    });
  }

  // 5. 检查 receiving_chain_key 是否存在（在所有 DH ratchet 更新之后）
  if (!session.receiving_chain_key) {
    throw new Error('接收链密钥不存在，无法解密消息。可能需要先接收对方的消息以建立接收链密钥。');
  }

  // 6. 处理乱序消息
  // 如果 counter 与 receive_counter 不匹配，需要调整链密钥
  // 注意：这里简化处理，假设 counter >= receive_counter（未来的消息）
  // 如果 counter < receive_counter（过去的消息），则需要使用缓存的消息密钥
  const counterDiff = encryptedMessage.counter - session.receive_counter;
  let chainKeyToUse = session.receiving_chain_key;
  
  if (counterDiff !== 0) {
    console.warn('检测到乱序消息，counter:', encryptedMessage.counter, 'receive_counter:', session.receive_counter, '差值:', counterDiff);
    
    if (counterDiff < 0) {
      // counter < receive_counter，这是过去的消息，应该已经处理过
      // 简化处理：仍然尝试解密，但可能会失败
      console.warn('⚠️ 收到过去的消息（counter < receive_counter），尝试解密可能会失败');
    } else if (counterDiff > 0) {
      // counter > receive_counter，这是未来的消息，说明之前可能有消息丢失
      // 由于 Double Ratchet 的特性，如果中间的消息丢失，我们无法解密未来的消息
      // 因为链密钥已经基于之前的消息推进了，我们无法回到之前的状态
      // 但是，为了兼容性，我们先尝试解密，如果失败再抛出错误
      console.warn('⚠️ 收到未来的消息（counter > receive_counter），可能之前有消息丢失，尝试解密可能会失败');
    }
  }

  // 7. 对称密钥棘轮：解密消息
  console.log('🔐 使用接收链密钥解密:', {
    chain_key_length: session.receiving_chain_key?.length,
    chain_key_preview: session.receiving_chain_key ? arrayBufferToBase64(session.receiving_chain_key).substring(0, 20) : null,
    ciphertext_base64_length: encryptedMessage.ciphertext?.length,
    iv_base64_length: encryptedMessage.iv?.length,
    auth_tag_base64_length: encryptedMessage.auth_tag?.length,
    ciphertext_base64_preview: encryptedMessage.ciphertext?.substring(0, 20),
  });
  
  // 解码 Base64
  let ciphertextBytes, ivBytes, authTagBytes;
  try {
    ciphertextBytes = base64ToArrayBuffer(encryptedMessage.ciphertext);
    ivBytes = base64ToArrayBuffer(encryptedMessage.iv);
    authTagBytes = base64ToArrayBuffer(encryptedMessage.auth_tag);
    console.log('✅ Base64 解码成功:', {
      ciphertext_bytes_length: ciphertextBytes.length,
      iv_bytes_length: ivBytes.length,
      auth_tag_bytes_length: authTagBytes.length,
    });
  } catch (error) {
    console.error('❌ Base64 解码失败:', error);
    throw new Error('Base64 解码失败: ' + error.message);
  }
  
  const oldReceivingChainKey = session.receiving_chain_key;
  
  try {
    const result = await receiveRatchet(
      session.receiving_chain_key,
      ciphertextBytes,
      ivBytes,
      authTagBytes
    );
    console.log('✅ receiveAndDecryptMessage 解密成功:', {
      plaintext_preview: result.plaintext?.substring(0, 20),
      plaintext_length: result.plaintext?.length,
    });

    // 5. 更新会话状态
    session.receiving_chain_key = result.newChainKey;
    session.receive_counter += 1;
    session.updated_at = Date.now();
    session.last_message_at = Date.now();

    console.log('🔄 receiveAndDecryptMessage 更新接收链密钥:', {
      old_chain_key_preview: arrayBufferToBase64(oldReceivingChainKey).substring(0, 20),
      new_chain_key_preview: arrayBufferToBase64(result.newChainKey).substring(0, 20),
      receive_counter: session.receive_counter,
    });

    // 6. 保存会话状态
    await put(STORES.SESSIONS, session);

    return result.plaintext;
  } catch (error) {
    console.error('❌ receiveAndDecryptMessage 解密失败:', {
      error: error.message,
      error_name: error.name,
      error_stack: error.stack?.substring(0, 200),
      receiving_chain_key_length: session.receiving_chain_key?.length,
      receiving_chain_key_preview: session.receiving_chain_key ? arrayBufferToBase64(session.receiving_chain_key).substring(0, 20) : null,
      ciphertext_length: encryptedMessage.ciphertext?.length,
      iv_length: encryptedMessage.iv?.length,
      auth_tag_length: encryptedMessage.auth_tag?.length,
      counter: encryptedMessage.counter,
      receive_counter: session.receive_counter,
    });
    throw error;
  }
}

/**
 * 检查会话是否存在
 * @param {string} contactId 
 * @returns {Promise<boolean>}
 */
export async function hasSession(contactId) {
  const session = await get(STORES.SESSIONS, contactId);
  return session !== null;
}

/**
 * 删除会话
 * @param {string} contactId 
 * @returns {Promise<void>}
 */
export async function deleteSession(contactId) {
  await remove(STORES.SESSIONS, contactId);
  console.log('会话已删除，contactId:', contactId);
}

/**
 * 只更新 ratchet_key（不重建会话）
 * 仅当会话存在但 ratchet_key 缺失，且其他状态完整时使用
 * @param {string} contactId - 联系人ID
 * @param {string} ratchetKey - Base64 编码的 ratchet_key
 * @returns {Promise<void>}
 */
export async function updateRatchetKeyOnly(contactId, ratchetKey) {
  const session = await get(STORES.SESSIONS, contactId);
  if (!session) {
    throw new Error('会话不存在，无法只更新 ratchet_key');
  }
  
  // 检查是否可以只更新 ratchet_key
  // 需要满足：root_key 和 receiving_chain_key 存在，且 receive_counter === 0
  const canUpdateOnly = 
    session.root_key && 
    session.receiving_chain_key && 
    session.receive_counter === 0 &&
    !session.receiving_ratchet_key_public; // 只有缺失时才更新
  
  if (!canUpdateOnly) {
    throw new Error('会话状态不完整，无法只更新 ratchet_key，需要重建会话');
  }
  
  // 确保 Uint8Array 字段是正确的类型
  if (session.root_key && !(session.root_key instanceof Uint8Array)) {
    session.root_key = new Uint8Array(session.root_key);
  }
  if (session.receiving_chain_key && !(session.receiving_chain_key instanceof Uint8Array)) {
    session.receiving_chain_key = new Uint8Array(session.receiving_chain_key);
  }
  
  // 只更新 ratchet_key
  session.receiving_ratchet_key_public = base64ToArrayBuffer(ratchetKey);
  session.updated_at = Date.now();
  
  await put(STORES.SESSIONS, session);
  console.log('✅ 已更新 ratchet_key（未重建会话）:', {
    contactId,
    ratchet_key_preview: ratchetKey.substring(0, 20),
    root_key_exists: !!session.root_key,
    receiving_chain_key_exists: !!session.receiving_chain_key,
  });
}

/**
 * 比较两个 Uint8Array 是否相等
 * @param {Uint8Array} a 
 * @param {Uint8Array} b 
 * @returns {boolean}
 */
function arraysEqual(a, b) {
  if (a.length !== b.length) return false;
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false;
  }
  return true;
}

