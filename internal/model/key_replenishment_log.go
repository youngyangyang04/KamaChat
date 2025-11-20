package model

import "time"

// KeyReplenishmentLog 密钥补充记录
type KeyReplenishmentLog struct {
	Id             int64     `gorm:"column:id;primaryKey;comment:自增id"`
	UserId         string    `gorm:"column:user_id;type:varchar(20);not null;index;comment:用户UUID"`
	ReplenishedAt  time.Time `gorm:"column:replenished_at;type:timestamp;not null;default:now();comment:补充时间"`
	KeysAdded      int       `gorm:"column:keys_added;not null;comment:本次添加的密钥数量"`
	RemainingKeys  int       `gorm:"column:remaining_keys;not null;comment:补充后剩余的密钥数量"`
}

func (KeyReplenishmentLog) TableName() string {
	return "key_replenishment_log"
}

