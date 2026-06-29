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
	GormDB *gorm.DB
	dbMu   sync.Mutex
)

func Init() error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if GormDB != nil {
		return nil
	}

	db, err := open(config.GetConfig())
	if err != nil {
		return err
	}

	if err := autoMigrate(db); err != nil {
		return err
	}

	GormDB = db
	return nil
}

func MustInit() {
	if err := Init(); err != nil {
		zlog.Fatal(err.Error())
	}
}

func BuildDSN(mysqlConfig config.MysqlConfig) string {
	base := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.DatabaseName,
	)

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
	dsn := BuildDSN(cfg.MysqlConfig)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.UserInfo{},
		&model.GroupInfo{},
		&model.UserContact{},
		&model.Session{},
		&model.ContactApply{},
		&model.Message{},
	)
}
