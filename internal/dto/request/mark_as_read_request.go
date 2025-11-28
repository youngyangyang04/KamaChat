package request

type MarkAsReadRequest struct {
	// UserId 从 JWT token 中获取，不需要前端传递
	NotificationIds []string `json:"notification_ids"` // 通知ID列表，为空则标记所有未读为已读
}

