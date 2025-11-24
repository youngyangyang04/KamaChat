package request

type GetNotificationListRequest struct {
	UserId   string `json:"user_id" binding:"required"` // 用户ID
	Page     int    `json:"page"`                       // 页码，从1开始
	PageSize int    `json:"page_size"`                  // 每页数量，默认20
	Type     *int8  `json:"type"`                       // 通知类型筛选（可选）
	Status   *int8  `json:"status"`                     // 状态筛选（可选，0.未读，1.已读）
}
