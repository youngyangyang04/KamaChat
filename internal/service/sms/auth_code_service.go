package sms

import (
	"fmt"
	"gochat/internal/config"
	"gochat/internal/service/redis"
	"gochat/pkg/constants"
	"gochat/pkg/util/random"
	"gochat/pkg/zlog"
	"strconv"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

var smsClient *dysmsapi20170525.Client

// createClient 使用AK&SK初始化账号Client
func createClient() (result *dysmsapi20170525.Client, err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	accessKeyID := config.GetConfig().AuthCodeConfig.AccessKeyID
	accessKeySecret := config.GetConfig().AuthCodeConfig.AccessKeySecret
	if smsClient == nil {
		openapiConfig := &openapi.Config{
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
			AccessKeyId: tea.String(accessKeyID),
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
			AccessKeySecret: tea.String(accessKeySecret),
		}
		// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
		openapiConfig.Endpoint = tea.String("dysmsapi.aliyuncs.com")
		smsClient, err = dysmsapi20170525.NewClient(openapiConfig)
	}
	return smsClient, err
}

func VerificationCode(telephone string) (string, int) {
	client, err := createClient()
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	key := "auth_code_" + telephone
	code, err := redis.GetKey(key)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}

	if code != "" {
		// 直接返回，验证码还没过期，用户应该去输验证码
		message := "目前还不能发送验证码，请输入已发送的验证码"
		zlog.Info(message)
		return message, -2
	}
	// 验证码过期，重新生成
	code = strconv.Itoa(random.GetRandomInt(6))
	fmt.Println(code)
	err = redis.SetKeyEx(key, code, time.Minute) // 1分钟有效
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"), // 短信模板
		PhoneNumbers:  tea.String(telephone),
		TemplateParam: tea.String("{\"code\":\"" + code + "\"}"),
	}

	runtime := &util.RuntimeOptions{}
	// 目前使用的是测试专用签名，签名必须是“阿里云短信测试”，模板code为“SMS_154950909”
	rsp, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	zlog.Info(*util.ToJSONString(rsp))
	return "验证码发送成功，请及时在对应电话查收短信", 0
}
