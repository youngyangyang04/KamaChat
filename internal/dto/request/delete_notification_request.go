package request

type DeleteNotificationRequest struct {
	UserId          string   `json:"user_id" binding:"required"`
	NotificationIds []string `json:"notification_ids" binding:"required"`
}

