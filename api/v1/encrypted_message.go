package v1

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"kama_chat_server/internal/config"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/internal/service/chat"
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

	// 确定消息类型：默认为文本(0)，如果指定了 FileMessageType 则使用指定的值
	msgType := int8(0)
	if req.FileMessageType != nil {
		msgType = int8(*req.FileMessageType)
		zlog.Info(fmt.Sprintf("加密文件消息类型: %d", msgType))
	}

	message := model.Message{
		Uuid:      "M" + generateRandomString(11),
		SessionId: req.SessionId,
		SendId:    uuid.(string),
		ReceiveId: req.ReceiverId,
		Type:      msgType,        // 消息类型: 0=文本, 4=加密图片, 5=加密文件
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

	// 通过 WebSocket 推送消息给发送方和接收方
	pushMessageViaWebSocket(&message)

	zlog.Info("加密消息发送成功")
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "消息发送成功",
		"data": gin.H{
			"message_id": message.Uuid,
		},
	})
}

// pushMessageViaWebSocket 通过 WebSocket 推送消息
func pushMessageViaWebSocket(message *model.Message) {
	// 获取用户信息（昵称和头像）
	var sender model.UserInfo
	if res := dao.GormDB.Select("nickname, avatar").Where("uuid = ?", message.SendId).First(&sender); res.Error != nil {
		zlog.Error("查询发送者信息失败: " + res.Error.Error())
		return
	}

	// 构建消息响应结构
	messageRsp := respond.GetMessageListRespond{
		Uuid:       message.Uuid,
		SendId:     message.SendId,
		SendName:   sender.Nickname,
		SendAvatar: sender.Avatar,
		ReceiveId:  message.ReceiveId,
		Type:       message.Type,
		Content:    message.Content, // 密文
		Url:        message.Url,
		FileType:   message.FileType,
		FileName:   message.FileName,
		FileSize:   message.FileSize,
		CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),

		// 加密相关字段
		IsEncrypted:                 message.IsEncrypted,
		EncryptionVersion:           message.EncryptionVersion,
		MessageType:                 message.MessageType,
		SenderIdentityKey:           message.SenderIdentityKey,
		SenderIdentityKeyCurve25519: message.SenderIdentityKeyCurve25519,
		SenderEphemeralKey:          message.SenderEphemeralKey,
		UsedOneTimePreKeyId:         message.UsedOneTimePreKeyId,
		RatchetKey:                  message.RatchetKey,
		Counter:                     message.Counter,
		PrevCounter:                 message.PrevCounter,
		IV:                          message.IV,
		AuthTag:                     message.AuthTag,
	}

	jsonMessage, err := json.Marshal(messageRsp)
	if err != nil {
		zlog.Error("序列化消息失败: " + err.Error())
		return
	}

	var messageBack = &chat.MessageBack{
		Message: jsonMessage,
		Uuid:    message.Uuid,
	}

	// 根据配置判断使用 channel 模式还是 kafka 模式
	kafkaConfig := config.GetConfig().KafkaConfig
	if kafkaConfig.MessageMode == "channel" {
		// 使用 channel 模式
		chat.ChatServer.PushMessage(messageBack, message.SendId, message.ReceiveId)
	} else {
		// 使用 kafka 模式
		chat.KafkaChatServer.PushMessage(messageBack, message.SendId, message.ReceiveId)
	}
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
