# 端到端加密测试指南

## 测试环境准备

### 1. 启动后端服务
```powershell
go run .\cmd\kama_chat_server\main.go
```

### 2. 启动前端服务
```powershell
cd web\chat-server
npm run serve
```

### 3. 打开浏览器
访问 `http://localhost:8080`

---

## 测试场景

### 场景 1：新用户注册（启用加密）

1. **注册用户 A**
   - 账号：`alice001`
   - 昵称：`Alice`
   - 密码：`password123`
   - 点击"注册"

2. **观察**：
   - 控制台应显示：`"正在生成加密密钥，请稍候..."`
   - 控制台应显示：`"加密密钥生成成功"`
   - 注册成功提示：`"注册成功 (端到端加密已启用)"`
   - 控制台显示：`"主密钥已保存到内存"`

3. **验证 IndexedDB**：
   - 打开浏览器开发者工具
   - 切换到 Application → IndexedDB
   - 应该看到 `wework_crypto` 数据库
   - 展开后看到 5 个 Object Stores：
     - `user_keys`：包含加密的私钥（101 条记录：1 个身份密钥 + 1 个签名预密钥 + 100 个 OTP 密钥）
     - `master_key_salt`：主密钥 salt
     - `sessions`：会话状态（暂时为空）
     - `message_keys_cache`：消息密钥缓存（暂时为空）
     - `encrypted_messages`：消息缓存（暂时为空）

4. **验证数据库**：
   ```sql
   -- 在 PostgreSQL 中查询
   SELECT uuid, account, nickname, identity_key_public, signed_pre_key_id 
   FROM user_info 
   WHERE account = 'alice001';
   
   -- 应该看到 identity_key_public 和 signed_pre_key_public 有值
   
   SELECT COUNT(*) FROM one_time_pre_keys WHERE user_id = 'U...';
   -- 应该返回 100
   ```

---

### 场景 2：新用户 B 注册

1. **退出登录**（用户 A）

2. **注册用户 B**
   - 账号：`bob001`
   - 昵称：`Bob`
   - 密码：`password456`

3. **观察**：
   - 同样的加密初始化过程
   - 成功后显示端到端加密已启用

---

### 场景 3：加密消息发送（首次会话）

1. **用户 A 登录**
   - 账号：`alice001`
   - 密码：`password123`
   - 观察控制台：`"主密钥已重新派生并保存到内存（加密功能已启用）"`

2. **添加用户 B 为好友**
   - 搜索：`bob001`
   - 发送好友申请

3. **用户 B 登录并通过好友申请**

4. **用户 A 发送第一条消息**
   - 打开与 Bob 的聊天
   - 输入：`Hello Bob! This is an encrypted message.`
   - 点击"发送"

5. **观察控制台**（关键输出）：
   ```
   准备发送消息： Hello Bob! This is an encrypted message.
   🔒 准备发送加密消息...
   会话不存在，正在建立加密会话...
   获取到公钥束: {...}
   Alice 开始 X3DH 协商...
   Alice X3DH 协商完成
   ✅ 加密会话已建立
   消息已加密: {...}
   发送加密消息到服务器...
   ✅ 加密消息发送成功
   ```

6. **验证界面**：
   - 消息前应显示 🔒 图标（绿色）
   - 消息内容正常显示

7. **验证 IndexedDB**：
   - `sessions` 应该有一条记录（contactId = Bob 的 UUID）
   - 查看会话状态，包含 root_key, chain_keys 等

8. **验证服务器数据库**：
   ```sql
   SELECT content, is_encrypted, message_type, ratchet_key 
   FROM message 
   WHERE send_id = 'Alice的UUID' AND receive_id = 'Bob的UUID';
   
   -- content 应该是 Base64 密文，无法识别原文
   -- is_encrypted 应该是 true
   -- message_type 应该是 'PreKeyMessage'
   ```

---

### 场景 4：加密消息接收

1. **用户 B 刷新页面**
   - 重新登录（主密钥重新派生）

2. **打开与 Alice 的聊天**

3. **观察控制台**：
   ```
   获取消息列表...
   准备解密 1 条消息...
   🔒 接收 PreKeyMessage，建立会话...
   Bob 开始 X3DH 协商...
   Bob X3DH 协商完成
   ✅ 会话已建立
   ✅ 消息已解密
   ✅ 1 条消息解密完成
   消息列表解密完成
   ```

4. **验证界面**：
   - 应该看到 Alice 发送的消息：`Hello Bob! This is an encrypted message.`
   - 消息前显示 🔒 图标

5. **验证 IndexedDB**：
   - Bob 的 `sessions` 应该有一条记录（contactId = Alice 的 UUID）
   - Bob 的一个 OTP 私钥被使用并删除（`user_keys` 减少 1 条）

---

### 场景 5：后续加密消息（双棘轮）

1. **用户 B 回复**
   - 输入：`Hi Alice! I got your encrypted message!`
   - 点击"发送"

2. **观察控制台**（B 的控制台）：
   ```
   🔒 准备发送加密消息...
   会话已存在，直接加密
   消息已加密: {...}
   ✅ 加密消息发送成功
   ```

3. **用户 A 接收**
   - 刷新页面或实时接收（如果实现了 WebSocket 集成）
   - 应该看到解密后的消息

4. **连续发送多条消息**
   - 发送 10 条消息
   - 观察每条消息的 `counter` 递增

5. **验证数据库**：
   ```sql
   SELECT counter, ratchet_key, message_type 
   FROM message 
   WHERE session_id = '会话ID' 
   ORDER BY counter;
   
   -- counter 应该是 0, 1, 2, 3, ...
   -- message_type 第一条是 'PreKeyMessage'，后续是 'SignalMessage'
   ```

---

### 场景 6：DH 棘轮更新（100 条消息后）

1. **连续发送 100 条消息**
   - 可以写个脚本自动化

2. **观察控制台**：
   ```
   第 99 条: counter=99
   第 100 条: 
     执行 DH 棘轮更新...
     counter=100
   ```

3. **验证**：
   - `ratchet_key` 应该更新为新值
   - 会话继续正常工作

---

### 场景 7：向后兼容（旧用户）

1. **创建旧用户**（使用原有注册接口）
   - 直接调用 `/register`（非 `/registerWithCrypto`）
   - 或使用之前注册的旧用户

2. **旧用户登录**
   - 控制台显示：`"未找到加密密钥，使用普通模式"`
   - 消息正常收发（明文）

3. **旧用户与新用户聊天**
   - 新用户发送：fallback 到明文（或提示对方未启用加密）
   - 消息前没有 🔒 图标

---

## 调试技巧

### 1. 查看主密钥状态
```javascript
// 在浏览器控制台输入
store.state.masterKey
// 应该是 Uint8Array(32) 或 null
```

### 2. 查看会话状态
```javascript
// 在浏览器控制台输入
const db = await indexedDB.databases();
console.log(db); // 应该看到 wework_crypto
```

### 3. 手动解密测试
```javascript
import { decryptMessage } from '@/utils/messageDecryptor';

// 测试解密单条消息
const message = { /* 从数据库获取的加密消息 */ };
const decrypted = await decryptMessage(message);
console.log(decrypted.content);
```

### 4. 查看服务器日志
```
保存用户公钥成功, userId: U...
标记一次性预密钥为已使用成功
加密消息发送成功
```

---

## 常见问题

### Q: 注册时提示"生成加密密钥失败"
**A**: 检查浏览器是否支持 Web Crypto API（需要 HTTPS 或 localhost）

### Q: 消息解密失败
**A**: 可能原因：
- 会话状态不同步（清除 IndexedDB 重试）
- 消息损坏
- 密钥错误

解决方法：
```javascript
// 清除所有加密数据
import { clearAllCryptoData } from '@/crypto';
await clearAllCryptoData();
```

### Q: 提示"对方未启用加密功能"
**A**: 对方是旧用户或使用旧注册接口。需要对方重新注册启用加密。

### Q: 查看原始密文
**A**: 在数据库中查询 message 表，content 字段存储的是 Base64 密文。

---

## 性能指标

### 预期性能
- 密钥生成（注册）：2-5 秒
- X3DH 协商（首次会话）：< 200ms
- 消息加密：< 20ms
- 消息解密：< 20ms
- 批量解密 100 条消息：< 2 秒

### 监控指标
- OTP 密钥剩余数量：`GET /crypto/getOneTimePreKeyCount`
- 会话建立成功率
- 加密/解密失败次数

---

## 安全验证

### 1. 服务器无法解密验证
```sql
-- 在数据库中查看消息内容
SELECT content FROM message WHERE is_encrypted = true LIMIT 1;

-- 应该看到类似这样的密文（Base64）：
-- "xK8j2Hs9pL4mQ..."
-- 服务器管理员无法看到明文
```

### 2. 密钥隔离验证
```javascript
// 主密钥只存在于内存
console.log(sessionStorage.getItem('masterKey')); // null

// 退出登录后主密钥应被清除
store.commit('cleanUserInfo');
console.log(store.state.masterKey); // null
```

### 3. 前向安全验证
- 发送消息后，刷新页面
- IndexedDB 中的 `message_keys_cache` 应该为空
- 无法从当前状态推导历史消息密钥

---

## 下一步优化建议

1. **WebSocket 实时加密消息推送**
   - 当前实现：加密消息存储到数据库
   - 改进：通过 WebSocket 实时推送加密消息给在线用户

2. **性能优化**
   - 使用 Web Worker 进行加密/解密
   - 批量解密优化

3. **用户体验**
   - 添加"正在建立安全连接"动画
   - 显示对方是否启用加密
   - 密钥安全码验证（防止中间人攻击）

4. **密钥管理**
   - OTP 密钥自动补充（< 20 个时）
   - 签名预密钥自动轮换（30 天）

5. **群组加密**
   - Sender Keys 实现
   - 群组密钥分发机制

---

**测试状态**: ✅ 核心功能已完整实现  
**建议测试时长**: 30-60 分钟  
**最后更新**: 2025-11-19

