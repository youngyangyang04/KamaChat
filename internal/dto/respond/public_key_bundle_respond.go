package respond

// PublicKeyBundleRespond 公钥束响应
type PublicKeyBundleRespond struct {
	UserId                  string           `json:"user_id"`
	IdentityKey             string           `json:"identity_key"` // Base64, Ed25519
	IdentityKeyCurve25519   string           `json:"identity_key_curve25519"` // Base64, Curve25519
	SignedPreKey            SignedPreKeyInfo `json:"signed_pre_key"`
	OneTimePreKey           *OneTimePreKeyInfo `json:"one_time_pre_key"` // 可能为 null
}

// SignedPreKeyInfo 签名预密钥信息
type SignedPreKeyInfo struct {
	KeyId     int    `json:"key_id"`
	PublicKey string `json:"public_key"` // Base64
	Signature string `json:"signature"`  // Base64
}

// OneTimePreKeyInfo 一次性预密钥信息
type OneTimePreKeyInfo struct {
	KeyId     int    `json:"key_id"`
	PublicKey string `json:"public_key"` // Base64
}

