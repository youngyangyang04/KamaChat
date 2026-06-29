# 分布式部署的仿微信项目GoChat

> **本项目目前只在[知识星球](https://programmercarl.com/other/kstar.html)答疑并维护**。

这次发布一个非常硬核的Go项目，分布式部署的仿微信项目：GoChat

越来越多的大厂开始在核心业务中采用 Go 语言，例如在云计算、分布式系统、微服务架构等领域，Go 语言都有着出色的表现。

其实我从24年就开始规划Go的项目，现在这个项目终于打磨好了。

这个项目整个前后端代码加起来有1.5w行（前后端分离），是一个功能非常庞大的项目。

如果对前后端有了解的同学，主要是理解业务流程，每天抽5h，应该在一个月左右能完全掌握这个项目。

如果对go、vue都不太熟悉，对前后端这种api调用方式不太熟悉，可能还需要提前大概去走一遍go和vue，具体时间得看你学基础的耗时，建议通过ai来快速学习。

在开发这个项目，工时大约在200h-300h左右（**相当于24小时都在工作，干了十天**）。

## 涉及技术面广

此项目构建了一个**全面且复杂的即时通讯系统**，前后端融合了广泛的技术体系，深入研究各个技术模块，能让人对即时通讯领域形成系统而完整的认知。

下面将详细阐述各个技术板块及其具体实现功能：

1、 **后台开发者管理**

这一模块赋予后台开发者强大的管理权限，主要用于对用户群聊和用户角色进行管控。

开发者可以对用户群聊执行禁用、启用和删除操作，以确保群聊的正常秩序和合规性。同时，开发者还能根据实际需求，将普通用户设置为管理员，协助进行系统管理。

2、**类似聊天软件的联系人体系**

该体系模拟了常见聊天软件的联系人管理功能，为用户提供了丰富的社交互动选项。

用户可以自由添加或删除联系人，还能对联系人进行拉黑操作。

当用户想要添加新联系人时，可以发起申请，对方则有权选择同意或拒绝该申请。

3、 **单聊与群聊**

单聊和群聊是即时通讯系统的核心功能之一，其实现依赖于后端的聊天服务器。

无论是单聊消息还是群聊消息，都会被发送到后端服务器，由服务器负责将消息准确无误地转发给相应的接收方。

4、 **多种消息类型的上传和下载（文本、文件、视频**）

系统支持多种类型的消息交互，包括文本、文件和视频。用户可以方便地上传和下载这些不同类型的消息，满足多样化的沟通需求。

5、 **Kafka**

Kafka 在项目中扮演着消息传输的重要角色。

它作为一个高效的消息队列系统，负责将客户端（client）发送的消息转发到服务器（server），确保消息的可靠传输和处理。

6、 **语音视频通话**

语音视频通话功能为用户提供了更加直观和实时的沟通方式。

用户可以发起通话邀请，对方可以选择接受或拒绝邀请。在通话过程中，任何一方都可以随时挂断通话，结束交流。

7、 **WebSocket**

WebSocket 为前端和客户端之间建立了实时、双向的通信连接。

它负责接收前端发送给客户端的消息，并将客户端的消息传递给服务器，实现了消息的高效传输和即时响应。

通过参与这个项目，对上述每个技术点进行深入钻研，开发者可以全面了解即时通讯系统的架构、原理和实现细节，从而在该领域积累丰富的经验和技能。

## 系统展示

单聊：

以下我只放了部分视频演示，该项目的各个部分都是视频演示，共计三十多个视频展示：


![image](https://file1.kamacoder.com/i/web/2025-08-18_12-28-59.jpg)

更多视频展示 可以在项目专栏的【项目展示以及测试】观看
![image](https://file1.kamacoder.com/i/web/2025-08-18_12-29-34.jpg)

## GoChat项目精讲

该项目的专栏是[知识星球](https://programmercarl.com/other/kstar.html)录友专享的。

项目专栏依然是将 「简历写法」给大家列出来了，大家学完就可以参考这个来写简历：

![image](https://file1.kamacoder.com/i/web/20250407103128.png)

做完该项目，面试中大概率会有哪些面试问题，以及如何回答，也列出好了。

面试问题，都是星球录友那这个项目面试遇到的实际问题：

![image](https://file1.kamacoder.com/i/web/20250407103212.png)

专栏中的项目面试题都掌握的话，这个项目在面试中基本没问题。

很多录友在做项目的时候，把项目运行起来 就是第一大难点！

本项目运行起来 需要依赖的环境很多，所以我给大家准备的 自动化环境配置脚本， **项目运行环境，一键配置！ 不需要大家去处理环境问题了**：

![image](https://file1.kamacoder.com/i/web/2025-08-18_12-32-49.jpg)

环境自动配置脚本执行中：

![image](https://file1.kamacoder.com/i/web/2025-08-18_12-33-20.jpg)

如果大家想进一步优化这个项目，这里给出可以深挖的点：

![image](https://file1.kamacoder.com/i/web/2025-08-18_12-33-47.jpg)



后端开发详细设计：包含 建表、日志库、架构、业务开发、群聊、后台管理、消息管理、管理员、文件上传下载、音视频通话等等：

![image](https://file1.kamacoder.com/i/web/2025-08-18_12-43-10.jpg)

![image](https://file1.kamacoder.com/i/web/2025-08-18_12-43-42.jpg)

前端开发详细设计：

![image](https://file1.kamacoder.com/i/web/2025-08-18_12-44-10.jpg)

![image](https://file1.kamacoder.com/i/web/2025-08-18_12-44-34.jpg)

## 本地开发与环境切换

当前仓库已经支持按 `base + env + 显式文件 + 环境变量覆盖` 的顺序加载配置：

1. `configs/config.toml`
2. `configs/config.<APP_ENV>.toml`
3. `CONFIG_FILE` 指定的额外配置文件
4. 运行时环境变量覆盖

### 本地开发

后端默认使用 `HTTP/WS`，适合先把主链路跑通：

```powershell
$env:APP_ENV = "dev"
go run .\cmd\gochat
```

前端默认读取 [web/chat-server/.env.development](web/chat-server/.env.development)，对应本地后端地址 `http://localhost:8000` 和 `ws://localhost:8000`：

```powershell
cd .\web\chat-server
npm install
npm run serve
```

### Docker Compose 一键启动

现在仓库里已经补了一套和 `CVAT` 类似的容器编排，可以一次性把前端、后端、`MySQL`、`Redis` 起在一起：

```powershell
docker compose up -d --build
```

默认会启动这些服务：

- `gochat-frontend`：前端页面，访问 `http://localhost:8080`
- `gochat-backend`：后端接口，访问 `http://localhost:8000`
- `gochat-mysql`：MySQL 8，映射到 `localhost:3306`
- `gochat-redis`：Redis 7，映射到 `localhost:6379`
- `gochat-seed`：一次性初始化测试用户，执行完成后自动退出

默认测试账号：

- 手机号：`13600000000`
- 密码：`123456`

容器编排默认使用 [configs/config.docker.toml](configs/config.docker.toml)。

如果你不想让前端构建时写死成 `localhost:8000`，可以在启动前覆盖：

```powershell
$env:FRONTEND_API_BASE_URL = "http://your-host:8000"
$env:FRONTEND_WS_BASE_URL = "ws://your-host:8000"
docker compose up -d --build
```

如果你希望把容器数据落到宿主机固定目录，而不是 Docker named volume，可以在仓库根目录放一个本地 `.env`，覆盖这些变量：

```powershell
GOCHAT_MYSQL_DATA=D:/develop/docker_workspace/gochat/mysql
GOCHAT_REDIS_DATA=D:/develop/docker_workspace/gochat/redis
GOCHAT_BACKEND_LOGS=D:/develop/docker_workspace/gochat/logs
GOCHAT_BACKEND_AVATARS=D:/develop/docker_workspace/gochat/avatars
GOCHAT_BACKEND_FILES=D:/develop/docker_workspace/gochat/files
```

仓库里附了示例文件 [`.env.compose.example`](</D:/develop/project/go-project/GoChat/.env.compose.example>)。这些变量不设置时，compose 会继续回退到原来的 named volume。

查看状态：

```powershell
docker compose ps
docker compose logs -f backend
docker compose logs -f frontend
```

停止并清理：

```powershell
docker compose down
```

如果连数据卷也一起清掉：

```powershell
docker compose down -v
```

当前 compose 默认仍然使用 `channel` 模式，因此 **不依赖 Kafka**。如果你后面想切到 Kafka，再单独补 broker 编排更合适。

### 生产部署

推荐让 `Nginx` 或其他反向代理终止 `TLS`，Go 服务继续监听内网 `HTTP`。可以参考 [configs/config.prod.example.toml](configs/config.prod.example.toml) 复制出自己的 `configs/config.prod.toml`，再通过环境变量注入密码、短信密钥和证书路径。

```powershell
$env:APP_ENV = "prod"
$env:GOCHAT_DB_PASSWORD = "replace-me"
go run .\cmd\gochat
```

前端生产构建默认读取 [web/chat-server/.env.production](web/chat-server/.env.production)。如果 `VUE_APP_API_BASE_URL` 和 `VUE_APP_WS_BASE_URL` 为空，会回退到当前页面所在的同源地址。

### 配置建议

- 放在配置文件里的内容：监听地址、端口、静态目录、日志路径、功能开关、数据库网络模式
- 放在环境变量里的内容：数据库密码、短信密钥、生产证书路径
- `WebRTC` 依赖安全上下文，本地如果要测音视频能力，请把前端 `DEV_SERVER_HTTPS=true`，并同时提供 `DEV_SERVER_HTTPS_CERT` 和 `DEV_SERVER_HTTPS_KEY`
- 本次环境切换改造没有处理短信注册/短信登录本身的可用性，短信相关配置为空时这部分能力仍需你自己按部署环境补齐

## 答疑

本项目在[知识星球](https://programmercarl.com/other/kstar.html)里为 文字专栏形式，大家不用担心，看不懂，星球里每个项目有专属答疑群，任何问题都可以在群里问，都会得到解答：

![](https://file1.kamacoder.com/i/web/2025-09-26_11-30-13.jpg)


## 获取本项目专栏

**本文档仅为星球内部专享，大家可以加入[知识星球](./kstar.md)里获取，在星球置顶一**

