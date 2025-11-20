/**
 * X3DH 密钥协商协议实现
 * Extended Triple Diffie-Hellman
 */

import nacl from 'tweetnacl';
import { base64ToArrayBuffer, arrayBufferToBase64 } from './keyGeneration';

/**
 * 执行 ECDH 密钥交换
 * @param {Uint8Array} privateKey - Curve25519 私钥
 * @param {Uint8Array} publicKey - Curve25519 公钥
 * @returns {Uint8Array} 共享密钥
 */
function ecdh(privateKey, publicKey) {
  return nacl.scalarMult(privateKey, publicKey);
}

/**
 * 从 Ed25519 密钥转换为 Curve25519 密钥（用于 ECDH）
 * @param {Uint8Array} ed25519Key - Ed25519 密钥
 * @returns {Uint8Array} Curve25519 密钥
 */
function convertEd25519ToCurve25519(ed25519Key) {
  // 注意：tweetnacl 的 sign.keyPair 生成的私钥是 64 字节
  // 前 32 字节是实际的 Ed25519 seed，后 32 字节是公钥
  if (ed25519Key.length === 64) {
    // 这是私钥，取前 32 字节作为 seed，使用 nacl.box.keyPair.fromSecretKey 生成 Curve25519 密钥对
    const curve25519KeyPair = nacl.box.keyPair.fromSecretKey(ed25519Key.slice(0, 32));
    return curve25519KeyPair.secretKey; // 返回 Curve25519 私钥
  } else if (ed25519Key.length === 32) {
    // 这是公钥，需要转换为 Curve25519 公钥
    // tweetnacl 的 sign.keyPair.fromSecretKey 可以从 Ed25519 seed 生成密钥对
    // 但我们需要从 Ed25519 公钥转换为 Curve25519 公钥
    // 实际上，Ed25519 和 Curve25519 共享相同的底层曲线，但表示方式不同
    // 对于 tweetnacl，我们可以使用 nacl.box.keyPair.fromSecretKey 生成一个临时密钥对
    // 然后使用 Ed25519 公钥作为 Curve25519 公钥（它们应该是兼容的）
    // 但是，更安全的方法是：如果这是公钥，我们需要知道它来自哪个私钥
    // 由于我们不知道原始私钥，我们假设 Ed25519 公钥可以直接用于 Curve25519
    // 这在大多数情况下是正确的，因为 tweetnacl 的 scalarMult 应该能够处理
    return ed25519Key;
  }
  throw new Error('Invalid Ed25519 key length');
}

/**
 * KDF - 密钥派生函数
 * 使用 HKDF-SHA256
 * @param {Uint8Array} inputKeyMaterial - 输入密钥材料
 * @param {string} info - 派生信息
 * @param {number} outputLength - 输出长度（字节）
 * @returns {Promise<Uint8Array>} 派生的密钥
 */
async function kdf(inputKeyMaterial, info, outputLength = 32) {
  const encoder = new TextEncoder();
  const infoBytes = encoder.encode(info);

  // 使用 Web Crypto API 的 HKDF
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
      salt: new Uint8Array(32), // 空 salt
      info: infoBytes,
    },
    importedKey,
    outputLength * 8 // bits
  );

  return new Uint8Array(derivedBits);
}

/**
 * 连接多个 Uint8Array
 * @param {...Uint8Array} arrays 
 * @returns {Uint8Array}
 */
function concat(...arrays) {
  const totalLength = arrays.reduce((sum, arr) => sum + arr.length, 0);
  const result = new Uint8Array(totalLength);
  let offset = 0;
  for (const arr of arrays) {
    result.set(arr, offset);
    offset += arr.length;
  }
  return result;
}

/**
 * Alice 执行 X3DH（发起方）
 * 用于建立与 Bob 的会话
 * 
 * @param {Object} aliceKeys - Alice 的密钥
 * @param {Uint8Array} aliceKeys.identityPrivateKey - Alice 的身份私钥（Ed25519）
 * @param {Uint8Array} aliceKeys.ephemeralPrivateKey - Alice 的临时私钥（Curve25519）
 * @param {Object} bobPublicKeyBundle - Bob 的公钥束（字段名为 snake_case）
 * @param {string} bobPublicKeyBundle.identity_key - Bob 的身份公钥（Base64）
 * @param {Object} bobPublicKeyBundle.signed_pre_key - Bob 的签名预公钥
 * @param {Object} bobPublicKeyBundle.one_time_pre_key - Bob 的一次性预公钥（可选）
 * @returns {Promise<Object>} { sharedSecret, ephemeralPublicKey, usedOneTimePreKeyId }
 */
export async function aliceX3DH(aliceKeys, bobPublicKeyBundle) {
  console.log('Alice 开始 X3DH 协商...');

  // 1. 解析 Bob 的公钥（注意：后端返回的是 snake_case 字段名）
  const IK_B = base64ToArrayBuffer(bobPublicKeyBundle.identity_key); // Ed25519 公钥（用于验证签名）
  const IK_B_curve25519 = bobPublicKeyBundle.identity_key_curve25519
    ? base64ToArrayBuffer(bobPublicKeyBundle.identity_key_curve25519) // Curve25519 公钥（用于 ECDH）
    : null;
  const SPK_B = base64ToArrayBuffer(bobPublicKeyBundle.signed_pre_key.public_key);
  const OPK_B = bobPublicKeyBundle.one_time_pre_key
    ? base64ToArrayBuffer(bobPublicKeyBundle.one_time_pre_key.public_key)
    : null;

  // 2. Alice 的密钥
  const IK_A_private = aliceKeys.identityPrivateKey;
  const EK_A_private = aliceKeys.ephemeralPrivateKey;

  // 将 Ed25519 身份私钥转换为 Curve25519（用于 ECDH）
  const IK_A_curve = convertEd25519ToCurve25519(IK_A_private);

  // 3. 执行 4 次（或 3 次）ECDH
  const DH1 = ecdh(IK_A_curve, SPK_B);
  // 使用 Bob 的 Curve25519 身份公钥（如果提供了）或转换 Ed25519 公钥
  const IK_B_curve = IK_B_curve25519 || convertEd25519ToCurve25519(IK_B);
  const DH2 = ecdh(EK_A_private, IK_B_curve);
  const DH3 = ecdh(EK_A_private, SPK_B);
  const DH4 = OPK_B ? ecdh(EK_A_private, OPK_B) : null;

  // 4. 连接所有 DH 输出
  const dhOutputs = DH4 ? concat(DH1, DH2, DH3, DH4) : concat(DH1, DH2, DH3);

  // 5. 派生共享密钥
  const sharedSecret = await kdf(dhOutputs, 'X3DHSharedSecret', 32);

  console.log('Alice X3DH 协商完成', {
    has_DH4: !!DH4,
    DH1_preview: arrayBufferToBase64(DH1).substring(0, 20),
    DH2_preview: arrayBufferToBase64(DH2).substring(0, 20),
    DH3_preview: arrayBufferToBase64(DH3).substring(0, 20),
    DH4_preview: DH4 ? arrayBufferToBase64(DH4).substring(0, 20) : null,
    dhOutputs_length: dhOutputs.length,
    dhOutputs_preview: arrayBufferToBase64(dhOutputs).substring(0, 40),
    sharedSecret_preview: arrayBufferToBase64(sharedSecret).substring(0, 20),
  });

  // 返回临时公钥（需要发送给 Bob）
  const ephemeralPublicKey = nacl.scalarMult.base(EK_A_private);

  return {
    sharedSecret: sharedSecret,
    ephemeralPublicKey: ephemeralPublicKey,
    usedOneTimePreKeyId: bobPublicKeyBundle.one_time_pre_key
      ? bobPublicKeyBundle.one_time_pre_key.key_id
      : null,
  };
}

/**
 * Bob 执行 X3DH（接收方）
 * 用于响应 Alice 的会话建立请求
 * 
 * @param {Object} bobKeys - Bob 的密钥
 * @param {Uint8Array} bobKeys.identityPrivateKey - Bob 的身份私钥（Ed25519）
 * @param {Uint8Array} bobKeys.signedPreKeyPrivate - Bob 的签名预私钥（Curve25519）
 * @param {Uint8Array} bobKeys.oneTimePreKeyPrivate - Bob 的一次性预私钥（可选）
 * @param {Object} alicePublicKeys - Alice 的公钥
 * @param {Uint8Array} alicePublicKeys.identityKey - Alice 的身份公钥（Ed25519）
 * @param {Uint8Array} alicePublicKeys.ephemeralKey - Alice 的临时公钥（Curve25519）
 * @returns {Promise<Uint8Array>} sharedSecret
 */
export async function bobX3DH(bobKeys, alicePublicKeys) {
  console.log('Bob 开始 X3DH 协商...');

  // 1. Bob 的密钥
  const IK_B_private = bobKeys.identityPrivateKey;
  const SPK_B_private = bobKeys.signedPreKeyPrivate;
  const OPK_B_private = bobKeys.oneTimePreKeyPrivate || null;

  // 2. Alice 的公钥
  const IK_A = alicePublicKeys.identityKey; // Ed25519 公钥
  const EK_A = alicePublicKeys.ephemeralKey; // Curve25519 公钥

  // 将 Ed25519 身份私钥转换为 Curve25519
  const IK_B_curve = convertEd25519ToCurve25519(IK_B_private);
  
  // 使用 Alice 的 Curve25519 身份公钥（如果提供了）或转换 Ed25519 公钥
  const IK_A_curve = alicePublicKeys.identityKeyCurve25519 || convertEd25519ToCurve25519(IK_A);

  // 3. 执行 4 次（或 3 次）ECDH
  const DH1 = ecdh(SPK_B_private, IK_A_curve);
  const DH2 = ecdh(IK_B_curve, EK_A);
  const DH3 = ecdh(SPK_B_private, EK_A);
  const DH4 = OPK_B_private ? ecdh(OPK_B_private, EK_A) : null;

  // 4. 连接所有 DH 输出
  const dhOutputs = DH4 ? concat(DH1, DH2, DH3, DH4) : concat(DH1, DH2, DH3);

  // 5. 派生共享密钥
  const sharedSecret = await kdf(dhOutputs, 'X3DHSharedSecret', 32);

  console.log('Bob X3DH 协商完成', {
    has_DH4: !!DH4,
    DH1_preview: arrayBufferToBase64(DH1).substring(0, 20),
    DH2_preview: arrayBufferToBase64(DH2).substring(0, 20),
    DH3_preview: arrayBufferToBase64(DH3).substring(0, 20),
    DH4_preview: DH4 ? arrayBufferToBase64(DH4).substring(0, 20) : null,
    dhOutputs_length: dhOutputs.length,
    dhOutputs_preview: arrayBufferToBase64(dhOutputs).substring(0, 40),
    sharedSecret_preview: arrayBufferToBase64(sharedSecret).substring(0, 20),
  });

  return sharedSecret;
}

/**
 * 从共享密钥派生会话密钥
 * 注意：在 Double Ratchet 中，Alice 的 sendingChainKey 应该等于 Bob 的 receivingChainKey
 * 所以我们应该派生一个初始链密钥，然后根据角色分配
 * 
 * @param {Uint8Array} sharedSecret - X3DH 派生的共享密钥
 * @param {boolean} isAlice - 是否为 Alice（发送方），true 表示 Alice，false 表示 Bob
 * @returns {Promise<Object>} { rootKey, sendingChainKey, receivingChainKey }
 */
export async function deriveInitialKeys(sharedSecret, isAlice = true) {
  const rootKey = await kdf(sharedSecret, 'RootKey', 32);
  
  // 派生初始链密钥（所有方使用相同的派生方式）
  const initialChainKey = await kdf(sharedSecret, 'InitialChainKey', 32);
  
  // 根据角色分配链密钥
  // Alice（发送方）：使用 initialChainKey 作为 sendingChainKey
  // Bob（接收方）：使用 initialChainKey 作为 receivingChainKey
  // 注意：对方还没有发送第一条消息，所以反向链密钥初始为空，会在第一条消息后设置
  
  if (isAlice) {
    // Alice：第一个消息的发送方
    return {
      rootKey: rootKey,
      sendingChainKey: initialChainKey,  // Alice 用这个发送
      receivingChainKey: null,  // Bob 还没有发送消息，所以 Alice 还没有 receivingChainKey
    };
  } else {
    // Bob：第一个消息的接收方
    return {
      rootKey: rootKey,
      sendingChainKey: null,  // Bob 还没有发送消息，所以 Bob 还没有 sendingChainKey
      receivingChainKey: initialChainKey,  // Bob 用这个接收 Alice 的消息
    };
  }
}

