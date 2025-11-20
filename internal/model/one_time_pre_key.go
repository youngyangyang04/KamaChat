package model

import (
	"database/sql"
	"time"
)

// OneTimePreKey 一次性预密钥
type OneTimePreKey struct {
	Id            int64        `gorm:"column:id;primaryKey;comment:自增id"`
	UserId        string       `gorm:"column:user_id;type:varchar(20);not null;index;comment:用户UUID"`
	KeyId         int          `gorm:"column:key_id;not null;comment:密钥ID"`
	PublicKey     string       `gorm:"column:public_key;type:text;not null;comment:公钥（Curve25519，Base64编码）"`
	UploadedAt    time.Time    `gorm:"column:uploaded_at;type:timestamp;not null;default:now();comment:上传时间"`
	Used          bool         `gorm:"column:used;default:false;comment:是否已被使用"`
	UsedAt        sql.NullTime `gorm:"column:used_at;type:timestamp;comment:使用时间"`
	UsedByUserId  string       `gorm:"column:used_by_user_id;type:varchar(20);comment:被谁使用"`
}

func (OneTimePreKey) TableName() string {
	return "one_time_pre_keys"
}

