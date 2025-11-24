package gorm

import (
	"fmt"
	"io"
	"kama_chat_server/internal/config"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type messageService struct {
}

var MessageService = new(messageService)

// GetMessageList 获取聊天记录
func (m *messageService) GetMessageList(userOneId, userTwoId string) (string, []respond.GetMessageListRespond, int) {
	zlog.Info(fmt.Sprintf("%s %s", userOneId, userTwoId))
	var messageList []model.Message
	if res := dao.GormDB.Where("(send_id = ? AND receive_id = ?) OR (send_id = ? AND receive_id = ?)", userOneId, userTwoId, userTwoId, userOneId).Order("created_at ASC").Find(&messageList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 收集所有唯一的 send_id，批量查询用户信息
	sendIdSet := make(map[string]bool)
	for _, message := range messageList {
		if message.SendId != "" {
			sendIdSet[message.SendId] = true
		}
	}

	// 批量查询用户信息
	var sendIdList []string
	for sendId := range sendIdSet {
		sendIdList = append(sendIdList, sendId)
	}

	userInfoMap := make(map[string]model.UserInfo)
	if len(sendIdList) > 0 {
		var users []model.UserInfo
		if res := dao.GormDB.Select("uuid, nickname, avatar").Where("uuid IN ?", sendIdList).Find(&users); res.Error != nil {
			zlog.Error("批量查询用户信息失败: " + res.Error.Error())
		} else {
			for _, user := range users {
				userInfoMap[user.Uuid] = user
			}
		}
	}

	var rspList []respond.GetMessageListRespond
	for _, message := range messageList {
		if message.IsEncrypted {
			zlog.Info(fmt.Sprintf("读取加密消息: Content长度=%d, IV长度=%d, AuthTag长度=%d",
				len(message.Content), len(message.IV), len(message.AuthTag)))
		}

		// 从用户信息 map 中获取发送者的名字和头像
		sendName := ""
		sendAvatar := ""
		if user, exists := userInfoMap[message.SendId]; exists {
			sendName = user.Nickname
			sendAvatar = user.Avatar
		}

		rspList = append(rspList, respond.GetMessageListRespond{
			Uuid:       message.Uuid, // 消息 UUID
			SendId:     message.SendId,
			SendName:   sendName,   // 从用户表动态查询
			SendAvatar: sendAvatar, // 从用户表动态查询
			ReceiveId:  message.ReceiveId,
			Content:    message.Content,
			Url:        message.Url,
			Type:       message.Type,
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
		})
	}
	return "获取聊天记录成功", rspList, 0
}

// GetGroupMessageList 获取群聊消息记录
func (m *messageService) GetGroupMessageList(groupId string) (string, []respond.GetGroupMessageListRespond, int) {
	var messageList []model.Message
	if res := dao.GormDB.Where("receive_id = ?", groupId).Order("created_at ASC").Find(&messageList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 收集所有唯一的 send_id，批量查询用户信息
	sendIdSet := make(map[string]bool)
	for _, message := range messageList {
		if message.SendId != "" {
			sendIdSet[message.SendId] = true
		}
	}

	// 批量查询用户信息
	var sendIdList []string
	for sendId := range sendIdSet {
		sendIdList = append(sendIdList, sendId)
	}

	userInfoMap := make(map[string]model.UserInfo)
	if len(sendIdList) > 0 {
		var users []model.UserInfo
		if res := dao.GormDB.Select("uuid, nickname, avatar").Where("uuid IN ?", sendIdList).Find(&users); res.Error != nil {
			zlog.Error("批量查询用户信息失败: " + res.Error.Error())
		} else {
			for _, user := range users {
				userInfoMap[user.Uuid] = user
			}
		}
	}

	var rspList []respond.GetGroupMessageListRespond
	for _, message := range messageList {
		// 从用户信息 map 中获取发送者的名字和头像
		sendName := ""
		sendAvatar := ""
		if user, exists := userInfoMap[message.SendId]; exists {
			sendName = user.Nickname
			sendAvatar = user.Avatar
		}

		rsp := respond.GetGroupMessageListRespond{
			SendId:     message.SendId,
			SendName:   sendName,   // 从用户表动态查询
			SendAvatar: sendAvatar, // 从用户表动态查询
			ReceiveId:  message.ReceiveId,
			Content:    message.Content,
			Url:        message.Url,
			Type:       message.Type,
			FileType:   message.FileType,
			FileName:   message.FileName,
			FileSize:   message.FileSize,
			CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rspList = append(rspList, rsp)
	}
	return "获取聊天记录成功", rspList, 0
}

// UploadAvatar 上传头像
func (m *messageService) UploadAvatar(c *gin.Context) (string, int) {
	if err := c.Request.ParseMultipartForm(constants.FILE_MAX_SIZE); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	mForm := c.Request.MultipartForm
	for key, _ := range mForm.File {
		file, fileHeader, err := c.Request.FormFile(key)
		if err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}
		defer file.Close()
		zlog.Info(fmt.Sprintf("文件名：%s，文件大小：%d", fileHeader.Filename, fileHeader.Size))
		// 原来Filename应该是213451545.xxx，将Filename修改为avatar_ownerId.xxx
		ext := filepath.Ext(fileHeader.Filename)
		zlog.Info(ext)
		localFileName := config.GetConfig().StaticAvatarPath + "/" + fileHeader.Filename
		out, err := os.Create(localFileName)
		if err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}
		zlog.Info("完成文件上传")
	}
	return "上传成功", 0
}

// UploadFile 上传文件
func (m *messageService) UploadFile(c *gin.Context) (string, int) {
	if err := c.Request.ParseMultipartForm(constants.FILE_MAX_SIZE); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	mForm := c.Request.MultipartForm
	for key, _ := range mForm.File {
		file, fileHeader, err := c.Request.FormFile(key)
		if err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}
		defer file.Close()
		zlog.Info(fmt.Sprintf("文件名：%s，文件大小：%d", fileHeader.Filename, fileHeader.Size))
		// 原来Filename应该是213451545.xxx，将Filename修改为avatar_ownerId.xxx
		ext := filepath.Ext(fileHeader.Filename)
		zlog.Info(ext)
		localFileName := config.GetConfig().StaticFilePath + "/" + fileHeader.Filename
		out, err := os.Create(localFileName)
		if err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}
		zlog.Info("完成文件上传")
	}
	return "上传成功", 0
}
