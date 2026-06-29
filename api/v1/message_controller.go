package v1

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/dto/request"
	"gochat/internal/service/gorm"
	"gochat/pkg/constants"
	"net/http"
)

// GetMessageList 获取聊天记录
func GetMessageList(c *gin.Context) {
	var req request.GetMessageListRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, rsp, ret := gorm.MessageService.GetMessageList(req.UserOneId, req.UserTwoId)
	JsonBack(c, message, ret, rsp)
}

// GetGroupMessageList 获取群聊消息记录
func GetGroupMessageList(c *gin.Context) {
	var req request.GetGroupMessageListRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, rsp, ret := gorm.MessageService.GetGroupMessageList(req.GroupId)
	JsonBack(c, message, ret, rsp)
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	message, ret := gorm.MessageService.UploadAvatar(c)
	JsonBack(c, message, ret, nil)
}

// UploadFile 上传头像
func UploadFile(c *gin.Context) {
	message, ret := gorm.MessageService.UploadFile(c)
	JsonBack(c, message, ret, nil)
}
