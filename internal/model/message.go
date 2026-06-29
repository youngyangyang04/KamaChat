package model

import (
	"database/sql"
	"time"
)

type Message struct {
	Id         int64        `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid       string       `gorm:"column:uuid;uniqueIndex;type:char(20);not null;comment:消息uuid"`
	SessionId  string       `gorm:"column:session_id;index;type:char(20);not null;comment:会话uuid"`
	Type       int8         `gorm:"column:type;not null;comment:消息类型，0.文本，1.语音，2.文件，3.通话"` // 通话不用存消息内容或者url
	Content    string       `gorm:"column:content;type:TEXT;comment:消息内容"`
	Url        string       `gorm:"column:url;type:char(255);comment:消息url"`
	SendId     string       `gorm:"column:send_id;index;type:char(20);not null;comment:发送者uuid"`
	SendName   string       `gorm:"column:send_name;type:varchar(20);not null;comment:发送者昵称"`
	SendAvatar string       `gorm:"column:send_avatar;type:varchar(255);not null;comment:发送者头像"`
	ReceiveId  string       `gorm:"column:receive_id;index;type:char(20);not null;comment:接受者uuid"`
	FileType   string       `gorm:"column:file_type;type:char(10);comment:文件类型"`
	FileName   string       `gorm:"column:file_name;type:varchar(50);comment:文件名"`
	FileSize   string       `gorm:"column:file_size;type:char(20);comment:文件大小"`
	Status     int8         `gorm:"column:status;not null;comment:状态，0.未发送，1.已发送"`
	CreatedAt  time.Time    `gorm:"column:created_at;not null;comment:创建时间"`
	SendAt     sql.NullTime `gorm:"column:send_at;comment:发送时间"`
	AVdata     string       `gorm:"column:av_data;comment:通话传递数据"`
}

func (Message) TableName() string {
	return "message"
}
