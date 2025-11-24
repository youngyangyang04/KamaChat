package dao

import (
	"fmt"
	"kama_chat_server/internal/config"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func init() {
	conf := config.GetConfig()
	user := conf.PostgresqlConfig.User
	password := conf.PostgresqlConfig.Password
	host := conf.PostgresqlConfig.Host
	port := conf.PostgresqlConfig.Port
	dbName := conf.PostgresqlConfig.DatabaseName

	// PostgreSQL DSN 格式
	var dsn string
	if password == "" {
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			host, user, dbName, port)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			host, user, password, dbName, port)
	}

	var err error
	GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zlog.Fatal(err.Error())
	}

	// 自动迁移，如果没有建表，会自动创建对应的表
	err = GormDB.AutoMigrate(
		&model.UserInfo{},
		&model.GroupInfo{},
		&model.UserContact{},
		&model.Session{},
		&model.ContactApply{},
		&model.Message{},
		&model.OneTimePreKey{},
		&model.KeyReplenishmentLog{},
		&model.Notification{},
	)
	if err != nil {
		zlog.Fatal(err.Error())
	}
}
