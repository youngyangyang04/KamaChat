package model

import (
	"database/sql"
	"time"
)

type Message struct {
	Id         int64     `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid       string    `gorm:"column:uuid;uniqueIndex;type:char(20);not null;comment:消息uuid"`
	SessionId  string    `gorm:"column:session_id;index;type:char(20);not null;comment:会话uuid"`
	Type       int8      `gorm:"column:type;not null;comment:消息类型，0.文本，1.语音，2.文件，3.通话"` // 通话不用存消息内容或者url
	Content    string    `gorm:"column:content;type:TEXT;comment:消息内容"`
	Url        string    `gorm:"column:url;type:char(255);comment:消息url"`
	SendId     string    `gorm:"column:send_id;index;type:char(20);not null;comment:发送者uuid"`
	SendName   string    `gorm:"column:send_name;type:varchar(20);not null;comment:发送者昵称"`
	SendAvatar string    `gorm:"column:send_avatar;type:varchar(255);not null;comment:发送者头像"`
	ReceiveId  string    `gorm:"column:receive_id;index;type:char(20);not null;comment:接受者uuid"`
	FileType   string    `gorm:"column:file_type;type:char(10);comment:文件类型"`
	FileName   string    `gorm:"column:file_name;type:varchar(50);comment:文件名"`
	FileSize   string    `gorm:"column:file_size;type:char(20);comment:文件大小"`
	Status     int8      `gorm:"column:status;not null;comment:状态，0.未发送，1.已发送"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;comment:创建时间"`
	SendAt     sql.NullTime `gorm:"column:send_at;comment:发送时间"`
	AVdata     string    `gorm:"column:av_data;comment:通话传递数据"`
	
	// 加密相关字段
	IsEncrypted                bool   `gorm:"column:is_encrypted;default:false;comment:是否为加密消息"`
	EncryptionVersion          int    `gorm:"column:encryption_version;default:1;comment:加密版本号"`
	MessageType                string `gorm:"column:message_type;type:varchar(20);default:'SignalMessage';comment:消息类型：PreKeyMessage 或 SignalMessage"`
	SenderIdentityKey          string `gorm:"column:sender_identity_key;type:text;comment:发送者身份公钥（Ed25519，PreKeyMessage专用）"`
	SenderIdentityKeyCurve25519 string `gorm:"column:sender_identity_key_curve25519;type:text;comment:发送者身份公钥（Curve25519，用于ECDH，PreKeyMessage专用）"`
	SenderEphemeralKey         string `gorm:"column:sender_ephemeral_key;type:text;comment:发送者临时公钥（PreKeyMessage专用）"`
	UsedOneTimePreKeyId        *int   `gorm:"column:used_one_time_pre_key_id;comment:使用的一次性预密钥ID（PreKeyMessage专用）"`
	RatchetKey            string `gorm:"column:ratchet_key;type:text;comment:当前DH公钥（Base64编码）"`
	Counter               *int   `gorm:"column:counter;comment:消息计数器"`
	PrevCounter           *int   `gorm:"column:prev_counter;comment:前一个链的计数器"`
	IV                    string `gorm:"column:iv;type:text;comment:AES-GCM初始化向量（Base64编码）"`
	AuthTag               string `gorm:"column:auth_tag;type:text;comment:AES-GCM认证标签（Base64编码）"`
}

func (Message) TableName() string {
	return "message"
}
