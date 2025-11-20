/**
 * 双棘轮算法实现
 * Double Ratchet Algorithm
 */

import nacl from 'tweetnacl';

/**
 * 执行 ECDH 密钥交换
 * @param {Uint8Array} privateKey 
 * @param {Uint8Array} publicKey 
 * @returns {Uint8Array}
 */
function ecdh(privateKey, publicKey) {
  return nacl.scalarMult(privateKey, publicKey);
}

/**
 * KDF - 密钥派生函数
 * @param {Uint8Array} inputKeyMaterial 
 * @param {string} info 
 * @param {number} outputLength 
 * @returns {Promise<Uint8Array>}
 */
async function kdf(inputKeyMaterial, info, outputLength = 32) {
  const encoder = new TextEncoder();
  const infoBytes = encoder.encode(info);

  const importedKey = await window.crypto.subtle.importKey(
    'raw',
    inputKeyMaterial,
    { name: 'HKDF' },
    false,
    ['deriveBits']
  );

  const derivedBits = await window.crypto.subtle.deriveBits(
    {
      name: 'HKDF',
      hash: 'SHA-256',
      salt: new Uint8Array(32),
      info: infoBytes,
    },
    importedKey,
    outputLength * 8
  );

  return new Uint8Array(derivedBits);
}

/**
 * 连接两个 Uint8Array
 * @param {Uint8Array} a 
 * @param {Uint8Array} b 
 * @returns {Uint8Array}
 */
function concat(a, b) {
  const result = new Uint8Array(a.length + b.length);
  result.set(a);
  result.set(b, a.length);
  return result;
}

/**
 * 从链密钥派生消息密钥
 * @param {Uint8Array} chainKey 
 * @returns {Promise<Uint8Array>}
 */
export async function deriveMessageKey(chainKey) {
  const messageKey = await kdf(chainKey, 'MessageKey', 32);
  // 注意：这里不记录完整日志，因为会频繁调用，只在调用方记录
  return messageKey;
}

/**
 * 更新链密钥
 * @param {Uint8Array} chainKey 
 * @returns {Promise<Uint8Array>}
 */
export async function advanceChainKey(chainKey) {
  const newChainKey = await kdf(chainKey, 'ChainKey', 32);
  // 注意：这里不记录完整日志，因为会频繁调用，只在调用方记录
  return newChainKey;
}

/**
 * 使用消息密钥加密消息
 * @param {Uint8Array} messageKey - 消息密钥
 * @param {string} plaintext - 明文消息
 * @returns {Promise<Object>} { ciphertext, iv, authTag }
 */
export async function encryptMessage(messageKey, plaintext) {
  // 生成随机 IV
  const iv = window.crypto.getRandomValues(new Uint8Array(12));

  // 导入消息密钥
  const cryptoKey = await window.crypto.subtle.importKey(
    'raw',
    messageKey,
    { name: 'AES-GCM' },
    false,
    ['encrypt']
  );

  // 将明文转换为 Uint8Array
  const encoder = new TextEncoder();
  const plaintextBytes = encoder.encode(plaintext);

  // 加密
  const ciphertext = await window.crypto.subtle.encrypt(
    { name: 'AES-GCM', iv: iv },
    cryptoKey,
    plaintextBytes
  );

  // AES-GCM 输出包含密文 + 认证标签（最后 16 字节）
  const ciphertextArray = new Uint8Array(ciphertext);
  const actualCiphertext = ciphertextArray.slice(0, -16);
  const authTag = ciphertextArray.slice(-16);

  return {
    ciphertext: actualCiphertext,
    iv: iv,
    authTag: authTag,
  };
}

/**
 * 使用消息密钥解密消息
 * @param {Uint8Array} messageKey - 消息密钥
 * @param {Uint8Array} ciphertext - 密文
 * @param {Uint8Array} iv - 初始化向量
 * @param {Uint8Array} authTag - 认证标签
 * @returns {Promise<string>} 明文消息
 */
export async function decryptMessage(messageKey, ciphertext, iv, authTag) {
  // 导入消息密钥
  const cryptoKey = await window.crypto.subtle.importKey(
    'raw',
    messageKey,
    { name: 'AES-GCM' },
    false,
    ['decrypt']
  );

  // 合并密文和认证标签
  const combined = new Uint8Array(ciphertext.length + authTag.length);
  combined.set(ciphertext);
  combined.set(authTag, ciphertext.length);

  // 解密
  const plaintext = await window.crypto.subtle.decrypt(
    { name: 'AES-GCM', iv: iv },
    cryptoKey,
    combined
  );

  // 转换为字符串
  const decoder = new TextDecoder();
  return decoder.decode(plaintext);
}

/**
 * DH 棘轮：更新 Root Key 和链密钥
 * @param {Uint8Array} rootKey - 当前根密钥
 * @param {Uint8Array} dhOutput - DH 输出
 * @param {boolean} isSender - 是否为发送方
 * @returns {Promise<Object>} { newRootKey, newChainKey }
 */
export async function dhRatchetStep(rootKey, dhOutput, isSender) {
  // 连接根密钥和 DH 输出
  const combined = concat(rootKey, dhOutput);

  // 派生新的根密钥
  const newRootKey = await kdf(combined, 'RootKey', 32);

  // 派生新的链密钥
  const chainKeyInfo = isSender ? 'SendingChainKey' : 'ReceivingChainKey';
  const newChainKey = await kdf(combined, chainKeyInfo, 32);

  return {
    newRootKey: newRootKey,
    newChainKey: newChainKey,
  };
}

/**
 * 生成新的 DH 棘轮密钥对
 * @returns {Object} { publicKey, privateKey }
 */
export function generateRatchetKeyPair() {
  const keyPair = nacl.box.keyPair();
  return {
    publicKey: keyPair.publicKey,
    privateKey: keyPair.secretKey,
  };
}

/**
 * 对称密钥棘轮：发送消息
 * @param {Uint8Array} chainKey - 当前链密钥
 * @param {string} plaintext - 明文消息
 * @returns {Promise<Object>} { encryptedMessage, newChainKey }
 */
export async function sendRatchet(chainKey, plaintext) {
  // 辅助函数：将 Uint8Array 转换为 Base64 预览
  const toBase64Preview = (arr) => {
    if (!arr) return null;
    const base64 = btoa(String.fromCharCode.apply(null, arr));
    return base64.substring(0, 20);
  };

  // 1. 派生消息密钥
  const messageKey = await deriveMessageKey(chainKey);

  // 2. 加密消息
  const encryptedMessage = await encryptMessage(messageKey, plaintext);

  // 3. 更新链密钥
  const newChainKey = await advanceChainKey(chainKey);

  console.log('🔐 sendRatchet 详情:', {
    chain_key_preview: toBase64Preview(chainKey),
    message_key_preview: toBase64Preview(messageKey),
    plaintext_length: plaintext.length,
    ciphertext_length: encryptedMessage.ciphertext.length,
    iv_preview: toBase64Preview(encryptedMessage.iv),
    auth_tag_preview: toBase64Preview(encryptedMessage.authTag),
    new_chain_key_preview: toBase64Preview(newChainKey),
  });

  // 4. 删除消息密钥（前向安全）
  // JavaScript 的垃圾回收会自动处理

  return {
    encryptedMessage: encryptedMessage,
    newChainKey: newChainKey,
  };
}

/**
 * 对称密钥棘轮：接收消息
 * @param {Uint8Array} chainKey - 当前链密钥
 * @param {Uint8Array} ciphertext - 密文
 * @param {Uint8Array} iv - 初始化向量
 * @param {Uint8Array} authTag - 认证标签
 * @returns {Promise<Object>} { plaintext, newChainKey }
 */
export async function receiveRatchet(chainKey, ciphertext, iv, authTag) {
  // 辅助函数：将 Uint8Array 转换为 Base64 预览
  const toBase64Preview = (arr) => {
    if (!arr) return null;
    const base64 = btoa(String.fromCharCode.apply(null, arr));
    return base64.substring(0, 20);
  };

  console.log('🔓 receiveRatchet 开始:', {
    chain_key_preview: toBase64Preview(chainKey),
    ciphertext_length: ciphertext?.length,
    iv_preview: toBase64Preview(iv),
    auth_tag_preview: toBase64Preview(authTag),
  });

  // 1. 派生消息密钥
  const messageKey = await deriveMessageKey(chainKey);
  console.log('🔑 receiveRatchet 派生消息密钥:', {
    message_key_preview: toBase64Preview(messageKey),
  });

  // 2. 解密消息
  let plaintext;
  try {
    plaintext = await decryptMessage(messageKey, ciphertext, iv, authTag);
    console.log('✅ receiveRatchet 解密成功:', {
      plaintext_length: plaintext?.length,
      plaintext_preview: plaintext?.substring(0, 20),
    });
  } catch (error) {
    console.error('❌ receiveRatchet 解密失败:', {
      error: error.message,
      error_name: error.name,
      chain_key_preview: toBase64Preview(chainKey),
      message_key_preview: toBase64Preview(messageKey),
      ciphertext_length: ciphertext?.length,
      iv_preview: toBase64Preview(iv),
      auth_tag_preview: toBase64Preview(authTag),
    });
    throw error;
  }

  // 3. 更新链密钥
  const newChainKey = await advanceChainKey(chainKey);
  console.log('🔄 receiveRatchet 更新链密钥:', {
    new_chain_key_preview: toBase64Preview(newChainKey),
  });

  // 4. 删除消息密钥
  // JavaScript 的垃圾回收会自动处理

  return {
    plaintext: plaintext,
    newChainKey: newChainKey,
  };
}

/**
 * 执行完整的 DH 棘轮更新（发送方）
 * @param {Uint8Array} rootKey - 当前根密钥
 * @param {Uint8Array} remoteRatchetPublicKey - 对方的棘轮公钥
 * @returns {Promise<Object>} { newRootKey, newSendingChainKey, newRatchetKeyPair }
 */
export async function performDHRatchetSend(rootKey, remoteRatchetPublicKey) {
  // 1. 生成新的 DH 密钥对
  const newRatchetKeyPair = generateRatchetKeyPair();

  // 2. 执行 ECDH
  const dhOutput = ecdh(newRatchetKeyPair.privateKey, remoteRatchetPublicKey);

  // 3. 更新根密钥和发送链密钥
  const result = await dhRatchetStep(rootKey, dhOutput, true);

  return {
    newRootKey: result.newRootKey,
    newSendingChainKey: result.newChainKey,
    newRatchetKeyPair: newRatchetKeyPair,
  };
}

/**
 * 执行完整的 DH 棘轮更新（接收方）
 * @param {Uint8Array} rootKey - 当前根密钥
 * @param {Uint8Array} localRatchetPrivateKey - 本地棘轮私钥
 * @param {Uint8Array} remoteRatchetPublicKey - 对方的新棘轮公钥
 * @returns {Promise<Object>} { newRootKey, newReceivingChainKey }
 */
export async function performDHRatchetReceive(
  rootKey,
  localRatchetPrivateKey,
  remoteRatchetPublicKey
) {
  // 1. 执行 ECDH
  const dhOutput = ecdh(localRatchetPrivateKey, remoteRatchetPublicKey);

  // 2. 更新根密钥和接收链密钥
  const result = await dhRatchetStep(rootKey, dhOutput, false);

  return {
    newRootKey: result.newRootKey,
    newReceivingChainKey: result.newChainKey,
  };
}

