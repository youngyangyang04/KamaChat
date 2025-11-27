package request

// UploadPublicKeyBundleRequest 上传公钥束请求
type UploadPublicKeyBundleRequest struct {
	UserId                      string              `json:"user_id" binding:"required"`
	IdentityKeyPublic           string              `json:"identity_key_public" binding:"required"`
	IdentityKeyPublicCurve25519 string              `json:"identity_key_public_curve25519" binding:"required"`
	SignedPreKey                SignedPreKeyData    `json:"signed_pre_key" binding:"required"`
	OneTimePreKeys              []OneTimePreKeyData `json:"one_time_pre_keys" binding:"required,min=1"`
}

