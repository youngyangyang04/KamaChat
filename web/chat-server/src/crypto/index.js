/**
 * 加密模块统一导出
 */

// 密钥派生
export {
  deriveKey,
  generateSalt,
  encryptWithMasterKey,
  decryptWithMasterKey,
  verifyMasterKey,
} from './keyDerivation';

// 密钥生成
export {
  generateIdentityKeyPair,
  generateSignedPreKey,
  generateOneTimePreKeys,
  verifySignedPreKey,
  encryptPrivateKey,
  arrayBufferToBase64,
  base64ToArrayBuffer,
} from './keyGeneration';

// IndexedDB 存储
export { initCryptoStore, put, get, getAll, remove, clear, STORES } from './cryptoStore';

// 密钥管理器
export {
  initializeUserKeys,
  loginAndDeriveMasterKey,
  getIdentityPrivateKey,
  getSignedPrePrivateKey,
  getOneTimePrePrivateKey,
  countOneTimePreKeys,
  clearAllCryptoData,
} from './keyManager';

// X3DH 协议
export { aliceX3DH, bobX3DH, deriveInitialKeys } from './x3dh';

// 双棘轮算法
export {
  deriveMessageKey,
  advanceChainKey,
  encryptMessage,
  decryptMessage,
  sendRatchet,
  receiveRatchet,
  generateRatchetKeyPair,
  performDHRatchetSend,
  performDHRatchetReceive,
} from './doubleRatchet';

// 会话管理器
export {
  createSession,
  acceptSession,
  encryptAndSendMessage,
  receiveAndDecryptMessage,
  hasSession,
  deleteSession,
  updateRatchetKeyOnly,
} from './sessionManager';

// 文件加密
export {
  generateFileKey,
  encryptFile,
  decryptFile,
  hashFile,
  bytesToBase64,
  base64ToBytes,
  generateThumbnail,
  encryptThumbnail,
  decryptThumbnail,
  readFileAsArrayBuffer,
  formatFileSize,
  isImageFile,
  getFileExtension,
} from './fileEncryption';

// OSS 上传服务
export {
  getUploadCredential,
  getDownloadURL,
  uploadEncryptedFile,
  uploadEncryptedImage,
  downloadAndDecryptFile,
  downloadAndDecryptFileAsURL,
  downloadAndSaveFile,
} from './ossUpload';
