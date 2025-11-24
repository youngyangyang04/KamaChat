package gorm

import (
	"database/sql"
	"fmt"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/constants"
	notification_status_enum "kama_chat_server/pkg/enum/notification/notification_status_enum"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
	"time"

	"gorm.io/gorm"
)

type notificationService struct {
}

var NotificationService = new(notificationService)

// CreateNotification 创建通知
// message: 返回给前端的提示
// data: 返回的通知对象
// ret: 0.成功，-1.系统错误，-2.业务错误
func (n *notificationService) CreateNotification(userId string, notificationType int8, title string, content string, relatedId string, relatedType string, senderId string, senderName string, senderAvatar string, extraData string) (message string, data *model.Notification, ret int) {
	notification := model.Notification{
		Uuid:        fmt.Sprintf("N%s", random.GetNowAndLenRandomString(11)),
		UserId:      userId,
		Type:        notificationType,
		Title:       title,
		Content:     content,
		RelatedId:   relatedId,
		RelatedType: relatedType,
		SenderId:    senderId,
		SenderName:  senderName,
		SenderAvatar: senderAvatar,
		Status:      notification_status_enum.Unread, // 默认未读
		ExtraData:   extraData,
		CreatedAt:   time.Now(),
	}

	if res := dao.GormDB.Create(&notification); res.Error != nil {
		zlog.Error("创建通知失败: " + res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	zlog.Info(fmt.Sprintf("创建通知成功: uuid=%s, user_id=%s, type=%d", notification.Uuid, userId, notificationType))
	return "通知创建成功", &notification, 0
}

// GetNotificationList 获取通知列表
func (n *notificationService) GetNotificationList(userId string, page int, pageSize int, notificationType *int8, status *int8) (message string, data respond.GetNotificationListResponse, ret int) {
	// 默认值处理
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大每页数量
	}

	// 构建查询
	query := dao.GormDB.Model(&model.Notification{}).Where("user_id = ?", userId)

	// 类型筛选
	if notificationType != nil {
		query = query.Where("type = ?", *notificationType)
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		zlog.Error("获取通知总数失败: " + err.Error())
		return constants.SYSTEM_ERROR, respond.GetNotificationListResponse{}, -1
	}

	// 分页查询
	var notifications []model.Notification
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error; err != nil {
		zlog.Error("获取通知列表失败: " + err.Error())
		return constants.SYSTEM_ERROR, respond.GetNotificationListResponse{}, -1
	}

	// 转换为响应格式
	var notificationList []respond.NotificationListItem
	for _, notification := range notifications {
		item := respond.NotificationListItem{
			Uuid:        notification.Uuid,
			Type:        notification.Type,
			Title:       notification.Title,
			Content:     notification.Content,
			RelatedId:   notification.RelatedId,
			RelatedType: notification.RelatedType,
			SenderId:    notification.SenderId,
			SenderName:  notification.SenderName,
			SenderAvatar: notification.SenderAvatar,
			Status:      notification.Status,
			ExtraData:   notification.ExtraData,
			CreatedAt:   notification.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if notification.ReadAt.Valid {
			item.ReadAt = notification.ReadAt.Time.Format("2006-01-02 15:04:05")
		}
		notificationList = append(notificationList, item)
	}

	return "获取通知列表成功", respond.GetNotificationListResponse{
		Data:  notificationList,
		Total: total,
	}, 0
}

// GetUnreadCount 获取未读通知数量
func (n *notificationService) GetUnreadCount(userId string, notificationType *int8) (message string, count int64, ret int) {
	query := dao.GormDB.Model(&model.Notification{}).
		Where("user_id = ? AND status = ?", userId, notification_status_enum.Unread)

	if notificationType != nil {
		query = query.Where("type = ?", *notificationType)
	}

	if err := query.Count(&count).Error; err != nil {
		zlog.Error("获取未读通知数量失败: " + err.Error())
		return constants.SYSTEM_ERROR, 0, -1
	}

	return "获取未读数量成功", count, 0
}

// MarkAsRead 标记通知为已读
func (n *notificationService) MarkAsRead(userId string, notificationIds []string) (message string, count int64, ret int) {
	now := time.Now()
	readAt := sql.NullTime{
		Time:  now,
		Valid: true,
	}

	var result *gorm.DB
	if len(notificationIds) == 0 {
		// 标记所有未读为已读
		result = dao.GormDB.Model(&model.Notification{}).
			Where("user_id = ? AND status = ?", userId, notification_status_enum.Unread).
			Updates(map[string]interface{}{
				"status":  notification_status_enum.Read,
				"read_at": readAt,
			})
	} else {
		// 标记指定通知为已读
		result = dao.GormDB.Model(&model.Notification{}).
			Where("user_id = ? AND uuid IN ? AND status = ?", userId, notificationIds, notification_status_enum.Unread).
			Updates(map[string]interface{}{
				"status":  notification_status_enum.Read,
				"read_at": readAt,
			})
	}

	if result.Error != nil {
		zlog.Error("标记通知为已读失败: " + result.Error.Error())
		return constants.SYSTEM_ERROR, 0, -1
	}

	count = result.RowsAffected
	zlog.Info(fmt.Sprintf("标记通知为已读成功: user_id=%s, count=%d", userId, count))
	return "标记为已读成功", count, 0
}

// DeleteNotification 删除通知（软删除）
func (n *notificationService) DeleteNotification(userId string, notificationIds []string) (message string, count int64, ret int) {
	result := dao.GormDB.Where("user_id = ? AND uuid IN ?", userId, notificationIds).
		Delete(&model.Notification{})

	if result.Error != nil {
		zlog.Error("删除通知失败: " + result.Error.Error())
		return constants.SYSTEM_ERROR, 0, -1
	}

	count = result.RowsAffected
	zlog.Info(fmt.Sprintf("删除通知成功: user_id=%s, count=%d", userId, count))
	return "删除通知成功", count, 0
}

// ClearAllNotification 清空所有通知（软删除）
func (n *notificationService) ClearAllNotification(userId string, notificationType *int8) (message string, count int64, ret int) {
	query := dao.GormDB.Where("user_id = ?", userId)

	if notificationType != nil {
		query = query.Where("type = ?", *notificationType)
	}

	result := query.Delete(&model.Notification{})

	if result.Error != nil {
		zlog.Error("清空通知失败: " + result.Error.Error())
		return constants.SYSTEM_ERROR, 0, -1
	}

	count = result.RowsAffected
	zlog.Info(fmt.Sprintf("清空通知成功: user_id=%s, type=%v, count=%d", userId, notificationType, count))
	return "清空通知成功", count, 0
}

