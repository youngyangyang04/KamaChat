package v1

import (
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/service/gorm"
	"kama_chat_server/pkg/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetNotificationList 获取通知列表
func GetNotificationList(c *gin.Context) {
	var req request.GetNotificationListRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, data, ret := gorm.NotificationService.GetNotificationList(
		req.UserId,
		req.Page,
		req.PageSize,
		req.Type,
		req.Status,
	)
	if ret == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": message,
			"data":    data.Data,
			"total":   data.Total,
		})
	} else if ret == -2 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": message,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": message,
		})
	}
}

// GetUnreadCount 获取未读通知数量
func GetUnreadCount(c *gin.Context) {
	var req request.GetUnreadCountRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, count, ret := gorm.NotificationService.GetUnreadCount(req.UserId, req.Type)
	JsonBack(c, message, ret, respond.GetUnreadCountResponse{Count: count})
}

// MarkAsRead 标记通知为已读
func MarkAsRead(c *gin.Context) {
	var req request.MarkAsReadRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, count, ret := gorm.NotificationService.MarkAsRead(req.UserId, req.NotificationIds)
	JsonBack(c, message, ret, respond.MarkAsReadResponse{Count: count})
}

// DeleteNotification 删除通知
func DeleteNotification(c *gin.Context) {
	var req request.DeleteNotificationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, count, ret := gorm.NotificationService.DeleteNotification(req.UserId, req.NotificationIds)
	JsonBack(c, message, ret, respond.DeleteNotificationResponse{Count: count})
}

// ClearAllNotification 清空所有通知
func ClearAllNotification(c *gin.Context) {
	var req request.ClearAllNotificationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, count, ret := gorm.NotificationService.ClearAllNotification(req.UserId, req.Type)
	JsonBack(c, message, ret, respond.ClearAllNotificationResponse{Count: count})
}

