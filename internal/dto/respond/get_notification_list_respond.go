package respond

type NotificationListItem struct {
	Uuid        string `json:"uuid"`
	Type        int8   `json:"type"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	RelatedId   string `json:"related_id"`
	RelatedType string `json:"related_type"`
	SenderId    string `json:"sender_id"`
	SenderName  string `json:"sender_name"`
	SenderAvatar string `json:"sender_avatar"`
	Status      int8   `json:"status"`
	ExtraData   string `json:"extra_data"`
	CreatedAt   string `json:"created_at"`
	ReadAt      string `json:"read_at,omitempty"`
}

type GetNotificationListResponse struct {
	Data  []NotificationListItem `json:"data"`
	Total int64                   `json:"total"`  // 总数量
}

