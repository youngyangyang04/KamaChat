package request

type ClearAllNotificationRequest struct {
	// UserId 从 JWT token 中获取，不需要前端传递
	Type *int8 `json:"type"` // 通知类型筛选（可选），为空则清空所有类型
}

