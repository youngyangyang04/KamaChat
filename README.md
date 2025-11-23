# WeWork Reborn

一个基于 [KamaChat](https://github.com/youngyangyang04/KamaChat) 开发的即时通讯系统，支持端到端加密（E2EE）的单聊功能。

## 📋 项目简介

WeWork Reborn 是一个功能完整的即时通讯系统，在 kama-chat 开源项目的基础上，新增了端到端加密功能，确保用户通信的隐私和安全。系统采用前后端分离架构，支持实时消息推送、用户管理、会话管理等功能。

## ✨ 主要特性

- 🔐 **端到端加密（E2EE）**：基于 Signal Protocol 设计，实现 X3DH 密钥协商和双棘轮算法
- 💬 **实时聊天**：基于 WebSocket 的实时消息推送
- 👥 **用户管理**：用户注册、登录、个人信息管理
- 📱 **会话管理**：会话列表、消息历史记录
- 🔑 **密钥管理**：主密钥可选持久化存储，支持会话刷新后自动恢复
- 📦 **消息队列**：支持 Kafka 和 Channel 两种消息模式
- 🔒 **安全认证**：JWT Token 认证机制

## 🛠️ 技术栈

### 后端
- **语言**：Go 1.24+
- **Web 框架**：Gin
- **数据库**：PostgreSQL
- **缓存**：Redis
- **消息队列**：Kafka（可选）
- **WebSocket**：Gorilla WebSocket
- **日志**：Zap
- **ORM**：GORM

### 前端
- **框架**：Vue 3
- **UI 组件库**：Element Plus
- **状态管理**：Vuex
- **路由**：Vue Router
- **HTTP 客户端**：Axios
- **加密库**：TweetNaCl.js

## 📁 项目结构

```
wework-reborn/
├── api/v1/                    # API 控制器
│   ├── chatroom_controller.go
│   ├── crypto_key.go          # 加密密钥 API
│   ├── encrypted_message.go   # 加密消息 API
│   ├── message_controller.go
│   ├── user_info_controller.go
│   └── ...
├── cmd/kama_chat_server/      # 主程序入口
│   └── main.go
├── configs/                    # 配置文件
│   └── config.toml
├── docs/                      # 文档
│   ├── 单聊加密实现总结.md
│   ├── 端到端加密方案.md
│   └── ...
├── internal/                   # 内部模块
│   ├── config/                # 配置管理
│   ├── dao/                   # 数据访问层
│   ├── dto/                   # 数据传输对象
│   ├── https_server/          # HTTP 服务器
│   ├── model/                 # 数据模型
│   ├── service/               # 业务逻辑层
│   │   ├── chat/              # 聊天服务
│   │   ├── gorm/              # GORM 服务
│   │   └── ...
│   └── ...
├── migrations/                # 数据库迁移脚本
│   └── 001_add_e2ee_support.sql
├── pkg/                       # 公共包
│   ├── middleware/            # 中间件
│   ├── util/                  # 工具函数
│   └── ...
├── web/chat-server/           # 前端项目
│   ├── src/
│   │   ├── crypto/            # 加密核心模块
│   │   │   ├── sessionManager.js
│   │   │   ├── doubleRatchet.js
│   │   │   ├── x3dh.js
│   │   │   └── ...
│   │   ├── views/             # 页面组件
│   │   ├── components/        # 公共组件
│   │   └── ...
│   └── package.json
└── README.md
```

## 🚀 快速开始

### 环境要求

- Go 1.24+
- Node.js 16+
- PostgreSQL 12+
- Redis 6+


## ⚙️ 配置说明

### 后端配置（configs/config.toml）

- **mainConfig**：服务器基本配置（端口、HTTPS 等）
- **postgresqlConfig**：PostgreSQL 数据库配置
- **redisConfig**：Redis 缓存配置
- **kafkaConfig**：Kafka 消息队列配置（可选）
- **jwtConfig**：JWT Token 配置
- **authCodeConfig**：短信验证码配置（阿里云）
- **logConfig**：日志配置

### 前端配置（web/chat-server/src/store/index.js）

- **backendUrl**：后端 API 地址
- **wsUrl**：WebSocket 服务器地址

## 🔐 端到端加密功能

### 加密协议

系统实现了基于 Signal Protocol 的端到端加密：

- **X3DH 密钥协商**：用于建立初始加密会话
- **双棘轮算法（Double Ratchet）**：提供前向安全性和后向安全性
- **AES-GCM 加密**：用于消息内容的对称加密

### 密钥管理

- **主密钥（Master Key）**：从用户密码派生，用于加密/解密私钥
- **身份密钥（Identity Key）**：长期密钥对，用于身份验证
- **签名预密钥（Signed Pre-Key）**：用于 X3DH 协商
- **一次性预密钥（One-Time Pre-Key）**：提供前向安全性
- **棘轮密钥（Ratchet Key）**：临时密钥对，用于双棘轮算法

### 主密钥持久化

用户可以在设置页面选择是否保存主密钥到本地存储：

- **启用**：主密钥保存在 sessionStorage，刷新页面后自动恢复，无需重新输入密码
- **禁用**：主密钥仅保存在内存中，刷新页面后需要重新登录

详细实现文档请参考：`docs/单聊加密实现总结.md`

## 📡 API 文档

### 认证相关
- `POST /login` - 用户登录
- `POST /register` - 用户注册
- `POST /registerWithCrypto` - 带加密密钥的注册

### 用户相关
- `POST /user/getUserInfo` - 获取用户信息
- `POST /user/updateUserInfo` - 更新用户信息
- `POST /user/getUserInfoList` - 获取用户列表

### 消息相关
- `POST /message/sendMessage` - 发送消息
- `POST /message/getMessageList` - 获取消息列表
- `POST /message/sendEncryptedMessage` - 发送加密消息

### 加密相关
- `POST /crypto/getPublicKeyBundle` - 获取公钥束
- `POST /crypto/uploadPublicKeyBundle` - 上传公钥束

更多 API 详情请查看 `api/v1/` 目录下的控制器文件。


## 📄 许可证

本项目基于 kama-chat 开源项目开发，请遵循原项目的许可证要求。


---

**注意**：本项目仅供学习和研究使用，生产环境使用前请进行充分的安全评估和测试。

