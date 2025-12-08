/**
 * OSS 文件上传服务
 * 处理加密文件的上传和下载
 */

import axios from 'axios';
import apiAxios from '@/utils/axios';  // 项目的 axios 实例，带有 baseURL 和 token
import {
  generateFileKey,
  encryptFile,
  decryptFile,
  hashFile,
  bytesToBase64,
  base64ToBytes,
  generateThumbnail,
  encryptThumbnail,
  readFileAsArrayBuffer,
  isImageFile,
  formatFileSize
} from './fileEncryption';

/**
 * 获取 OSS 上传凭证
 * @param {string} fileType - 文件类型：'image' 或 'file'
 * @returns {Promise<Object>} 上传凭证
 */
export async function getUploadCredential(fileType = 'file') {
  const response = await apiAxios.get('/oss/upload-credential', {
    params: { file_type: fileType }
  });
  
  if (response.data.code !== 200) {
    throw new Error(response.data.message || '获取上传凭证失败');
  }
  
  return response.data.data;
}

/**
 * 获取 OSS 下载 URL
 * @param {string} ossKey - OSS 文件键
 * @returns {Promise<string>} 签名下载 URL
 */
export async function getDownloadURL(ossKey) {
  const response = await apiAxios.post('/oss/download-url', {
    oss_key: ossKey
  });
  
  if (response.data.code !== 200) {
    throw new Error(response.data.message || '获取下载 URL 失败');
  }
  
  return response.data.data.url;
}

/**
 * 上传加密文件到 OSS
 * @param {File} file - 原始文件
 * @param {Function} onProgress - 进度回调 (percent: number) => void
 * @returns {Promise<Object>} 包含 ossKey、fileKey、fileIv、fileHash 等信息
 */
export async function uploadEncryptedFile(file, onProgress = null) {
  // 1. 获取上传凭证
  const fileType = isImageFile(file) ? 'image' : 'file';
  const credential = await getUploadCredential(fileType);
  
  // 2. 读取文件
  if (onProgress) onProgress(5);
  const fileData = await readFileAsArrayBuffer(file);
  
  // 3. 计算原始文件哈希
  if (onProgress) onProgress(10);
  const fileHash = await hashFile(fileData);
  
  // 4. 生成加密密钥
  const { key: fileKey, iv: fileIv } = await generateFileKey();
  
  // 5. 加密文件
  if (onProgress) onProgress(20);
  const { encryptedData } = await encryptFile(fileData, fileKey, fileIv);
  
  // 6. 上传到 OSS（使用 STS 临时凭证和 V4 签名）
  // 注意：file 必须是最后一个表单域
  const formData = new FormData();
  formData.append('success_action_status', '200');
  formData.append('policy', credential.policy);
  formData.append('x-oss-signature', credential.signature);
  formData.append('x-oss-signature-version', credential.x_oss_signature_version);
  formData.append('x-oss-credential', credential.x_oss_credential);
  formData.append('x-oss-date', credential.x_oss_date);
  formData.append('x-oss-security-token', credential.security_token);
  formData.append('key', credential.key);
  formData.append('file', new Blob([encryptedData]), file.name + '.enc'); // file 必须放最后
  
  await axios.post(credential.host, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        // 20% 给加密，80% 给上传
        const uploadPercent = Math.round((progressEvent.loaded / progressEvent.total) * 80);
        onProgress(20 + uploadPercent);
      }
    }
  });
  
  if (onProgress) onProgress(100);
  
  return {
    ossKey: credential.key,
    fileKey: bytesToBase64(fileKey),
    fileIv: bytesToBase64(fileIv),
    fileHash: fileHash,
    fileName: file.name,
    fileSize: file.size,
    mimeType: file.type
  };
}

/**
 * 上传加密图片（包含缩略图）
 * @param {File} file - 图片文件
 * @param {Function} onProgress - 进度回调
 * @returns {Promise<Object>} 包含图片信息和缩略图
 */
export async function uploadEncryptedImage(file, onProgress = null) {
  // 1. 生成缩略图
  if (onProgress) onProgress(5);
  let thumbnailInfo = null;
  let imageWidth = 0;
  let imageHeight = 0;
  
  try {
    const thumbResult = await generateThumbnail(file, 200, 200, 0.6);
    imageWidth = thumbResult.originalWidth;
    imageHeight = thumbResult.originalHeight;
    
    // 加密缩略图
    if (onProgress) onProgress(10);
    const encryptedThumb = await encryptThumbnail(thumbResult.thumbnail);
    thumbnailInfo = {
      data: encryptedThumb.data,
      key: encryptedThumb.key,
      iv: encryptedThumb.iv,
      width: thumbResult.width,
      height: thumbResult.height
    };
  } catch (error) {
    console.warn('生成缩略图失败:', error);
  }
  
  // 2. 上传原图
  const uploadResult = await uploadEncryptedFile(file, (percent) => {
    // 缩略图占 10%，原图上传占 90%
    if (onProgress) onProgress(10 + Math.round(percent * 0.9));
  });
  
  return {
    ...uploadResult,
    width: imageWidth,
    height: imageHeight,
    thumbnail: thumbnailInfo
  };
}

/**
 * 下载并解密文件
 * @param {string} ossKey - OSS 文件键
 * @param {string} fileKeyBase64 - Base64 编码的文件密钥
 * @param {string} fileIvBase64 - Base64 编码的 IV
 * @param {string} expectedHash - 预期的文件哈希（可选，用于验证）
 * @param {Function} onProgress - 进度回调
 * @returns {Promise<ArrayBuffer>} 解密后的文件数据
 */
export async function downloadAndDecryptFile(
  ossKey,
  fileKeyBase64,
  fileIvBase64,
  expectedHash = null,
  onProgress = null
) {
  // 1. 获取下载 URL
  if (onProgress) onProgress(5);
  const downloadURL = await getDownloadURL(ossKey);
  
  // 2. 下载加密文件
  if (onProgress) onProgress(10);
  const response = await axios.get(downloadURL, {
    responseType: 'arraybuffer',
    onDownloadProgress: (progressEvent) => {
      if (onProgress && progressEvent.total) {
        const downloadPercent = Math.round((progressEvent.loaded / progressEvent.total) * 70);
        onProgress(10 + downloadPercent);
      }
    }
  });
  
  // 3. 解密文件
  if (onProgress) onProgress(85);
  const fileKey = base64ToBytes(fileKeyBase64);
  const fileIv = base64ToBytes(fileIvBase64);
  const decryptedData = await decryptFile(response.data, fileKey, fileIv);
  
  // 4. 验证哈希（如果提供）
  if (expectedHash) {
    if (onProgress) onProgress(95);
    const actualHash = await hashFile(decryptedData);
    if (actualHash !== expectedHash) {
      throw new Error('文件完整性验证失败：哈希不匹配');
    }
  }
  
  if (onProgress) onProgress(100);
  return decryptedData;
}

/**
 * 下载并解密文件，返回 Blob URL
 * @param {Object} fileInfo - 文件信息对象
 * @param {Function} onProgress - 进度回调
 * @returns {Promise<string>} Blob URL
 */
export async function downloadAndDecryptFileAsURL(fileInfo, onProgress = null) {
  const decryptedData = await downloadAndDecryptFile(
    fileInfo.ossKey,
    fileInfo.fileKey,
    fileInfo.fileIv,
    fileInfo.fileHash,
    onProgress
  );
  
  const blob = new Blob([decryptedData], { type: fileInfo.mimeType || 'application/octet-stream' });
  return URL.createObjectURL(blob);
}

/**
 * 下载并保存文件
 * @param {Object} fileInfo - 文件信息对象
 * @param {Function} onProgress - 进度回调
 */
export async function downloadAndSaveFile(fileInfo, onProgress = null) {
  const blobURL = await downloadAndDecryptFileAsURL(fileInfo, onProgress);
  
  // 创建下载链接
  const a = document.createElement('a');
  a.href = blobURL;
  a.download = fileInfo.fileName || 'download';
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  
  // 延迟释放 URL
  setTimeout(() => URL.revokeObjectURL(blobURL), 1000);
}

export { formatFileSize, isImageFile };
