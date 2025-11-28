package request

type DeleteNotificationRequest struct {
	// UserId 从 JWT token 中获取，不需要前端传递
	NotificationIds []string `json:"notification_ids" binding:"required"`
}

