/**
 * 密钥生成模块
 * 生成身份密钥、签名预密钥、一次性预密钥
 */

import nacl from 'tweetnacl';
import { encryptWithMasterKey } from './keyDerivation';

/**
 * 生成身份密钥对（Ed25519）
 * @returns {Object} { publicKey, privateKey }
 */
export function generateIdentityKeyPair() {
  const keyPair = nacl.sign.keyPair();
  return {
    publicKey: keyPair.publicKey,  // Uint8Array(32)
    privateKey: keyPair.secretKey,  // Uint8Array(64)
  };
}

/**
 * 生成签名预密钥对（Curve25519）
 * @param {Uint8Array} identityPrivateKey - 身份私钥（用于签名）
 * @param {number} keyId - 密钥ID
 * @returns {Object} { keyId, publicKey, privateKey, signature }
 */
export function generateSignedPreKey(identityPrivateKey, keyId = 1) {
  // 生成 Curve25519 密钥对
  const keyPair = nacl.box.keyPair();

  // 用身份私钥签名公钥
  const signature = nacl.sign.detached(keyPair.publicKey, identityPrivateKey);

  return {
    keyId: keyId,
    publicKey: keyPair.publicKey,  // Uint8Array(32)
    privateKey: keyPair.secretKey,  // Uint8Array(32)
    signature: signature,           // Uint8Array(64)
  };
}

/**
 * 批量生成一次性预密钥（Curve25519）
 * @param {number} count - 生成数量，默认 100
 * @param {number} startId - 起始ID，默认 1
 * @returns {Array} [{ keyId, publicKey, privateKey }, ...]
 */
export function generateOneTimePreKeys(count = 100, startId = 1) {
  const keys = [];
  for (let i = 0; i < count; i++) {
    const keyPair = nacl.box.keyPair();
    keys.push({
      keyId: startId + i,
      publicKey: keyPair.publicKey,  // Uint8Array(32)
      privateKey: keyPair.secretKey,  // Uint8Array(32)
    });
  }
  return keys;
}

/**
 * 验证签名预密钥的签名
 * @param {Uint8Array} publicKey - 签名预公钥
 * @param {Uint8Array} signature - 签名
 * @param {Uint8Array} identityPublicKey - 身份公钥
 * @returns {boolean} 签名是否有效
 */
export function verifySignedPreKey(publicKey, signature, identityPublicKey) {
  return nacl.sign.detached.verify(publicKey, signature, identityPublicKey);
}

/**
 * 加密并存储私钥
 * @param {Uint8Array} masterKey - 主密钥
 * @param {Uint8Array} privateKey - 私钥
 * @returns {Promise<Object>} { ciphertext, iv, authTag }
 */
export async function encryptPrivateKey(masterKey, privateKey) {
  return await encryptWithMasterKey(masterKey, privateKey);
}

/**
 * 将 Uint8Array 转换为 Base64
 * @param {Uint8Array} buffer 
 * @returns {string}
 */
export function arrayBufferToBase64(buffer) {
  const bytes = new Uint8Array(buffer);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return window.btoa(binary);
}

/**
 * 将 Base64 转换为 Uint8Array
 * @param {string} base64 
 * @returns {Uint8Array}
 */
export function base64ToArrayBuffer(base64) {
  const binary = window.atob(base64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
}

