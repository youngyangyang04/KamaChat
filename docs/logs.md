## 修改日志

### 2025-11-18 - 数据库从 MySQL 迁移到 PostgreSQL
- 修改 `configs/config.toml`：配置改为 PostgreSQL（端口 5432，用户 postgres，空密码）
- 修改 `internal/config/config.go`：`MysqlConfig` → `PostgresqlConfig`，配置加载路径改为本地
- 修改 `internal/dao/gorm.go`：数据库驱动改为 `gorm.io/driver/postgres`，更新 DSN 连接格式
- 添加依赖：`gorm.io/driver/postgres v1.6.0` 及相关 PostgreSQL 驱动包
- Git：创建 dev 分支并推送 (commit: 936c849)

---