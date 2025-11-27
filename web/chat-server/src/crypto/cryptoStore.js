/**
 * IndexedDB 加密存储模块
 * 存储用户密钥、会话状态、消息密钥缓存
 */

const DB_NAME_PREFIX = 'wework_crypto';
const DB_VERSION = 3; // 增加版本号以添加新的 Object Store

// Object Store 名称
const STORES = {
  USER_KEYS: 'user_keys',
  SESSIONS: 'sessions',
  MESSAGE_KEYS_CACHE: 'message_keys_cache',
  MASTER_KEY_SALT: 'master_key_salt',
  ENCRYPTED_MESSAGES: 'encrypted_messages',
  SENT_MESSAGES: 'sent_messages', // 存储发送方的明文（仅发送方可见）
  RECEIVED_MESSAGES: 'received_messages', // 存储接收方的明文（仅接收方可见）
};

let dbInstance = null;
let currentUserId = null;

/**
 * 设置当前用户 ID（用于数据库隔离）
 * @param {string} userId 
 */
export function setCurrentUserId(userId) {
  if (currentUserId !== userId) {
    console.log(`🔄 [CryptoStore] 切换用户: ${currentUserId} → ${userId}`);
    currentUserId = userId;
    dbInstance = null; // 切换用户时清空数据库实例，强制重新打开
  }
}

/**
 * 获取当前用户的数据库名称
 * @returns {string}
 */
function getDBName() {
  if (!currentUserId) {
    console.warn('⚠️ [CryptoStore] 未设置用户 ID，使用默认数据库');
    return DB_NAME_PREFIX;
  }
  return `${DB_NAME_PREFIX}_${currentUserId}`;
}

/**
 * 初始化 IndexedDB
 * @returns {Promise<IDBDatabase>}
 */
export async function initCryptoStore() {
  if (dbInstance) {
    return dbInstance;
  }

  const dbName = getDBName();
  console.log(`🗄️ [CryptoStore] 打开数据库: ${dbName}`);

  return new Promise((resolve, reject) => {
    const request = indexedDB.open(dbName, DB_VERSION);

    request.onerror = () => {
      reject(new Error('Failed to open IndexedDB'));
    };

    request.onsuccess = (event) => {
      dbInstance = event.target.result;
      console.log('CryptoStore initialized');
      resolve(dbInstance);
    };

    request.onupgradeneeded = (event) => {
      const db = event.target.result;

      // 1. user_keys - 用户密钥
      if (!db.objectStoreNames.contains(STORES.USER_KEYS)) {
        const userKeysStore = db.createObjectStore(STORES.USER_KEYS, {
          keyPath: 'key_type',
        });
        userKeysStore.createIndex('created_at', 'created_at');
      }

      // 2. sessions - 会话状态
      if (!db.objectStoreNames.contains(STORES.SESSIONS)) {
        const sessionsStore = db.createObjectStore(STORES.SESSIONS, {
          keyPath: 'contact_id',
        });
        sessionsStore.createIndex('updated_at', 'updated_at');
      }

      // 3. message_keys_cache - 消息密钥缓存
      if (!db.objectStoreNames.contains(STORES.MESSAGE_KEYS_CACHE)) {
        const cacheStore = db.createObjectStore(STORES.MESSAGE_KEYS_CACHE, {
          keyPath: 'id',
          autoIncrement: true,
        });
        cacheStore.createIndex('contact_ratchet_counter', [
          'contact_id',
          'ratchet_key',
          'counter',
        ], { unique: true });
        cacheStore.createIndex('created_at', 'created_at');
      }

      // 4. master_key_salt - 主密钥盐值
      if (!db.objectStoreNames.contains(STORES.MASTER_KEY_SALT)) {
        db.createObjectStore(STORES.MASTER_KEY_SALT, { keyPath: 'id' });
      }

      // 5. encrypted_messages - 加密消息缓存
      if (!db.objectStoreNames.contains(STORES.ENCRYPTED_MESSAGES)) {
        const msgsStore = db.createObjectStore(STORES.ENCRYPTED_MESSAGES, {
          keyPath: 'message_id',
        });
        msgsStore.createIndex('contact_id', 'contact_id');
        msgsStore.createIndex('created_at', 'created_at');
      }

      // 6. sent_messages - 发送方的明文（仅发送方可见，用于历史记录）
      if (!db.objectStoreNames.contains(STORES.SENT_MESSAGES)) {
        const sentMsgsStore = db.createObjectStore(STORES.SENT_MESSAGES, {
          keyPath: 'message_id',
        });
        sentMsgsStore.createIndex('contact_id', 'contact_id');
        sentMsgsStore.createIndex('created_at', 'created_at');
        console.log('✅ 创建 SENT_MESSAGES Object Store');
      } else {
        console.log('ℹ️ SENT_MESSAGES Object Store 已存在');
      }

      // 7. received_messages - 接收方的明文（仅接收方可见，用于历史记录）
      if (!db.objectStoreNames.contains(STORES.RECEIVED_MESSAGES)) {
        const receivedMsgsStore = db.createObjectStore(STORES.RECEIVED_MESSAGES, {
          keyPath: 'message_id',
        });
        receivedMsgsStore.createIndex('contact_id', 'contact_id');
        receivedMsgsStore.createIndex('created_at', 'created_at');
        console.log('✅ 创建 RECEIVED_MESSAGES Object Store');
      } else {
        console.log('ℹ️ RECEIVED_MESSAGES Object Store 已存在');
      }

      console.log('CryptoStore schema upgraded');
    };
  });
}

/**
 * 存储数据
 * @param {string} storeName 
 * @param {Object} data 
 * @returns {Promise<void>}
 */
export async function put(storeName, data) {
  const db = await initCryptoStore();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);
    const request = store.put(data);

    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error);
  });
}

/**
 * 获取数据
 * @param {string} storeName 
 * @param {any} key 
 * @returns {Promise<Object|null>}
 */
export async function get(storeName, key) {
  const db = await initCryptoStore();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([storeName], 'readonly');
    const store = transaction.objectStore(storeName);
    const request = store.get(key);

    request.onsuccess = () => resolve(request.result || null);
    request.onerror = () => reject(request.error);
  });
}

/**
 * 获取所有数据
 * @param {string} storeName 
 * @returns {Promise<Array>}
 */
export async function getAll(storeName) {
  const db = await initCryptoStore();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([storeName], 'readonly');
    const store = transaction.objectStore(storeName);
    const request = store.getAll();

    request.onsuccess = () => resolve(request.result || []);
    request.onerror = () => reject(request.error);
  });
}

/**
 * 删除数据
 * @param {string} storeName 
 * @param {any} key 
 * @returns {Promise<void>}
 */
export async function remove(storeName, key) {
  const db = await initCryptoStore();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);
    const request = store.delete(key);

    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error);
  });
}

/**
 * 清空 Object Store
 * @param {string} storeName 
 * @returns {Promise<void>}
 */
export async function clear(storeName) {
  const db = await initCryptoStore();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);
    const request = store.clear();

    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error);
  });
}

// 导出 Store 名称常量
export { STORES };

