package notification_type_enum

const (
	// 好友申请
	FriendApply = iota + 1
	// 群聊邀请
	GroupInvite
	// 新消息（用户不在当前聊天窗口时）
	NewMessage
	// 系统通知
	System
	// 群聊申请（用户申请加入群聊）
	GroupApply
	// 好友申请已通过
	ContactAccepted
	// 好友申请已拒绝
	ContactRejected
	// 群聊申请已通过
	GroupAccepted
	// 群聊申请已拒绝
	GroupRejected
)

