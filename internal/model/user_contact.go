package model

import (
	"gorm.io/gorm"
	"time"
)

type UserContact struct {
	Id          int64          `gorm:"column:id;primaryKey;comment:自增id"`
	UserId      string         `gorm:"column:user_id;index;type:char(20);not null;comment:用户唯一id"`
	ContactId   string         `gorm:"column:contact_id;index;type:char(20);not null;comment:对应联系id"`
	ContactType int8           `gorm:"column:contact_type;not null;comment:联系类型，0.用户，1.群聊"`
	Status      int8           `gorm:"column:status;not null;comment:联系状态，0.正常，1.拉黑，2.被拉黑，3.删除好友，4.被删除好友，5.被禁言，6.退出群聊，7.被踢出群聊"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp;not null;comment:创建时间"`
	UpdateAt    time.Time      `gorm:"column:update_at;type:timestamp;not null;comment:更新时间"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index;comment:删除时间"`
}

func (UserContact) TableName() string {
	return "user_contact"
}
