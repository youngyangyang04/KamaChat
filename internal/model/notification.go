package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	Id          int64          `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid        string         `gorm:"column:uuid;uniqueIndex;type:char(20);not null;comment:通知uuid"`
	UserId      string         `gorm:"column:user_id;index;type:char(20);not null;comment:接收者用户id"`
	Type        int8           `gorm:"column:type;not null;comment:通知类型，1.好友申请，2.群聊邀请，3.新消息，4.系统通知，5.群聊申请，6.好友申请已通过，7.好友申请已拒绝，8.群聊申请已通过，9.群聊申请已拒绝"`
	Title       string         `gorm:"column:title;type:varchar(100);not null;comment:通知标题"`
	Content     string         `gorm:"column:content;type:TEXT;comment:通知内容"`
	RelatedId   string         `gorm:"column:related_id;index;type:char(20);comment:关联id（如申请id、消息id等）"`
	RelatedType string         `gorm:"column:related_type;type:varchar(20);comment:关联类型（如contact_apply、message等）"`
	SenderId    string         `gorm:"column:sender_id;index;type:char(20);comment:发送者id（如果是系统通知则为空）"`
	SenderName  string         `gorm:"column:sender_name;type:varchar(20);comment:发送者昵称"`
	SenderAvatar string        `gorm:"column:sender_avatar;type:varchar(255);comment:发送者头像"`
	Status      int8           `gorm:"column:status;not null;default:0;comment:通知状态，0.未读，1.已读，2.已删除"`
	ExtraData   string         `gorm:"column:extra_data;type:TEXT;comment:额外数据（JSON格式，用于存储扩展信息）"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;comment:创建时间"`
	ReadAt      sql.NullTime   `gorm:"column:read_at;comment:已读时间"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;type:timestamp;comment:删除时间"`
}

func (Notification) TableName() string {
	return "notification"
}

