## 修改日志

### 2025-11-19 - 实现单聊端到端加密核心功能（阶段1-3完成）

#### 数据库层
- 创建数据库迁移脚本 `migrations/001_add_e2ee_support.sql`
- 创建 Go models：`OneTimePreKey`, `KeyReplenishmentLog`

#### 后端服务层
- 创建公钥管理服务：`internal/service/gorm/crypto_key_service.go`
- 创建加密 API 控制器：`api/v1/crypto_key.go`
- 注册路由：
  - `GET /crypto/getPublicKeyBundle`：获取公钥束
  - `GET /crypto/getOneTimePreKeyCount`：查询剩余密钥
  - `POST /crypto/replenishOneTimePreKeys`：补充密钥
  - `POST /crypto/rotateSignedPreKey`：轮换签名预密钥

#### 前端加密模块（完整实现）
- **基础设施**：
  - `crypto/keyDerivation.js`：PBKDF2 主密钥派生、AES-GCM 加密
  - `crypto/keyGeneration.js`：Ed25519/Curve25519 密钥生成
  - `crypto/cryptoStore.js`：IndexedDB 封装（5个存储）
  - `crypto/keyManager.js`：统一密钥管理接口
- **核心协议**：
  - `crypto/x3dh.js`：X3DH 密钥协商协议（会话建立）
  - `crypto/doubleRatchet.js`：双棘轮算法（消息加密）
  - `crypto/sessionManager.js`：会话管理器（整合所有功能）
- **统一导出**：
  - `crypto/index.js`：统一导出所有加密API

#### 依赖
- 添加 `tweetnacl@1.0.3`（前端加密库）✅

### 2025-11-19 - 集成注册/登录流程支持加密

#### 后端集成
- 修改 `internal/service/gorm/user_info_service.go`：
  - 新增 `RegisterWithCrypto` 函数（使用事务保存用户和密钥）
  - 保持原有 `Register` 函数向后兼容
- 新增控制器 `RegisterWithCrypto` in `api/v1/user_info_controller.go`
- 注册路由 `POST /registerWithCrypto`

#### 前端集成
- 修改 `web/chat-server/src/views/access/Register.vue`：
  - 调用 `initializeUserKeys` 生成加密密钥（100 个 OTP 密钥）
  - 调用 `/registerWithCrypto` 接口上传公钥
  - 注册成功后重新派生主密钥保存到内存
- 修改 `web/chat-server/src/views/access/Login.vue`：
  - 登录成功后尝试重新派生主密钥
  - 支持加密和非加密用户无缝切换
- 修改 `web/chat-server/src/store/index.js`：
  - 新增 `masterKey` state（仅内存，不持久化）
  - 新增 `setMasterKey` mutation
  - 修改 `cleanUserInfo` mutation 清除主密钥

#### 特性
- ✅ 新用户注册自动启用加密
- ✅ 旧用户登录保持普通模式（向后兼容）
- ✅ 主密钥仅存储在内存，退出登录自动清除（安全）
- ✅ 事务保证用户创建和密钥保存的原子性

### 2025-11-19 - 集成消息加密/解密功能（阶段4-5完成）

#### 后端集成
- 创建加密消息 API：`api/v1/encrypted_message.go`
  - `POST /message/sendEncryptedMessage`：发送加密消息
- 修改 `internal/model/message.go`：添加加密相关字段（11个新字段）
- 修改 `internal/model/user_info.go`：添加加密公钥字段（6个新字段）
- 修改 `internal/service/gorm/message_service.go`：GetMessageList 返回加密字段
- 修改 `internal/dto/respond/get_message_list_respond.go`：响应包含加密字段
- 创建 DTO：`send_encrypted_message_request.go`
- 注册路由：`POST /message/sendEncryptedMessage`
- 更新 `internal/dao/gorm.go`：AutoMigrate 包含新 model

#### 前端集成
- 创建消息解密工具：`web/chat-server/src/utils/messageDecryptor.js`
  - `decryptMessage`：解密单条消息
  - `decryptMessageList`：批量解密消息列表
  - 自动处理 PreKeyMessage 会话建立
- 修改 `web/chat-server/src/views/chat/contact/ContactChat.vue`：
  - 导入加密模块
  - `sendMessage` 函数：检测加密状态，自动选择加密/明文发送
  - `sendEncryptedMessage` 函数：完整的加密发送流程
    - 检查会话存在性
    - 首次会话：获取公钥束、建立会话（X3DH）
    - 加密消息（双棘轮）
    - 发送到服务器
  - `getMessageList` 函数：获取后自动解密
  - UI 添加 🔒 加密状态指示器

#### 数据库迁移
- 执行 `migrations/001_add_e2ee_support.sql`
- 创建表：`one_time_pre_keys`, `key_replenishment_log`
- 扩展表：`user_info` (6字段), `message` (11字段)
- 创建索引：加密消息查询优化

#### 功能特性
- ✅ 自动检测加密状态（有主密钥则加密）
- ✅ 首次会话自动建立（X3DH 协议）
- ✅ 消息自动加密/解密（双棘轮算法）
- ✅ 向后兼容明文消息
- ✅ UI 显示加密状态（🔒 图标）
- ✅ 批量解密历史消息
- ✅ 服务器无法看到明文（端到端加密）

### 2025-11-18 - 数据库从 MySQL 迁移到 PostgreSQL
- 修改 `configs/config.toml`：配置改为 PostgreSQL（端口 5432，用户 postgres，空密码）
- 修改 `internal/config/config.go`：`MysqlConfig` → `PostgresqlConfig`，配置加载路径改为本地
- 修改 `internal/dao/gorm.go`：数据库驱动改为 `gorm.io/driver/postgres`，更新 DSN 连接格式（支持空密码）
- 修改所有 model 文件：将时间字段类型从 `type:datetime` 改为 `type:timestamp`（PostgreSQL 不支持 datetime）
- 修改 `cmd/kama_chat_server/main.go`：改用 HTTP（注释 TLS）方便本地测试
- 修改 `web/chat-server/vue.config.js`：前端改用 HTTP 在 8080 端口运行
- 添加依赖：`gorm.io/driver/postgres v1.6.0` 及相关 PostgreSQL 驱动包

### 2025-11-19 - 修复发送消息时后端 panic
- 修复 `internal/service/chat/server.go`：normalizePath 函数处理外部 URL 时 slice 越界
- **问题**：发送消息时后端崩溃，`panic: runtime error: slice bounds out of range [-1:]`
- **原因**：头像 URL 是外部 HTTPS 链接，不包含 `/static/`，`staticIndex` 为 -1，导致 `path[-1:]` 越界
- **修复**：检查路径是否为 HTTP/HTTPS URL，如果是则直接返回；如果找不到 `/static/` 也返回原路径

### 2025-11-19 - 修复 WebSocket 连接需要 JWT 认证
- 修改 `pkg/middleware/jwt_auth.go`：JWT 中间件支持从 query 参数读取 token
- 修改前端 WebSocket 连接：在 URL 中添加 token 参数
  - `web/chat-server/src/App.vue`
  - `web/chat-server/src/views/access/Login.vue`
  - `web/chat-server/src/views/access/Register.vue`
  - `web/chat-server/src/views/access/SmsLogin.vue`
- 修改 `web/chat-server/src/views/chat/contact/ContactChat.vue`：sendMessage 函数添加完整的错误检查和调试日志
- **问题**：WebSocket 连接失败，点击发送按钮显示"WebSocket 连接已断开"
- **原因**：/wss 路由需要 JWT 认证，但 WebSocket 连接时未传递 token
- **修复**：后端支持从 query 参数读取 token，前端在 WebSocket URL 中添加 token 参数

### 2025-11-19 - 修复退出登录 401 错误
- 修复 `web/chat-server/src/components/NavigationModal.vue`：调整 logout 函数执行顺序
- **问题**：退出登录时报 401 错误
- **原因**：先清除了 token，再调用退出登录接口，导致接口无权限
- **修复**：改为先调用接口，成功后再清除 token，并添加错误处理

### 2025-11-19 - 修复 WebSocket onerror 回调错误
- 修复 `web/chat-server/src/App.vue`：WebSocket onerror 回调添加 error 参数
- 修复 `web/chat-server/src/views/access/Login.vue`：WebSocket onerror 回调添加 error 参数
- 修复 `web/chat-server/src/views/access/Register.vue`：WebSocket onerror 回调添加 error 参数
- 修复 `web/chat-server/src/views/access/SmsLogin.vue`：WebSocket onerror 回调添加 error 参数
- **问题**：刷新页面时报错 "ReferenceError: error is not defined"
- **原因**：onerror 回调函数未声明 error 参数就使用了 error 变量

### 2025-11-19 - 引入 JWT Token 认证机制
**后端修改：**
- 新增 `pkg/util/jwt/jwt.go`：JWT token 生成、解析和刷新功能
- 新增 `pkg/middleware/jwt_auth.go`：JWT 认证中间件，验证请求头中的 token
- 修改 `configs/config.toml`：添加 JWT 配置（secretKey、expireHours、refreshHours）
- 修改 `internal/config/config.go`：添加 `JWTConfig` 结构体
- 修改 `pkg/util/jwt/jwt.go`：在 init 函数中自动读取配置并初始化（通过依赖链自动执行）
- 修改 `internal/service/gorm/user_info_service.go`：
  - 登录时生成 token 并返回
  - 注册时生成 token 并返回
- 修改 `internal/dto/respond/login_respond.go` 和 `register_respond.go`：添加 `Token` 字段
- 修改 `internal/https_server/https_server.go`：
  - 登录和注册接口无需认证
  - 其他所有接口需要 JWT 认证
- 添加依赖：`github.com/golang-jwt/jwt/v5`

**前端修改：**
- 修改 `web/chat-server/src/store/index.js`：
  - 添加 `token` 状态管理
  - 登录/注册时自动保存 token
  - 退出登录时清除 token
- 新增 `web/chat-server/src/utils/axios.js`：
  - 请求拦截器自动在请求头添加 Authorization: Bearer <token>
  - 响应拦截器处理 401 未授权错误，自动跳转登录
- 修改所有 Vue 组件：将 `import axios from 'axios'` 改为 `import axios from '@/utils/axios'`
  - views 目录：Login.vue, Register.vue, SmsLogin.vue, Manager.vue, OwnInfo.vue, SessionList.vue, ContactList.vue, ContactChat.vue
  - components 目录：SetAdminModal.vue, DisableUserModal.vue, NavigationModal.vue, ContactListModal.vue 等

**Token 特性：**
- 有效期：24小时（可配置）
- 格式：JWT (HS256)
- 携带信息：UUID、账号、昵称、是否管理员
- 自动续期：可通过刷新 token API 续期（7天有效期）

### 2025-11-19 - 添加密码加密功能
- 新增 `pkg/util/password/password.go`：使用 bcrypt 算法实现密码加密和验证
- 修改 `internal/service/gorm/user_info_service.go`：
  - 注册时使用 bcrypt 加密密码后存储
  - 登录时使用 bcrypt 验证密码
- 添加依赖：`golang.org/x/crypto/bcrypt`
- **重要**：数据库中原有的明文密码将无法登录，需要重新注册

### 2025-11-19 - 添加 HTTPS 开关配置
- 修改 `configs/config.toml`：添加 `enableHTTPS`、`certFile`、`keyFile` 配置项，可通过配置切换 HTTP/HTTPS
- 修改 `internal/config/config.go`：`MainConfig` 结构体添加 HTTPS 相关字段
- 修改 `internal/https_server/https_server.go`：根据配置决定是否启用 TLS 重定向中间件
- 修改 `cmd/kama_chat_server/main.go`：根据配置自动选择 HTTP 或 HTTPS 启动模式

### 2025-11-19 - 修复端口冲突问题
- 修改 `configs/config.toml`：后端端口从 8000 改为 8888（避免 Windows 保留端口范围 7959-8058）
- 修改 `web/chat-server/src/store/index.js`：前端配置改为本地开发地址 `http://127.0.0.1:8888` 和 `ws://127.0.0.1:8888`

### 2025-11-18 - 简化注册登录流程
**后端修改：**
- 修改 `internal/model/user_info.go`：添加 `Account` 字段（唯一索引），`Password` 字段长度改为 varchar(100)
- 修改 `internal/dto/request/register_request.go`：注册只需 account、nickname、password
- 修改 `internal/dto/request/login_request.go`：登录只需 account、password
- 修改 `internal/dto/respond/login_respond.go` 和 `register_respond.go`：添加 `Account` 字段
- 修改 `internal/service/gorm/user_info_service.go`：
  - 添加 `checkAccountExist()` 方法检查账号唯一性
  - `Register()` 移除短信验证，使用账号注册
  - `Login()` 使用账号登录而非手机号

**前端修改：**
- 修改 `web/chat-server/src/views/access/Register.vue`：
  - 注册表单改为：账号（3-20字符）、昵称（1-20字符）、密码（6-50字符）
  - 移除短信验证码相关代码
  - 移除"验证码登录"按钮
- 修改 `web/chat-server/src/views/access/Login.vue`：
  - 登录表单改为：账号、密码
  - 移除手机号验证逻辑
  - 移除"验证码登录"按钮

**字段限制：**
- 账号：3-20个字符，唯一
- 昵称：1-20个字符，可重复
- 密码：6-50个字符

---