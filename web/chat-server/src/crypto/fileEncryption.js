/**
 * 文件加密模块
 * 使用 AES-256-GCM 加密文件，密钥通过 Signal Protocol 安全传输
 */

/**
 * 生成随机文件加密密钥
 * @returns {Promise<{key: Uint8Array, iv: Uint8Array}>}
 */
export async function generateFileKey() {
  // 生成 256 位 AES 密钥
  const key = new Uint8Array(32);
  crypto.getRandomValues(key);
  
  // 生成 96 位 IV（AES-GCM 推荐长度）
  const iv = new Uint8Array(12);
  crypto.getRandomValues(iv);
  
  return { key, iv };
}

/**
 * 加密文件
 * @param {ArrayBuffer|Uint8Array} fileData - 原始文件数据
 * @param {Uint8Array} key - 256 位 AES 密钥
 * @param {Uint8Array} iv - 96 位初始化向量
 * @returns {Promise<{encryptedData: ArrayBuffer, authTag: Uint8Array}>}
 */
export async function encryptFile(fileData, key, iv) {
  // 导入密钥
  const cryptoKey = await crypto.subtle.importKey(
    'raw',
    key,
    { name: 'AES-GCM' },
    false,
    ['encrypt']
  );
  
  // 加密
  const encryptedData = await crypto.subtle.encrypt(
    {
      name: 'AES-GCM',
      iv: iv,
      tagLength: 128 // 认证标签长度（位）
    },
    cryptoKey,
    fileData
  );
  
  // AES-GCM 的输出包含密文 + 认证标签
  // 认证标签在最后 16 字节
  const encryptedArray = new Uint8Array(encryptedData);
  
  return {
    encryptedData: encryptedData,
    // 注意：Web Crypto API 的 AES-GCM 输出已包含 auth tag
    // 不需要单独提取，解密时会自动验证
  };
}

/**
 * 解密文件
 * @param {ArrayBuffer|Uint8Array} encryptedData - 加密的文件数据
 * @param {Uint8Array} key - 256 位 AES 密钥
 * @param {Uint8Array} iv - 96 位初始化向量
 * @returns {Promise<ArrayBuffer>} 解密后的原始文件数据
 */
export async function decryptFile(encryptedData, key, iv) {
  // 导入密钥
  const cryptoKey = await crypto.subtle.importKey(
    'raw',
    key,
    { name: 'AES-GCM' },
    false,
    ['decrypt']
  );
  
  // 解密（自动验证认证标签）
  const decryptedData = await crypto.subtle.decrypt(
    {
      name: 'AES-GCM',
      iv: iv,
      tagLength: 128
    },
    cryptoKey,
    encryptedData
  );
  
  return decryptedData;
}

/**
 * 计算文件的 SHA-256 哈希
 * @param {ArrayBuffer|Uint8Array} data - 文件数据
 * @returns {Promise<string>} 十六进制哈希字符串
 */
export async function hashFile(data) {
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  const hashArray = new Uint8Array(hashBuffer);
  return Array.from(hashArray)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');
}

/**
 * 将 Uint8Array 转换为 Base64 字符串
 * @param {Uint8Array} bytes 
 * @returns {string}
 */
export function bytesToBase64(bytes) {
  let binary = '';
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}

/**
 * 将 Base64 字符串转换为 Uint8Array
 * @param {string} base64 
 * @returns {Uint8Array}
 */
export function base64ToBytes(base64) {
  const binary = atob(base64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
}

/**
 * 生成图片缩略图
 * @param {File} file - 图片文件
 * @param {number} maxWidth - 最大宽度
 * @param {number} maxHeight - 最大高度
 * @param {number} quality - JPEG 质量 (0-1)
 * @returns {Promise<{thumbnail: Blob, width: number, height: number, originalWidth: number, originalHeight: number}>}
 */
export async function generateThumbnail(file, maxWidth = 200, maxHeight = 200, quality = 0.7) {
  return new Promise((resolve, reject) => {
    const img = new Image();
    const url = URL.createObjectURL(file);
    
    img.onload = () => {
      URL.revokeObjectURL(url);
      
      const originalWidth = img.width;
      const originalHeight = img.height;
      
      // 计算缩放比例
      let width = originalWidth;
      let height = originalHeight;
      
      if (width > maxWidth) {
        height = (height * maxWidth) / width;
        width = maxWidth;
      }
      if (height > maxHeight) {
        width = (width * maxHeight) / height;
        height = maxHeight;
      }
      
      // 创建 canvas 绘制缩略图
      const canvas = document.createElement('canvas');
      canvas.width = width;
      canvas.height = height;
      
      const ctx = canvas.getContext('2d');
      ctx.drawImage(img, 0, 0, width, height);
      
      // 转换为 Blob
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve({
              thumbnail: blob,
              width: Math.round(width),
              height: Math.round(height),
              originalWidth,
              originalHeight
            });
          } else {
            reject(new Error('生成缩略图失败'));
          }
        },
        'image/jpeg',
        quality
      );
    };
    
    img.onerror = () => {
      URL.revokeObjectURL(url);
      reject(new Error('加载图片失败'));
    };
    
    img.src = url;
  });
}

/**
 * 加密缩略图并转为 Base64
 * @param {Blob} thumbnailBlob - 缩略图 Blob
 * @returns {Promise<{data: string, key: string, iv: string}>}
 */
export async function encryptThumbnail(thumbnailBlob) {
  const { key, iv } = await generateFileKey();
  const arrayBuffer = await thumbnailBlob.arrayBuffer();
  const { encryptedData } = await encryptFile(arrayBuffer, key, iv);
  
  return {
    data: bytesToBase64(new Uint8Array(encryptedData)),
    key: bytesToBase64(key),
    iv: bytesToBase64(iv)
  };
}

/**
 * 解密缩略图
 * @param {string} encryptedData - Base64 编码的加密数据
 * @param {string} keyBase64 - Base64 编码的密钥
 * @param {string} ivBase64 - Base64 编码的 IV
 * @returns {Promise<Blob>}
 */
export async function decryptThumbnail(encryptedData, keyBase64, ivBase64) {
  const key = base64ToBytes(keyBase64);
  const iv = base64ToBytes(ivBase64);
  const encrypted = base64ToBytes(encryptedData);
  
  const decrypted = await decryptFile(encrypted, key, iv);
  return new Blob([decrypted], { type: 'image/jpeg' });
}

/**
 * 读取文件为 ArrayBuffer
 * @param {File} file 
 * @returns {Promise<ArrayBuffer>}
 */
export function readFileAsArrayBuffer(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result);
    reader.onerror = () => reject(reader.error);
    reader.readAsArrayBuffer(file);
  });
}

/**
 * 获取友好的文件大小字符串
 * @param {number} bytes 
 * @returns {string}
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * 检查是否为图片文件
 * @param {File} file 
 * @returns {boolean}
 */
export function isImageFile(file) {
  return file.type.startsWith('image/');
}

/**
 * 获取文件扩展名
 * @param {string} filename 
 * @returns {string}
 */
export function getFileExtension(filename) {
  const parts = filename.split('.');
  return parts.length > 1 ? parts.pop().toLowerCase() : '';
}
