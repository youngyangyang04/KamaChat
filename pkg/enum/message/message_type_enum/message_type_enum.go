package message_type_enum

const (
	Text = iota
	// 语音
	Voice
	// 文件
	File
	// 通话
	AudioOrVideo
	// 加密图片
	EncryptedImage
	// 加密文件
	EncryptedFile
)
