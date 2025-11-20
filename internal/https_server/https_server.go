package https_server

import (
	v1 "kama_chat_server/api/v1"
	"kama_chat_server/internal/config"
	"kama_chat_server/pkg/middleware"
	"kama_chat_server/pkg/ssl"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var GE *gin.Engine

func init() {
	GE = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	GE.Use(cors.New(corsConfig))

	// 根据配置决定是否启用 HTTPS TLS 重定向
	if config.GetConfig().MainConfig.EnableHTTPS {
		GE.Use(ssl.TlsHandler(config.GetConfig().MainConfig.Host, config.GetConfig().MainConfig.Port))
	}

	// 静态文件服务
	GE.Static("/static/avatars", config.GetConfig().StaticAvatarPath)
	GE.Static("/static/files", config.GetConfig().StaticFilePath)

	// 无需认证的路由
	GE.POST("/login", v1.Login)
	GE.POST("/register", v1.Register)
	GE.POST("/registerWithCrypto", v1.RegisterWithCrypto) // 带加密密钥的注册

	// 需要认证的路由组
	auth := GE.Group("/")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("user/updateUserInfo", v1.UpdateUserInfo)
		auth.POST("user/getUserInfoList", v1.GetUserInfoList)
		auth.POST("user/ableUsers", v1.AbleUsers)
		auth.POST("user/getUserInfo", v1.GetUserInfo)
		auth.POST("user/disableUsers", v1.DisableUsers)
		auth.POST("user/deleteUsers", v1.DeleteUsers)
		auth.POST("user/setAdmin", v1.SetAdmin)
		auth.POST("user/sendSmsCode", v1.SendSmsCode)
		auth.POST("user/smsLogin", v1.SmsLogin)
		auth.POST("user/wsLogout", v1.WsLogout)
		auth.POST("group/createGroup", v1.CreateGroup)
		auth.POST("group/loadMyGroup", v1.LoadMyGroup)
		auth.POST("group/checkGroupAddMode", v1.CheckGroupAddMode)
		auth.POST("group/enterGroupDirectly", v1.EnterGroupDirectly)
		auth.POST("group/leaveGroup", v1.LeaveGroup)
		auth.POST("group/dismissGroup", v1.DismissGroup)
		auth.POST("group/getGroupInfo", v1.GetGroupInfo)
		auth.POST("group/getGroupInfoList", v1.GetGroupInfoList)
		auth.POST("group/deleteGroups", v1.DeleteGroups)
		auth.POST("group/setGroupsStatus", v1.SetGroupsStatus)
		auth.POST("group/updateGroupInfo", v1.UpdateGroupInfo)
		auth.POST("group/getGroupMemberList", v1.GetGroupMemberList)
		auth.POST("group/removeGroupMembers", v1.RemoveGroupMembers)
		auth.POST("session/openSession", v1.OpenSession)
		auth.POST("session/getUserSessionList", v1.GetUserSessionList)
		auth.POST("session/getGroupSessionList", v1.GetGroupSessionList)
		auth.POST("session/deleteSession", v1.DeleteSession)
		auth.POST("session/checkOpenSessionAllowed", v1.CheckOpenSessionAllowed)
		auth.POST("contact/getUserList", v1.GetUserList)
		auth.POST("contact/loadMyJoinedGroup", v1.LoadMyJoinedGroup)
		auth.POST("contact/getContactInfo", v1.GetContactInfo)
		auth.POST("contact/deleteContact", v1.DeleteContact)
		auth.POST("contact/applyContact", v1.ApplyContact)
		auth.POST("contact/getNewContactList", v1.GetNewContactList)
		auth.POST("contact/passContactApply", v1.PassContactApply)
		auth.POST("contact/blackContact", v1.BlackContact)
		auth.POST("contact/cancelBlackContact", v1.CancelBlackContact)
		auth.POST("contact/getAddGroupList", v1.GetAddGroupList)
		auth.POST("contact/refuseContactApply", v1.RefuseContactApply)
		auth.POST("contact/blackApply", v1.BlackApply)
		auth.POST("message/getMessageList", v1.GetMessageList)
		auth.POST("message/getGroupMessageList", v1.GetGroupMessageList)
		auth.POST("message/uploadAvatar", v1.UploadAvatar)
		auth.POST("message/uploadFile", v1.UploadFile)
		auth.POST("message/sendEncryptedMessage", v1.SendEncryptedMessage)
		auth.POST("chatroom/getCurContactListInChatRoom", v1.GetCurContactListInChatRoom)

		// 加密相关接口
		auth.GET("crypto/getPublicKeyBundle", v1.GetPublicKeyBundle)
		auth.GET("crypto/getOneTimePreKeyCount", v1.GetOneTimePreKeyCount)
		auth.POST("crypto/replenishOneTimePreKeys", v1.ReplenishOneTimePreKeys)
		auth.POST("crypto/rotateSignedPreKey", v1.RotateSignedPreKey)

		auth.GET("wss", v1.WsLogin)
	}

}
