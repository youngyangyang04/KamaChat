/**
 * 密钥管理器
 * 统一管理用户密钥的生成、存储、加载
 */

import nacl from 'tweetnacl';
import {
  deriveKey,
  generateSalt,
  encryptWithMasterKey,
  decryptWithMasterKey,
  verifyMasterKey,
} from './keyDerivation';

import {
  generateIdentityKeyPair,
  generateSignedPreKey,
  generateOneTimePreKeys,
  arrayBufferToBase64,
  base64ToArrayBuffer,
} from './keyGeneration';

import { put, get, getAll, clear, STORES } from './cryptoStore';

/**
 * 初始化用户密钥（注册时调用）
 * @param {string} password - 用户密码
 * @returns {Promise<Object>} 包含公钥信息和主密钥相关数据
 */
export async function initializeUserKeys(password) {
  console.log('开始初始化用户密钥...');

  // 1. 生成 Salt 并派生主密钥
  const salt = generateSalt();
  const masterKey = await deriveKey(password, salt);

  // 2. 生成身份密钥对（Ed25519）
  const identityKeyPair = generateIdentityKeyPair();

  // 2.1. 从 Ed25519 身份私钥生成 Curve25519 密钥对（用于 ECDH）
  const identityCurve25519KeyPair = nacl.box.keyPair.fromSecretKey(identityKeyPair.privateKey.slice(0, 32));

  // 3. 生成签名预密钥对
  const signedPreKey = generateSignedPreKey(identityKeyPair.privateKey, 1);

  // 4. 生成 100 个一次性预密钥
  const oneTimePreKeys = generateOneTimePreKeys(100, 1);

  // 5. 加密并存储身份私钥
  const encryptedIdentityKey = await encryptWithMasterKey(
    masterKey,
    identityKeyPair.privateKey
  );
  await put(STORES.USER_KEYS, {
    key_type: 'identity',
    private_key: encryptedIdentityKey.ciphertext,
    public_key: identityKeyPair.publicKey,
    iv: encryptedIdentityKey.iv,
    auth_tag: encryptedIdentityKey.authTag,
    created_at: Date.now(),
  });

  // 6. 加密并存储签名预私钥
  const encryptedSignedPreKey = await encryptWithMasterKey(
    masterKey,
    signedPreKey.privateKey
  );
  await put(STORES.USER_KEYS, {
    key_type: 'signed_pre_key',
    key_id: signedPreKey.keyId,
    private_key: encryptedSignedPreKey.ciphertext,
    public_key: signedPreKey.publicKey,
    signature: signedPreKey.signature,
    iv: encryptedSignedPreKey.iv,
    auth_tag: encryptedSignedPreKey.authTag,
    created_at: Date.now(),
    expires_at: Date.now() + 30 * 24 * 60 * 60 * 1000, // 30 天后过期
  });

  // 7. 加密并存储一次性预私钥
  for (const key of oneTimePreKeys) {
    const encrypted = await encryptWithMasterKey(masterKey, key.privateKey);
    await put(STORES.USER_KEYS, {
      key_type: `one_time_pre_${key.keyId}`,
      key_id: key.keyId,
      private_key: encrypted.ciphertext,
      public_key: key.publicKey,
      iv: encrypted.iv,
      auth_tag: encrypted.authTag,
      created_at: Date.now(),
    });
  }

  // 8. 存储 Salt（用于登录时重新派生主密钥）
  // 重要：将 Uint8Array 转换为 base64 字符串，避免 IndexedDB 序列化问题
  await put(STORES.MASTER_KEY_SALT, {
    id: 1,
    salt: arrayBufferToBase64(salt),
    iterations: 100000,
    created_at: Date.now(),
  });

  // 9. 创建验证数据（用于验证主密钥是否正确）
  const testData = new Uint8Array([1, 2, 3, 4, 5]);
  const encryptedTestData = await encryptWithMasterKey(masterKey, testData);
  await put(STORES.MASTER_KEY_SALT, {
    id: 2,
    test_ciphertext: encryptedTestData.ciphertext,
    test_iv: encryptedTestData.iv,
    test_auth_tag: encryptedTestData.authTag,
  });

  console.log('用户密钥初始化完成');

  // 返回公钥数据（需要上传到服务器）
  return {
    identity_key_public: arrayBufferToBase64(identityKeyPair.publicKey), // Ed25519 公钥（用于验证签名）
    identity_key_public_curve25519: arrayBufferToBase64(identityCurve25519KeyPair.publicKey), // Curve25519 公钥（用于 ECDH）
    signed_pre_key: {
      key_id: signedPreKey.keyId,
      public_key: arrayBufferToBase64(signedPreKey.publicKey),
      signature: arrayBufferToBase64(signedPreKey.signature),
    },
    one_time_pre_keys: oneTimePreKeys.map((key) => ({
      key_id: key.keyId,
      public_key: arrayBufferToBase64(key.publicKey),
    })),
  };
}

/**
 * 登录时重新派生主密钥并验证
 * @param {string} password - 用户密码
 * @returns {Promise<Uint8Array|null>} 主密钥，如果密码错误则返回 null
 */
export async function loginAndDeriveMasterKey(password) {
  console.log('重新派生主密钥...');

  // 1. 读取 Salt
  const saltData = await get(STORES.MASTER_KEY_SALT, 1);
  if (!saltData) {
    throw new Error('未找到 Salt，请先注册');
  }

  // 2. 将 base64 字符串转换回 Uint8Array
  const salt = base64ToArrayBuffer(saltData.salt);

  // 3. 派生主密钥
  const masterKey = await deriveKey(password, salt, saltData.iterations);

  // 4. 验证主密钥
  const testData = await get(STORES.MASTER_KEY_SALT, 2);
  if (testData) {
    const isValid = await verifyMasterKey(masterKey, {
      ciphertext: testData.test_ciphertext,
      iv: testData.test_iv,
      authTag: testData.test_auth_tag,
    });

    if (!isValid) {
      console.error('密码错误');
      return null;
    }
  }

  console.log('主密钥验证成功');
  return masterKey;
}

/**
 * 解密并获取身份私钥
 * @param {Uint8Array} masterKey 
 * @returns {Promise<Uint8Array>}
 */
export async function getIdentityPrivateKey(masterKey) {
  console.log('🔑 [keyManager] 开始获取身份私钥，masterKey 长度:', masterKey?.length);
  
  const keyData = await get(STORES.USER_KEYS, 'identity');
  if (!keyData) {
    console.error('❌ [keyManager] 未找到身份密钥数据');
    throw new Error('未找到身份密钥');
  }
  
  console.log('🔑 [keyManager] 找到身份密钥数据，开始解密...');
  console.log('🔑 [keyManager] 密钥数据预览:', {
    has_private_key: !!keyData.private_key,
    private_key_length: keyData.private_key?.length,
    has_iv: !!keyData.iv,
    iv_length: keyData.iv?.length,
    has_auth_tag: !!keyData.auth_tag,
    auth_tag_length: keyData.auth_tag?.length,
  });

  try {
    const privateKey = await decryptWithMasterKey(
      masterKey,
      keyData.private_key,
      keyData.iv,
      keyData.auth_tag
    );
    console.log('✅ [keyManager] 身份私钥解密成功');
    return privateKey;
  } catch (error) {
    console.error('❌ [keyManager] 身份私钥解密失败:', error);
    console.error('❌ [keyManager] 这通常意味着主密钥不正确或密钥数据已损坏');
    throw error;
  }
}

/**
 * 解密并获取签名预私钥
 * @param {Uint8Array} masterKey 
 * @returns {Promise<Uint8Array>}
 */
export async function getSignedPrePrivateKey(masterKey) {
  const keyData = await get(STORES.USER_KEYS, 'signed_pre_key');
  if (!keyData) {
    throw new Error('未找到签名预密钥');
  }

  return await decryptWithMasterKey(
    masterKey,
    keyData.private_key,
    keyData.iv,
    keyData.auth_tag
  );
}

/**
 * 解密并获取一次性预私钥
 * @param {Uint8Array} masterKey 
 * @param {number} keyId 
 * @returns {Promise<Uint8Array>}
 */
export async function getOneTimePrePrivateKey(masterKey, keyId) {
  const keyData = await get(STORES.USER_KEYS, `one_time_pre_${keyId}`);
  if (!keyData) {
    throw new Error(`未找到一次性预密钥 ${keyId}`);
  }

  return await decryptWithMasterKey(
    masterKey,
    keyData.private_key,
    keyData.iv,
    keyData.auth_tag
  );
}

/**
 * 检查一次性预密钥数量
 * @returns {Promise<number>}
 */
export async function countOneTimePreKeys() {
  const allKeys = await getAll(STORES.USER_KEYS);
  const otpKeys = allKeys.filter((key) => key.key_type.startsWith('one_time_pre_'));
  return otpKeys.length;
}

/**
 * 清除所有加密数据（退出登录时调用）
 * @returns {Promise<void>}
 */
export async function clearAllCryptoData() {
  console.log('清除所有加密数据...');
  const stores = Object.values(STORES);
  for (const store of stores) {
    await clear(store);
  }
  console.log('加密数据已清除');
}

