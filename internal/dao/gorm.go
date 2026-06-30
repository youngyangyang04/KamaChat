package dao

import (
	"fmt"
	"strings"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gochat/internal/config"
	"gochat/internal/model"
	"gochat/pkg/zlog"
)

var (
	// GormDB 保存全局数据库连接，供 service / controller 复用。
	GormDB *gorm.DB
	// dbMu 保护 Init 只会在并发场景下真正执行一次初始化。
	dbMu sync.Mutex
)

func Init() error {
	dbMu.Lock()
	defer dbMu.Unlock()

	// 已经初始化过就直接返回，避免重复建连和重复迁移。
	if GormDB != nil {
		return nil
	}

	// 从全局配置中读取 MySQL 参数，组装 DSN 并打开连接。
	db, err := open(config.GetConfig())
	if err != nil {
		return err
	}

	// 启动时自动同步表结构，保证基础表存在。
	if err := autoMigrate(db); err != nil {
		return err
	}

	// 初始化成功后把连接保存到全局变量。
	GormDB = db
	return nil
}

func MustInit() {
	// 确保初始化成功,报错了直接把服务停掉
	if err := Init(); err != nil {
		zlog.Fatal(err.Error())
	}
}

func BuildDSN(mysqlConfig config.MysqlConfig) string {
	// 默认按 tcp 方式拼接 DSN。
	base := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.DatabaseName,
	)

	// 如果显式配置 network=unix，就切换成 Unix Socket 连接方式。
	if strings.EqualFold(mysqlConfig.Network, "unix") {
		base = fmt.Sprintf(
			"%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlConfig.User,
			mysqlConfig.Password,
			mysqlConfig.SocketPath,
			mysqlConfig.DatabaseName,
		)
	}

	return base
}

func open(cfg *config.Config) (*gorm.DB, error) {
	// 根据配置组装 DSN，再交给 gorm.Open 创建连接。
	dsn := BuildDSN(cfg.MysqlConfig)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func autoMigrate(db *gorm.DB) error {
	// 把项目需要的 model 一次性迁移到数据库。
	// 这里属于“启动即保证表结构可用”的策略。
	return db.AutoMigrate(
		&model.UserInfo{},
		&model.GroupInfo{},
		&model.UserContact{},
		&model.Session{},
		&model.ContactApply{},
		&model.Message{},
	)
}
