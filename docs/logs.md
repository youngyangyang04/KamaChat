## 修改日志

### 2025-11-18 - 数据库从 MySQL 迁移到 PostgreSQL
- 修改 `configs/config.toml`：配置改为 PostgreSQL（端口 5432，用户 postgres，空密码）
- 修改 `internal/config/config.go`：`MysqlConfig` → `PostgresqlConfig`，配置加载路径改为本地
- 修改 `internal/dao/gorm.go`：数据库驱动改为 `gorm.io/driver/postgres`，更新 DSN 连接格式（支持空密码）
- 修改所有 model 文件：将时间字段类型从 `type:datetime` 改为 `type:timestamp`（PostgreSQL 不支持 datetime）
- 修改 `cmd/kama_chat_server/main.go`：改用 HTTP（注释 TLS）方便本地测试
- 修改 `web/chat-server/vue.config.js`：前端改用 HTTP 在 8080 端口运行
- 添加依赖：`gorm.io/driver/postgres v1.6.0` 及相关 PostgreSQL 驱动包

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