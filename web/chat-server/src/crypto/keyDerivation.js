/**
 * 密钥派生模块
 * 使用 PBKDF2 从密码派生主密钥
 */

/**
 * 从密码派生主密钥
 * @param {string} password - 用户密码
 * @param {Uint8Array} salt - 盐值（32字节）
 * @param {number} iterations - 迭代次数，默认 100000
 * @returns {Promise<Uint8Array>} 主密钥（32字节）
 */
export async function deriveKey(password, salt, iterations = 100000) {
  // 将密码转换为 ArrayBuffer
  const encoder = new TextEncoder();
  const passwordBuffer = encoder.encode(password);

  // 导入密码为 CryptoKey
  const keyMaterial = await window.crypto.subtle.importKey(
    'raw',
    passwordBuffer,
    { name: 'PBKDF2' },
    false,
    ['deriveBits', 'deriveKey']
  );

  // 派生密钥
  const derivedKey = await window.crypto.subtle.deriveKey(
    {
      name: 'PBKDF2',
      salt: salt,
      iterations: iterations,
      hash: 'SHA-256',
    },
    keyMaterial,
    { name: 'AES-GCM', length: 256 },
    true,
    ['encrypt', 'decrypt']
  );

  // 导出为原始字节
  const exportedKey = await window.crypto.subtle.exportKey('raw', derivedKey);
  return new Uint8Array(exportedKey);
}

/**
 * 生成随机盐值
 * @param {number} length - 盐值长度，默认 32 字节
 * @returns {Uint8Array} 随机盐值
 */
export function generateSalt(length = 32) {
  return window.crypto.getRandomValues(new Uint8Array(length));
}

/**
 * 使用主密钥加密数据
 * @param {Uint8Array} masterKey - 主密钥
 * @param {Uint8Array} data - 要加密的数据
 * @returns {Promise<Object>} { ciphertext, iv, authTag }
 */
export async function encryptWithMasterKey(masterKey, data) {
  // 生成随机 IV
  const iv = window.crypto.getRandomValues(new Uint8Array(12));

  // 导入主密钥
  const cryptoKey = await window.crypto.subtle.importKey(
    'raw',
    masterKey,
    { name: 'AES-GCM' },
    false,
    ['encrypt']
  );

  // 加密
  const ciphertext = await window.crypto.subtle.encrypt(
    { name: 'AES-GCM', iv: iv },
    cryptoKey,
    data
  );

  // AES-GCM 的输出包含密文和认证标签（最后 16 字节）
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
 * 使用主密钥解密数据
 * @param {Uint8Array} masterKey - 主密钥
 * @param {Uint8Array} ciphertext - 密文
 * @param {Uint8Array} iv - 初始化向量
 * @param {Uint8Array} authTag - 认证标签
 * @returns {Promise<Uint8Array>} 解密后的数据
 */
export async function decryptWithMasterKey(masterKey, ciphertext, iv, authTag) {
  // 导入主密钥
  const cryptoKey = await window.crypto.subtle.importKey(
    'raw',
    masterKey,
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

  return new Uint8Array(plaintext);
}

/**
 * 验证主密钥是否正确
 * @param {Uint8Array} masterKey - 待验证的主密钥
 * @param {Object} testData - 测试数据 { ciphertext, iv, authTag }
 * @returns {Promise<boolean>} 是否正确
 */
export async function verifyMasterKey(masterKey, testData) {
  try {
    await decryptWithMasterKey(
      masterKey,
      testData.ciphertext,
      testData.iv,
      testData.authTag
    );
    return true;
  } catch (error) {
    console.error('主密钥验证失败:', error);
    return false;
  }
}

