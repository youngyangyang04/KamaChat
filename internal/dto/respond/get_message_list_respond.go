package respond

type GetMessageListRespond struct {
	Uuid       string `json:"uuid"` // 消息 UUID（用于标识消息）
	SendId     string `json:"send_id"`
	SendName   string `json:"send_name"`
	SendAvatar string `json:"send_avatar"`
	ReceiveId  string `json:"receive_id"`
	Type       int8   `json:"type"`
	Content    string `json:"content"`
	Url        string `json:"url"`
	FileType   string `json:"file_type"`
	FileName   string `json:"file_name"`
	FileSize   string `json:"file_size"`
	CreatedAt  string `json:"created_at"` // 先用CreatedAt排序，后面考虑改成SentAt
	
	// 加密相关字段
	IsEncrypted                bool   `json:"is_encrypted"`
	EncryptionVersion          int    `json:"encryption_version,omitempty"`
	MessageType                string `json:"message_type,omitempty"`
	SenderIdentityKey          string `json:"sender_identity_key,omitempty"`
	SenderIdentityKeyCurve25519 string `json:"sender_identity_key_curve25519,omitempty"`
	SenderEphemeralKey         string `json:"sender_ephemeral_key,omitempty"`
	UsedOneTimePreKeyId        *int   `json:"used_one_time_pre_key_id,omitempty"`
	RatchetKey            string `json:"ratchet_key,omitempty"`
	Counter               *int   `json:"counter,omitempty"`
	PrevCounter           *int   `json:"prev_counter,omitempty"`
	IV                    string `json:"iv,omitempty"`
	AuthTag               string `json:"auth_tag,omitempty"`
}
