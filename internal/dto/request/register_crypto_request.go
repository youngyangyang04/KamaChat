package request

// RegisterCryptoRequest 注册时的加密密钥数据
type RegisterCryptoRequest struct {
	Account                   string           `json:"account" binding:"required,min=3,max=20"`
	Nickname                  string           `json:"nickname" binding:"required,min=1,max=20"`
	Password                  string           `json:"password" binding:"required,min=6,max=50"`
	IdentityKeyPublic         string           `json:"identity_key_public" binding:"required"` // Base64, Ed25519
	IdentityKeyPublicCurve25519 string         `json:"identity_key_public_curve25519" binding:"required"` // Base64, Curve25519
	SignedPreKey              SignedPreKeyData `json:"signed_pre_key" binding:"required"`
	OneTimePreKeys            []OneTimePreKeyData `json:"one_time_pre_keys" binding:"required,min=1"`
}

// SignedPreKeyData 签名预密钥数据
type SignedPreKeyData struct {
	KeyId     int    `json:"key_id" binding:"required"`
	PublicKey string `json:"public_key" binding:"required"` // Base64
	Signature string `json:"signature" binding:"required"`  // Base64
}

// OneTimePreKeyData 一次性预密钥数据
type OneTimePreKeyData struct {
	KeyId     int    `json:"key_id" binding:"required"`
	PublicKey string `json:"public_key" binding:"required"` // Base64
}

