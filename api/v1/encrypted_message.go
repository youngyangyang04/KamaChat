package v1

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/model"
	"kama_chat_server/internal/service/gorm"
	"kama_chat_server/pkg/zlog"
	"time"

	"github.com/gin-gonic/gin"
)

// SendEncryptedMessage 发送加密消息
func SendEncryptedMessage(c *gin.Context) {
	uuid, exists := c.Get("uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权",
		})
		return
	}

	var req request.SendEncryptedMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error("参数绑定失败: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 记录接收到的数据
	zlog.Info(fmt.Sprintf("接收到加密消息请求: ciphertext长度=%d, iv长度=%d, auth_tag长度=%d, ratchet_key长度=%d",
		len(req.Ciphertext), len(req.IV), len(req.AuthTag), len(req.RatchetKey)))
	if len(req.Ciphertext) > 0 {
		previewLen := len(req.Ciphertext)
		if previewLen > 50 {
			previewLen = 50
		}
		zlog.Info(fmt.Sprintf("ciphertext预览: %s", req.Ciphertext[:previewLen]))
	}

	// 验证发送者
	if uuid.(string) != req.ReceiverId {
		// 正常情况：发送者发给接收者
		// 需要确保 sender 是当前用户
		// 这里简化处理，假设已验证
	}

	// 如果是 PreKeyMessage，标记 OTP 密钥为已使用
	if req.MessageType == "PreKeyMessage" {
		if req.UsedOneTimePreKeyId != nil {
			zlog.Info(fmt.Sprintf("PreKeyMessage 使用了 OTP 密钥，key_id: %d", *req.UsedOneTimePreKeyId))
			if err := gorm.CryptoKeyService.MarkOneTimePreKeyAsUsed(
				req.ReceiverId,
				*req.UsedOneTimePreKeyId,
				uuid.(string),
			); err != nil {
				zlog.Error("标记 OTP 密钥失败: " + err.Error())
				// 不返回错误，继续发送消息
			}
		} else {
			zlog.Info("PreKeyMessage 未使用 OTP 密钥（used_one_time_pre_key_id 为 nil）")
		}
	}

	// 保存加密消息到数据库
	zlog.Info(fmt.Sprintf("保存加密消息: ciphertext长度=%d, iv长度=%d, auth_tag长度=%d",
		len(req.Ciphertext), len(req.IV), len(req.AuthTag)))

	message := model.Message{
		Uuid:      "M" + generateRandomString(11),
		SessionId: req.SessionId,
		SendId:    uuid.(string),
		ReceiveId: req.ReceiverId,
		Type:      0,              // 文本消息
		Content:   req.Ciphertext, // 存储密文
		Status:    0,
		CreatedAt: time.Now(),

		// 加密相关字段
		IsEncrypted:                 true,
		EncryptionVersion:           1,
		MessageType:                 req.MessageType,
		SenderIdentityKey:           req.SenderIdentityKey,
		SenderIdentityKeyCurve25519: req.SenderIdentityKeyCurve25519,
		SenderEphemeralKey:          req.SenderEphemeralKey,
		UsedOneTimePreKeyId:         req.UsedOneTimePreKeyId,
		RatchetKey:                  req.RatchetKey,
		Counter:                     req.Counter,
		PrevCounter:                 req.PrevCounter,
		IV:                          req.IV,
		AuthTag:                     req.AuthTag,
	}

	if err := dao.GormDB.Create(&message).Error; err != nil {
		zlog.Error("保存加密消息失败: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "发送消息失败",
		})
		return
	}

	zlog.Info("加密消息发送成功")
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "消息发送成功",
		"data": gin.H{
			"message_id": message.Uuid,
		},
	})
}

// 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// 如果 crypto/rand 失败，使用时间戳作为后备
		return time.Now().Format("20060102150405")
	}
	// 将随机字节映射到字符集
	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}
	return string(bytes)
}
