package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	Id            int64          `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid          string         `gorm:"column:uuid;uniqueIndex;type:char(20);comment:会话uuid"`
	SendId        string         `gorm:"column:send_id;Index;type:char(20);not null;comment:创建会话人id"`
	ReceiveId     string         `gorm:"column:receive_id;Index;type:char(20);not null;comment:接受会话人id"`
	ReceiveName   string         `gorm:"column:receive_name;type:varchar(20);not null;comment:名称"`
	Avatar        string         `gorm:"column:avatar;type:char(255);default:default_avatar.png;not null;comment:头像"`
	LastMessage   string         `gorm:"column:last_message;type:TEXT;comment:最新的消息"`
	LastMessageAt sql.NullTime   `gorm:"column:last_message_at;type:datetime;comment:最近接收时间"`
	CreatedAt     time.Time      `gorm:"column:created_at;Index;type:datetime;comment:创建时间"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;Index;type:datetime;comment:删除时间"`
}

func (Session) TableName() string {
	return "session"
}
