package request

type GetUnreadCountRequest struct {
	UserId string `json:"user_id" binding:"required"`
	Type   *int8  `json:"type"`  // 通知类型筛选（可选）
}

