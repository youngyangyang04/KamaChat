package request

type MarkAsReadRequest struct {
	UserId          string   `json:"user_id" binding:"required"`
	NotificationIds []string `json:"notification_ids"`  // 通知ID列表，为空则标记所有未读为已读
}

