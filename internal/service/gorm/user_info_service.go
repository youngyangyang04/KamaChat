package gorm

import (
	"encoding/json"
	"errors"
	"fmt"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	myredis "kama_chat_server/internal/service/redis"
	"kama_chat_server/internal/service/sms"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/enum/user_info/user_status_enum"
	jwtutil "kama_chat_server/pkg/util/jwt"
	"kama_chat_server/pkg/util/password"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
	"regexp"
	"time"

	redis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
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
func (u *userInfoService) Login(loginReq request.LoginRequest) (string, *respond.LoginRespond, int) {
	var user model.UserInfo
	res := dao.GormDB.First(&user, "account = ?", loginReq.Account)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "账号不存在，请注册"
			zlog.Error(message)
			return message, nil, -2
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 验证密码
	if !password.VerifyPassword(user.Password, loginReq.Password) {
		message := "密码不正确，请重试"
		zlog.Error(message)
		return message, nil, -2
	}

	// 生成 JWT token
	token, err := jwtutil.GenerateToken(user.Uuid, user.Account, user.Nickname, user.IsAdmin)
	if err != nil {
		zlog.Error("生成 token 失败: " + err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	loginRsp := &respond.LoginRespond{
		Token:     token,
		Uuid:      user.Uuid,
		Account:   user.Account,
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

	return "登录成功", loginRsp, 0
}

// SmsLogin 验证码登录
func (u *userInfoService) SmsLogin(req request.SmsLoginRequest) (string, *respond.LoginRespond, int) {
	var user model.UserInfo
	res := dao.GormDB.First(&user, "telephone = ?", req.Telephone)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "用户不存在，请注册"
			zlog.Error(message)
			return message, nil, -2
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	key := "auth_code_" + req.Telephone
	code, err := myredis.GetKey(key)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	if code != req.SmsCode {
		message := "验证码不正确，请重试"
		zlog.Info(message)
		return message, nil, -2
	} else {
		if err := myredis.DelKeyIfExists(key); err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}

	loginRsp := &respond.LoginRespond{
		Uuid:      user.Uuid,
		Account:   user.Account,
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

	return "登陆成功", loginRsp, 0
}

// SendSmsCode 发送短信验证码 - 验证码登录
func (u *userInfoService) SendSmsCode(telephone string) (string, int) {
	return sms.VerificationCode(telephone)
}

// checkTelephoneExist 检查手机号是否存在
func (u *userInfoService) checkTelephoneExist(telephone string) (string, int) {
	var user model.UserInfo
	// gorm默认排除软删除，所以翻译过来的select语句是SELECT * FROM `user_info` WHERE telephone = '18089596095' AND `user_info`.`deleted_at` IS NULL ORDER BY `user_info`.`id` LIMIT 1
	if res := dao.GormDB.Where("telephone = ?", telephone).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("该电话不存在，可以注册")
			return "", 0
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	zlog.Info("该电话已经存在，注册失败")
	return "该电话已经存在，注册失败", -2
}

// checkAccountExist 检查账号是否存在
func (u *userInfoService) checkAccountExist(account string) (string, int) {
	var user model.UserInfo
	if res := dao.GormDB.Where("account = ?", account).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("该账号不存在，可以注册")
			return "", 0
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	message := "该账号已被注册"
	zlog.Info(message)
	return message, -2
}

// Register 注册，返回(message, register_respond_string, error)
func (u *userInfoService) Register(registerReq request.RegisterRequest) (string, *respond.RegisterRespond, int) {
	// 检查账号是否已存在
	message, ret := u.checkAccountExist(registerReq.Account)
	if ret != 0 {
		return message, nil, ret
	}

	// 创建新用户
	var newUser model.UserInfo
	newUser.Uuid = "U" + random.GetNowAndLenRandomString(11)
	newUser.Account = registerReq.Account
	newUser.Nickname = registerReq.Nickname

	// 加密密码
	hashedPassword, err := password.HashPassword(registerReq.Password)
	if err != nil {
		zlog.Error("密码加密失败: " + err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	newUser.Password = hashedPassword

	newUser.Avatar = "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
	newUser.CreatedAt = time.Now()
	newUser.IsAdmin = 0
	newUser.Status = user_status_enum.NORMAL

	res := dao.GormDB.Create(&newUser)
	if res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 生成 JWT token
	token, err := jwtutil.GenerateToken(newUser.Uuid, newUser.Account, newUser.Nickname, newUser.IsAdmin)
	if err != nil {
		zlog.Error("生成 token 失败: " + err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	registerRsp := &respond.RegisterRespond{
		Token:     token,
		Uuid:      newUser.Uuid,
		Account:   newUser.Account,
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

	return "注册成功", registerRsp, 0
}

// RegisterWithCrypto 注册（带加密密钥），返回(message, register_respond_string, error)
func (u *userInfoService) RegisterWithCrypto(registerReq request.RegisterCryptoRequest) (string, *respond.RegisterRespond, int) {
	// 检查账号是否已存在
	message, ret := u.checkAccountExist(registerReq.Account)
	if ret != 0 {
		return message, nil, ret
	}

	// 创建新用户
	var newUser model.UserInfo
	newUser.Uuid = "U" + random.GetNowAndLenRandomString(11)
	newUser.Account = registerReq.Account
	newUser.Nickname = registerReq.Nickname

	// 加密密码
	hashedPassword, err := password.HashPassword(registerReq.Password)
	if err != nil {
		zlog.Error("密码加密失败: " + err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	newUser.Password = hashedPassword

	newUser.Avatar = "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
	newUser.CreatedAt = time.Now()
	newUser.IsAdmin = 0
	newUser.Status = user_status_enum.NORMAL

	// 使用事务确保用户创建和密钥保存的原子性
	tx := dao.GormDB.Begin()
	if tx.Error != nil {
		zlog.Error("开始事务失败: " + tx.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 创建用户
	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		zlog.Error("创建用户失败: " + err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 保存加密公钥（在同一事务中）
	// 直接传入事务实例，避免修改全局变量
	err = CryptoKeyService.SaveUserPublicKeys(newUser.Uuid, &registerReq, tx)

	if err != nil {
		tx.Rollback()
		zlog.Error("保存加密密钥失败: " + err.Error())
		return "保存加密密钥失败", nil, -1
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zlog.Error("提交事务失败: " + err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 生成 JWT token
	token, err := jwtutil.GenerateToken(newUser.Uuid, newUser.Account, newUser.Nickname, newUser.IsAdmin)
	if err != nil {
		zlog.Error("生成 token 失败: " + err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	registerRsp := &respond.RegisterRespond{
		Token:     token,
		Uuid:      newUser.Uuid,
		Account:   newUser.Account,
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

	zlog.Info("用户注册成功（带加密）: " + newUser.Account)
	return "注册成功", registerRsp, 0
}

// UpdateUserInfo 修改用户信息
// 某用户修改了信息，可能会影响contact_user_list，不需要删除redis的contact_user_list，timeout之后会自己更新
// 但是需要更新redis的user_info，因为可能影响用户搜索
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
	//if err := myredis.DelKeysWithPattern("user_info_" + updateReq.Uuid); err != nil {
	//	zlog.Error(err.Error())
	//}
	return "修改用户信息成功", 0
}

// GetUserInfoList 获取用户列表除了ownerId之外 - 管理员
// 管理员少，而且如果用户更改了，那么管理员会一直频繁删除redis，更新redis，比较麻烦，所以管理员暂时不使用redis缓存
func (u *userInfoService) GetUserInfoList(ownerId string) (string, []respond.GetUserListRespond, int) {
	// redis中没有数据，从数据库中获取
	var users []model.UserInfo
	// 获取所有的用户
	if res := dao.GormDB.Unscoped().Where("uuid != ?", ownerId).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	var rsp []respond.GetUserListRespond
	for _, user := range users {
		rp := respond.GetUserListRespond{
			Uuid:      user.Uuid,
			Telephone: user.Telephone,
			Nickname:  user.Nickname,
			Status:    user.Status,
			IsAdmin:   user.IsAdmin,
		}
		if user.DeletedAt.Valid {
			rp.IsDeleted = true
		} else {
			rp.IsDeleted = false
		}
		rsp = append(rsp, rp)
	}
	return "获取用户列表成功", rsp, 0
}

// AbleUsers 启用用户
// 用户是否启用禁用需要实时更新contact_user_list状态，所以redis的contact_user_list需要删除
func (u *userInfoService) AbleUsers(uuidList []string) (string, int) {
	var users []model.UserInfo
	if res := dao.GormDB.Model(model.UserInfo{}).Where("uuid in (?)", uuidList).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for _, user := range users {
		user.Status = user_status_enum.NORMAL
		if res := dao.GormDB.Save(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}
	// 删除所有"contact_user_list"开头的key
	//if err := myredis.DelKeysWithPrefix("contact_user_list"); err != nil {
	//	zlog.Error(err.Error())
	//}
	return "启用用户成功", 0
}

// DisableUsers 禁用用户
// 用户是否启用禁用需要实时更新contact_user_list状态，所以redis的contact_user_list需要删除
func (u *userInfoService) DisableUsers(uuidList []string) (string, int) {
	var users []model.UserInfo
	if res := dao.GormDB.Model(model.UserInfo{}).Where("uuid in (?)", uuidList).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for _, user := range users {
		user.Status = user_status_enum.DISABLE
		if res := dao.GormDB.Save(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
		var sessionList []model.Session
		if res := dao.GormDB.Where("send_id = ? or receive_id = ?", user.Uuid, user.Uuid).Find(&sessionList); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
		for _, session := range sessionList {
			var deletedAt gorm.DeletedAt
			deletedAt.Time = time.Now()
			deletedAt.Valid = true
			session.DeletedAt = deletedAt
			if res := dao.GormDB.Save(&session); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
	}
	// 删除所有"contact_user_list"开头的key
	//if err := myredis.DelKeysWithPrefix("contact_user_list"); err != nil {
	//	zlog.Error(err.Error())
	//}
	return "禁用用户成功", 0
}

// DeleteUsers 删除用户
// 用户是否启用禁用需要实时更新contact_user_list状态，所以redis的contact_user_list需要删除
func (u *userInfoService) DeleteUsers(uuidList []string) (string, int) {
	var users []model.UserInfo
	if res := dao.GormDB.Model(model.UserInfo{}).Where("uuid in (?)", uuidList).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for _, user := range users {
		user.DeletedAt.Valid = true
		user.DeletedAt.Time = time.Now()
		if res := dao.GormDB.Save(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}

		// 删除会话
		var sessionList []model.Session
		if res := dao.GormDB.Where("send_id = ? or receive_id = ?", user.Uuid, user.Uuid).Find(&sessionList); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info(res.Error.Error())
			} else {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
		for _, session := range sessionList {
			var deletedAt gorm.DeletedAt
			deletedAt.Time = time.Now()
			deletedAt.Valid = true
			session.DeletedAt = deletedAt
			if res := dao.GormDB.Save(&session); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}

		// 删除联系人
		var contactList []model.UserContact
		if res := dao.GormDB.Where("user_id = ? or contact_id = ?", user.Uuid, user.Uuid).Find(&contactList); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info(res.Error.Error())
			} else {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
		for _, contact := range contactList {
			var deletedAt gorm.DeletedAt
			deletedAt.Time = time.Now()
			deletedAt.Valid = true
			contact.DeletedAt = deletedAt
			if res := dao.GormDB.Save(&contact); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}

		// 删除申请记录
		var applyList []model.ContactApply
		if res := dao.GormDB.Where("user_id = ? or contact_id = ?", user.Uuid, user.Uuid).Find(&applyList); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info(res.Error.Error())
			} else {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
		for _, apply := range applyList {
			var deletedAt gorm.DeletedAt
			deletedAt.Time = time.Now()
			deletedAt.Valid = true
			apply.DeletedAt = deletedAt
			if res := dao.GormDB.Save(&apply); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}

	}
	// 删除所有"contact_user_list"开头的key
	//if err := myredis.DelKeysWithPrefix("contact_user_list"); err != nil {
	//	zlog.Error(err.Error())
	//}
	return "删除用户成功", 0
}

// GetUserInfo 获取用户信息
func (u *userInfoService) GetUserInfo(uuid string) (string, *respond.GetUserInfoRespond, int) {
	// redis
	zlog.Info(uuid)
	rspString, err := myredis.GetKeyNilIsErr("user_info_" + uuid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Info(err.Error())
			var user model.UserInfo
			if res := dao.GormDB.Where("uuid = ?", uuid).Find(&user); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, nil, -1
			}
			rsp := respond.GetUserInfoRespond{
				Uuid:      user.Uuid,
				Telephone: user.Telephone,
				Nickname:  user.Nickname,
				Avatar:    user.Avatar,
				Birthday:  user.Birthday,
				Email:     user.Email,
				Gender:    user.Gender,
				Signature: user.Signature,
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
				IsAdmin:   user.IsAdmin,
				Status:    user.Status,
			}
			//rspString, err := json.Marshal(rsp)
			//if err != nil {
			//	zlog.Error(err.Error())
			//}
			//if err := myredis.SetKeyEx("user_info_"+uuid, string(rspString), constants.REDIS_TIMEOUT*time.Minute); err != nil {
			//	zlog.Error(err.Error())
			//}
			return "获取用户信息成功", &rsp, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	var rsp respond.GetUserInfoRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取用户信息成功", &rsp, 0
}

// SetAdmin 设置管理员
func (u *userInfoService) SetAdmin(uuidList []string, isAdmin int8) (string, int) {
	var users []model.UserInfo
	if res := dao.GormDB.Where("uuid = (?)", uuidList).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for _, user := range users {
		user.IsAdmin = isAdmin
		if res := dao.GormDB.Save(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}
	return "设置管理员成功", 0
}
