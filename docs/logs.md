## 修改日志

### 2025-11-18 - 数据库从 MySQL 迁移到 PostgreSQL
- 修改 `configs/config.toml`：配置改为 PostgreSQL（端口 5432，用户 postgres，空密码）
- 修改 `internal/config/config.go`：`MysqlConfig` → `PostgresqlConfig`，配置加载路径改为本地
- 修改 `internal/dao/gorm.go`：数据库驱动改为 `gorm.io/driver/postgres`，更新 DSN 连接格式
- 修改所有 model 文件：将时间字段类型从 `type:datetime` 改为 `type:timestamp`（PostgreSQL 不支持 datetime）
  - `internal/model/user_info.go`
  - `internal/model/user_contact.go`
  - `internal/model/session.go`
  - `internal/model/group_info.go`
  - `internal/model/contact_apply.go`
- 添加依赖：`gorm.io/driver/postgres v1.6.0` 及相关 PostgreSQL 驱动包

---