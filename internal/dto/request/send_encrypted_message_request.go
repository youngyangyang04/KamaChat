package request

// SendEncryptedMessageRequest 发送加密消息请求
type SendEncryptedMessageRequest struct {
	SessionId       string `json:"session_id" binding:"required"`
	ReceiverId      string `json:"receiver_id" binding:"required"`
	MessageType     string `json:"message_type" binding:"required"` // "PreKeyMessage" 或 "SignalMessage"
	FileMessageType *int   `json:"file_message_type,omitempty"`     // 可选：4=加密图片, 5=加密文件

	// PreKeyMessage 专用字段（首次会话）
	SenderIdentityKey           string `json:"sender_identity_key,omitempty"`
	SenderIdentityKeyCurve25519 string `json:"sender_identity_key_curve25519,omitempty"` // Curve25519 格式的身份公钥（用于 ECDH）
	SenderEphemeralKey          string `json:"sender_ephemeral_key,omitempty"`
	UsedOneTimePreKeyId         *int   `json:"used_one_time_pre_key_id,omitempty"`

	// SignalMessage 通用字段（对于 PreKeyMessage，这些字段也是必需的）
	RatchetKey  string `json:"ratchet_key" binding:"required"`
	Counter     *int   `json:"counter"`      // 使用指针类型，允许为空（PreKeyMessage 时可能不需要）
	PrevCounter *int   `json:"prev_counter"` // 使用指针类型，允许为空

	// AES-GCM 加密数据
	Ciphertext string `json:"ciphertext" binding:"required"` // Base64
	IV         string `json:"iv" binding:"required"`         // Base64
	AuthTag    string `json:"auth_tag" binding:"required"`   // Base64
}
