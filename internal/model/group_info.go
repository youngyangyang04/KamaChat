package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type GroupInfo struct {
	Id        int64           `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid      string          `gorm:"column:uuid;uniqueIndex;type:char(20);not null;comment:群组唯一id"`
	Name      string          `gorm:"column:name;type:varchar(20);not null;comment:群名称"`
	Notice    string          `gorm:"column:notice;type:varchar(500);comment:群公告"`
	Members   json.RawMessage `gorm:"column:members;type:json;comment:群组成员"`
	MemberCnt int             `gorm:"column:member_cnt;default:1;comment:群人数"` // 默认群主1人
	OwnerId   string          `gorm:"column:owner_id;type:char(20);not null;comment:群主uuid"`
	AddMode   bool            `gorm:"column:add_mode;type:tinyint(1);default:false;comment:加群方式，0.直接，1.审核"`
	Avatar    string          `gorm:"column:avatar;type:char(255);default:default_avatar.png;not null;comment:头像"`
	CreatedAt time.Time       `gorm:"column:created_at;index;type:datetime;not null;comment:创建时间"`
	UpdatedAt time.Time       `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at;index;comment:删除时间"`
}

func (GroupInfo) TableName() string {
	return "group_info"
}
