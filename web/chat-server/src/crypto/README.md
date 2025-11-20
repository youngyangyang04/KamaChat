# 端到端加密模块使用指南

## 模块概览

本加密模块实现了完整的端到端加密聊天功能，基于 Signal Protocol 设计，包括：

- ✅ **X3DH 密钥协商**：安全的会话建立
- ✅ **双棘轮算法**：前向安全的消息加密
- ✅ **密钥管理**：本地加密存储
- ✅ **IndexedDB 持久化**：会话状态管理

---

## 快速开始

### 1. 注册时初始化密钥

```javascript
import { initializeUserKeys } from '@/crypto';

// 用户注册时
const password = 'user_password';
const cryptoKeys = await initializeUserKeys(password);

// cryptoKeys 包含：
// - identity_key_public
// - signed_pre_key { key_id, public_key, signature }
// - one_time_pre_keys [{ key_id, public_key }, ...]

// 将这些公钥上传到服务器
await axios.post('/register', {
  account: 'user_account',
  nickname: 'user_nickname',
  password: password,
  ...cryptoKeys,
});
```

### 2. 登录时重新派生主密钥

```javascript
import { loginAndDeriveMasterKey } from '@/crypto';

// 用户登录时
const password = 'user_password';
const masterKey = await loginAndDeriveMasterKey(password);

if (!masterKey) {
  alert('密码错误');
  return;
}

// 将 masterKey 保存到 Vuex store（内存中）
store.commit('setMasterKey', masterKey);
```

### 3. 建立加密会话（发起方）

```javascript
import { createSession } from '@/crypto';
import axios from '@/utils/axios';

// Alice 给 Bob 发第一条消息时
const contactId = 'bob_uuid';
const masterKey = store.state.masterKey;

// 1. 获取 Bob 的公钥束
const response = await axios.get('/crypto/getPublicKeyBundle', {
  params: { user_id: contactId },
});
const bobPublicKeyBundle = response.data.data;

// 2. 建立会话
const initData = await createSession(masterKey, contactId, bobPublicKeyBundle);

// 3. initData 需要包含在首条消息中发送给 Bob
```

### 4. 接受加密会话（接收方）

```javascript
import { acceptSession } from '@/crypto';

// Bob 收到 Alice 的首条加密消息时
const contactId = 'alice_uuid';
const masterKey = store.state.masterKey;

// 从消息中提取初始化数据
const aliceInitData = {
  identity_key: message.sender_identity_key,
  ephemeral_key: message.sender_ephemeral_key,
  ratchet_key: message.ratchet_key,
  used_one_time_pre_key_id: message.used_one_time_pre_key_id,
};

// 建立会话
await acceptSession(masterKey, contactId, aliceInitData);
```

### 5. 发送加密消息

```javascript
import { encryptAndSendMessage } from '@/crypto';

const contactId = 'bob_uuid';
const plaintext = 'Hello, Bob!';

// 加密消息
const encryptedMessage = await encryptAndSendMessage(contactId, plaintext);

// 发送到服务器
await axios.post('/message/sendEncryptedMessage', {
  receiver_id: contactId,
  is_encrypted: true,
  message_type: 'SignalMessage',
  ...encryptedMessage,
});
```

### 6. 接收并解密消息

```javascript
import { receiveAndDecryptMessage } from '@/crypto';

// 收到加密消息时
const contactId = message.sender_id;
const encryptedMessage = {
  ratchet_key: message.ratchet_key,
  counter: message.counter,
  prev_counter: message.prev_counter,
  ciphertext: message.content, // Base64
  iv: message.iv,
  auth_tag: message.auth_tag,
};

// 解密消息
const plaintext = await receiveAndDecryptMessage(contactId, encryptedMessage);

console.log('解密后的消息:', plaintext);
```

---

## API 参考

### 密钥管理

#### `initializeUserKeys(password)`
初始化用户密钥（注册时调用）

**参数**:
- `password` (string): 用户密码

**返回**: Promise<Object>
```javascript
{
  identity_key_public: "Base64...",
  signed_pre_key: {
    key_id: 1,
    public_key: "Base64...",
    signature: "Base64..."
  },
  one_time_pre_keys: [
    { key_id: 1, public_key: "Base64..." },
    // ... 100 个
  ]
}
```

#### `loginAndDeriveMasterKey(password)`
登录时重新派生主密钥

**参数**:
- `password` (string): 用户密码

**返回**: Promise<Uint8Array|null>
- 成功: 主密钥（Uint8Array）
- 失败: null（密码错误）

---

### 会话管理

#### `createSession(masterKey, contactId, remotePublicKeyBundle)`
创建新会话（作为发起方）

**参数**:
- `masterKey` (Uint8Array): 主密钥
- `contactId` (string): 对方用户ID
- `remotePublicKeyBundle` (Object): 对方的公钥束

**返回**: Promise<Object> - 初始化数据（需要发送给对方）

#### `acceptSession(masterKey, contactId, aliceInitData)`
接受会话（作为接收方）

**参数**:
- `masterKey` (Uint8Array): 主密钥
- `contactId` (string): 对方用户ID
- `aliceInitData` (Object): 对方的初始化数据

**返回**: Promise<void>

#### `hasSession(contactId)`
检查会话是否存在

**参数**:
- `contactId` (string): 对方用户ID

**返回**: Promise<boolean>

---

### 消息加密/解密

#### `encryptAndSendMessage(contactId, plaintext)`
加密并准备发送消息

**参数**:
- `contactId` (string): 对方用户ID
- `plaintext` (string): 明文消息

**返回**: Promise<Object>
```javascript
{
  ratchet_key: "Base64...",
  counter: 15,
  prev_counter: 12,
  ciphertext: "Base64...",
  iv: "Base64...",
  auth_tag: "Base64..."
}
```

#### `receiveAndDecryptMessage(contactId, encryptedMessage)`
接收并解密消息

**参数**:
- `contactId` (string): 对方用户ID
- `encryptedMessage` (Object): 加密消息数据

**返回**: Promise<string> - 明文消息

---

## 数据存储

所有敏感数据存储在 IndexedDB 中，包括：

### 1. `user_keys` - 用户密钥
- 身份私钥（加密存储）
- 签名预私钥（加密存储）
- 一次性预私钥（加密存储）

### 2. `sessions` - 会话状态
- 根密钥（Root Key）
- 链密钥（Chain Keys）
- DH 棘轮密钥
- 计数器

### 3. `message_keys_cache` - 消息密钥缓存
- 用于处理乱序消息

### 4. `master_key_salt` - 主密钥盐值
- PBKDF2 的 salt
- 验证数据

### 5. `encrypted_messages` - 加密消息缓存
- 可选，用于性能优化

---

## 安全特性

### ✅ 端到端加密
- 服务器无法解密任何消息内容
- 密钥仅在客户端生成和存储

### ✅ 前向安全
- 每条消息使用不同的密钥
- 密钥链单向派生
- 旧密钥泄露不影响新消息

### ✅ 后向安全（部分）
- DH 棘轮定期更新
- 新密钥泄露不影响旧消息

### ✅ 身份验证
- Ed25519 签名验证
- 防止公钥替换攻击

### ✅ 认证加密
- AES-256-GCM
- 同时提供加密和完整性

---

## 性能优化建议

### 1. 密钥预加载
```javascript
// 登录时预加载常用联系人的会话
const contacts = ['bob_uuid', 'carol_uuid'];
for (const contactId of contacts) {
  const session = await get(STORES.SESSIONS, contactId);
  // 会话已加载到内存
}
```

### 2. 使用 Web Worker
```javascript
// 在 Web Worker 中执行加密/解密
// 避免阻塞 UI 线程
const worker = new Worker('/crypto-worker.js');
worker.postMessage({ action: 'encrypt', data: plaintext });
```

### 3. 批量解密
```javascript
// 批量解密历史消息
const messages = [...]; // 多条加密消息
const plaintexts = await Promise.all(
  messages.map(msg => receiveAndDecryptMessage(contactId, msg))
);
```

---

## 故障排除

### 问题：密码验证失败
```javascript
const masterKey = await loginAndDeriveMasterKey(password);
if (!masterKey) {
  console.error('密码错误或数据损坏');
  // 可能需要重新注册
}
```

### 问题：会话不存在
```javascript
try {
  await encryptAndSendMessage(contactId, plaintext);
} catch (error) {
  if (error.message.includes('会话不存在')) {
    // 需要先建立会话
    await createSession(masterKey, contactId, publicKeyBundle);
  }
}
```

### 问题：解密失败
```javascript
try {
  const plaintext = await receiveAndDecryptMessage(contactId, encryptedMsg);
} catch (error) {
  console.error('解密失败:', error);
  // 可能原因：
  // 1. 消息损坏
  // 2. 会话状态不同步
  // 3. 密钥错误
}
```

---

## 开发注意事项

1. **主密钥不存储**: 主密钥仅存在于内存中，用户退出登录后会被清除
2. **定期轮换**: 签名预密钥建议每 30 天轮换一次
3. **密钥补充**: 一次性预密钥少于 20 个时自动补充
4. **乱序处理**: 当前实现简化了乱序消息处理，完整实现需要缓存跳过的消息密钥
5. **DH 棘轮频率**: 每 100 条消息执行一次 DH 棘轮更新

---

## 技术栈

- **加密库**: [TweetNaCl](https://github.com/dchest/tweetnacl-js)
- **密钥派生**: PBKDF2 (100,000 次迭代)
- **签名算法**: Ed25519
- **密钥交换**: Curve25519 (ECDH)
- **对称加密**: AES-256-GCM
- **密钥派生函数**: HKDF-SHA256

---

## 参考资料

- [Signal Protocol 规范](https://signal.org/docs/)
- [X3DH 密钥协商](https://signal.org/docs/specifications/x3dh/)
- [Double Ratchet Algorithm](https://signal.org/docs/specifications/doubleratchet/)
- [TweetNaCl 文档](https://tweetnacl.js.org/)

---

**版本**: 1.0  
**最后更新**: 2025-11-19

