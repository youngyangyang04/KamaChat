package gorm

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/enum/user_info/user_status_enum"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
	"regexp"
	"time"
)

type userInfoService struct {
}

var UserInfoService = new(userInfoService)

// dao层加不了校验，在service层加
// checkTelephoneValid 检验电话是否有效
func (u *userInfoService) checkTelephoneValid(telephone string) bool {
	pattern := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	match, err := regexp.MatchString(pattern, telephone)
	if err != nil {
		zlog.Error(err.Error())
	}
	return match
}

// checkEmailValid 校验邮箱是否有效
func (u *userInfoService) checkEmailValid(email string) bool {
	pattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		zlog.Error(err.Error())
	}
	return match
}

// checkUserIsAdminOrNot 检验用户是否为管理员
func (u *userInfoService) checkUserIsAdminOrNot(user model.UserInfo) int8 {
	return user.IsAdmin
}

// Login 登录
func (u *userInfoService) Login(c *gin.Context, loginReq request.LoginRequest) (string, *respond.LoginRespond, int) {
	password := loginReq.Password
	var user model.UserInfo
	res := dao.GormDB.First(&user, "telephone = ?", loginReq.Telephone)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "用户不存在，请注册"
			zlog.Error(message)
			return message, nil, -2
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	if user.Password != password {
		message := "密码不正确，请重试"
		zlog.Error(message)
		return message, nil, -2
	}
	// 手机号验证，最后一步才调用api，省钱hhh
	//if err := sms.VerificationCode(loginReq.Telephone); err != nil {
	//	zlog.Error(err.Error())
	//	return "", err
	//}

	loginRsp := &respond.LoginRespond{
		Uuid:      user.Uuid,
		Telephone: user.Telephone,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		Signature: user.Signature,
		IsAdmin:   user.IsAdmin,
		Status:    user.Status,
	}
	year, month, day := user.CreatedAt.Date()
	loginRsp.CreatedAt = fmt.Sprintf("%d.%d.%d", year, month, day)

	return "登陆成功!", loginRsp, 0
}

// Register 注册，返回(message, register_respond_string, error)
func (u *userInfoService) Register(c *gin.Context, registerReq request.RegisterRequest) (string, *respond.RegisterRespond, int) {
	// 不用校验手机号，前端校验

	var newUser model.UserInfo
	res := dao.GormDB.First(&newUser, "telephone = ?", registerReq.Telephone)
	if res.Error == nil {
		// 用户已经存在，注册失败
		message := "用户已经存在，注册失败"
		zlog.Error(message)
		return message, nil, -2
	} else {
		// 其他报错
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
		// 可以继续注册
	}
	newUser.Uuid = "U" + random.GetNowAndLenRandomString(11)
	newUser.Telephone = registerReq.Telephone
	newUser.Password = registerReq.Password
	newUser.Nickname = registerReq.Nickname
	newUser.Avatar = "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
	newUser.CreatedAt = time.Now()
	newUser.IsAdmin = u.checkUserIsAdminOrNot(newUser)
	newUser.Status = user_status_enum.NORMAL
	// 手机号验证，最后一步才调用api，省钱hhh
	//err := sms.VerificationCode(registerReq.Telephone)
	//if err != nil {
	//	zlog.Error(err.Error())
	//	return "", err
	//}
	res = dao.GormDB.Create(&newUser)
	if res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	// 注册成功，chat client建立
	//if err := chat.NewClientInit(c, newUser.Uuid); err != nil {
	//	return "", err
	//}
	newUser.LastOnlineAt = time.Now()
	registerRsp := &respond.RegisterRespond{
		Uuid:      newUser.Uuid,
		Telephone: newUser.Telephone,
		Nickname:  newUser.Nickname,
		Email:     newUser.Email,
		Avatar:    newUser.Avatar,
		Gender:    newUser.Gender,
		Birthday:  newUser.Birthday,
		Signature: newUser.Signature,
		IsAdmin:   newUser.IsAdmin,
		Status:    newUser.Status,
	}
	year, month, day := newUser.CreatedAt.Date()
	registerRsp.CreatedAt = fmt.Sprintf("%d.%d.%d", year, month, day)

	return "注册成功!", registerRsp, 0
}

// UpdateUserInfo 修改用户信息
func (u *userInfoService) UpdateUserInfo(updateReq request.UpdateUserInfoRequest) (string, int) {
	var user model.UserInfo
	if res := dao.GormDB.First(&user, "uuid = ?", updateReq.Uuid); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	if updateReq.Email != "" {
		user.Email = updateReq.Email
	}
	if updateReq.Nickname != "" {
		user.Nickname = updateReq.Nickname
	}
	if updateReq.Birthday != "" {
		user.Birthday = updateReq.Birthday
	}
	if updateReq.Signature != "" {
		user.Signature = updateReq.Signature
	}
	if updateReq.Avatar != "" {
		user.Avatar = updateReq.Avatar
	}
	if res := dao.GormDB.Save(&user); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	return "修改用户信息成功", 0
}
