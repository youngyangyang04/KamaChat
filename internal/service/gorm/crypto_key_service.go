package gorm

import (
	"errors"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"
	"time"

	"gorm.io/gorm"
)

type cryptoKeyService struct{}

var CryptoKeyService = new(cryptoKeyService)

// SaveUserPublicKeys 保存用户公钥（注册时调用）
// 如果 db 为 nil，则使用默认的 dao.GormDB
func (c *cryptoKeyService) SaveUserPublicKeys(userId string, req *request.RegisterCryptoRequest, db ...*gorm.DB) error {
	var gormDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		gormDB = db[0] // 使用传入的 DB（可能是事务实例）
	} else {
		gormDB = dao.GormDB // 使用默认的 DB
	}

	// 1. 更新 user_info 表的公钥字段
	updates := map[string]interface{}{
		"identity_key_public":            req.IdentityKeyPublic,
		"identity_key_public_curve25519": req.IdentityKeyPublicCurve25519,
		"signed_pre_key_id":              req.SignedPreKey.KeyId,
		"signed_pre_key_public":          req.SignedPreKey.PublicKey,
		"signed_pre_key_signature":       req.SignedPreKey.Signature,
		"signed_pre_key_timestamp":       time.Now(),
		"key_generation":                 1,
	}

	if err := gormDB.Model(&model.UserInfo{}).Where("uuid = ?", userId).Updates(updates).Error; err != nil {
		zlog.Error("更新用户公钥失败: " + err.Error())
		return err
	}

	// 2. 批量插入一次性预密钥
	oneTimePreKeys := make([]model.OneTimePreKey, 0, len(req.OneTimePreKeys))
	now := time.Now()
	for _, key := range req.OneTimePreKeys {
		oneTimePreKeys = append(oneTimePreKeys, model.OneTimePreKey{
			UserId:     userId,
			KeyId:      key.KeyId,
			PublicKey:  key.PublicKey,
			UploadedAt: now,
			Used:       false,
		})
	}

	if err := gormDB.Create(&oneTimePreKeys).Error; err != nil {
		zlog.Error("保存一次性预密钥失败: " + err.Error())
		return err
	}

	zlog.Info("保存用户公钥成功, userId: " + userId)
	return nil
}

// GetPublicKeyBundle 获取用户的公钥束（用于建立会话）
func (c *cryptoKeyService) GetPublicKeyBundle(userId string) (*respond.PublicKeyBundleRespond, error) {
	// 1. 获取用户信息（包含身份公钥和签名预公钥）
	var userInfo model.UserInfo
	if err := dao.GormDB.Where("uuid = ?", userId).First(&userInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		zlog.Error("查询用户信息失败: " + err.Error())
		return nil, err
	}

	// 检查是否有加密密钥
	if userInfo.IdentityKeyPublic == "" {
		return nil, errors.New("该用户未启用加密功能")
	}

	// 2. 构造响应
	bundle := &respond.PublicKeyBundleRespond{
		UserId:                userId,
		IdentityKey:           userInfo.IdentityKeyPublic,
		IdentityKeyCurve25519: userInfo.IdentityKeyPublicCurve25519,
		SignedPreKey: respond.SignedPreKeyInfo{
			KeyId:     userInfo.SignedPreKeyId,
			PublicKey: userInfo.SignedPreKeyPublic,
			Signature: userInfo.SignedPreKeySignature,
		},
	}

	// 3. 获取一个未使用的一次性预密钥（可选）
	var oneTimePreKey model.OneTimePreKey
	err := dao.GormDB.Where("user_id = ? AND used = ?", userId, false).
		Order("id ASC").
		First(&oneTimePreKey).Error

	if err == nil {
		// 找到了未使用的 OTP 密钥
		bundle.OneTimePreKey = &respond.OneTimePreKeyInfo{
			KeyId:     oneTimePreKey.KeyId,
			PublicKey: oneTimePreKey.PublicKey,
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 其他错误
		zlog.Error("查询一次性预密钥失败: " + err.Error())
		return nil, err
	}
	// 如果是 ErrRecordNotFound，oneTimePreKey 为 nil，表示没有可用的 OTP 密钥

	zlog.Info("获取公钥束成功, userId: " + userId)
	return bundle, nil
}

// MarkOneTimePreKeyAsUsed 标记一次性预密钥为已使用
func (c *cryptoKeyService) MarkOneTimePreKeyAsUsed(userId string, keyId int, usedByUserId string) error {
	updates := map[string]interface{}{
		"used":            true,
		"used_at":         time.Now(),
		"used_by_user_id": usedByUserId,
	}

	result := dao.GormDB.Model(&model.OneTimePreKey{}).
		Where("user_id = ? AND key_id = ?", userId, keyId).
		Updates(updates)

	if result.Error != nil {
		zlog.Error("标记一次性预密钥为已使用失败: " + result.Error.Error())
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("未找到指定的一次性预密钥")
	}

	zlog.Info("标记一次性预密钥为已使用成功")
	return nil
}

// GetOneTimePreKeyCount 获取剩余的一次性预密钥数量
func (c *cryptoKeyService) GetOneTimePreKeyCount(userId string) (int64, error) {
	var count int64
	err := dao.GormDB.Model(&model.OneTimePreKey{}).
		Where("user_id = ? AND used = ?", userId, false).
		Count(&count).Error

	if err != nil {
		zlog.Error("查询一次性预密钥数量失败: " + err.Error())
		return 0, err
	}

	return count, nil
}

// ReplenishOneTimePreKeys 补充一次性预密钥
func (c *cryptoKeyService) ReplenishOneTimePreKeys(userId string, keys []request.OneTimePreKeyData) error {
	// 1. 批量插入新密钥
	oneTimePreKeys := make([]model.OneTimePreKey, 0, len(keys))
	now := time.Now()
	for _, key := range keys {
		oneTimePreKeys = append(oneTimePreKeys, model.OneTimePreKey{
			UserId:     userId,
			KeyId:      key.KeyId,
			PublicKey:  key.PublicKey,
			UploadedAt: now,
			Used:       false,
		})
	}

	if err := dao.GormDB.Create(&oneTimePreKeys).Error; err != nil {
		zlog.Error("补充一次性预密钥失败: " + err.Error())
		return err
	}

	// 2. 查询剩余数量
	remainingCount, err := c.GetOneTimePreKeyCount(userId)
	if err != nil {
		return err
	}

	// 3. 记录日志
	log := model.KeyReplenishmentLog{
		UserId:        userId,
		ReplenishedAt: now,
		KeysAdded:     len(keys),
		RemainingKeys: int(remainingCount),
	}

	if err := dao.GormDB.Create(&log).Error; err != nil {
		zlog.Error("记录密钥补充日志失败: " + err.Error())
		// 不返回错误，日志记录失败不影响主流程
	}

	zlog.Info("补充一次性预密钥成功")
	return nil
}

// RotateSignedPreKey 轮换签名预密钥
func (c *cryptoKeyService) RotateSignedPreKey(userId string, newKey request.SignedPreKeyData) error {
	updates := map[string]interface{}{
		"signed_pre_key_id":        newKey.KeyId,
		"signed_pre_key_public":    newKey.PublicKey,
		"signed_pre_key_signature": newKey.Signature,
		"signed_pre_key_timestamp": time.Now(),
		"key_generation":           gorm.Expr("key_generation + 1"),
	}

	if err := dao.GormDB.Model(&model.UserInfo{}).Where("uuid = ?", userId).Updates(updates).Error; err != nil {
		zlog.Error("轮换签名预密钥失败: " + err.Error())
		return err
	}

	zlog.Info("轮换签名预密钥成功, userId: " + userId)
	return nil
}
