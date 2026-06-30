package sms

import (
	"encoding/json"
	"strconv"
	"strings"

	"gochat/internal/config"
	"gochat/pkg/constants"
	"gochat/pkg/zlog"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dypnsapi20170525 "github.com/alibabacloud-go/dypnsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
)

var smsClient *dypnsapi20170525.Client

// createClient 使用 AK&SK 初始化号码认证客户端。
func createClient() (result *dypnsapi20170525.Client, err error) {
	authCfg := config.GetConfig().AuthCodeConfig
	if smsClient == nil {
		zlog.Info(
			"初始化阿里云短信认证客户端",
			zap.String("endpoint", "dypnsapi.aliyuncs.com"),
		)
		openapiConfig := &openapi.Config{
			AccessKeyId:     tea.String(authCfg.AccessKeyID),
			AccessKeySecret: tea.String(authCfg.AccessKeySecret),
		}
		openapiConfig.Endpoint = tea.String("dypnsapi.aliyuncs.com")
		smsClient, err = dypnsapi20170525.NewClient(openapiConfig)
		if err != nil {
			zlog.Error("初始化阿里云短信认证客户端失败", zap.Error(err))
			return nil, err
		}
	}
	return smsClient, err
}

// SendVerificationCode 发送短信验证码。验证码由阿里云生成，不再由本地 Redis 保存。
func SendVerificationCode(telephone string) (string, int) {
	client, err := createClient()
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}

	authCfg := config.GetConfig().AuthCodeConfig
	if message, ret := validateAuthCodeConfig(authCfg); ret != 0 {
		return message, ret
	}

	trimmedSchemeName := strings.TrimSpace(authCfg.SchemeName)
	zlog.Info(
		"开始发送短信验证码",
		zap.String("telephone", maskTelephone(telephone)),
		zap.String("scheme_name", trimmedSchemeName),
		zap.String("country_code", authCfg.CountryCode),
		zap.Int("valid_time_seconds", authCfg.ValidTime),
		zap.Int("interval_seconds", authCfg.Interval),
		zap.Int("code_length", authCfg.CodeLength),
	)

	templateParam, err := buildTemplateParam()
	if err != nil {
		zlog.Error(
			"构造短信模板参数失败",
			zap.String("telephone", maskTelephone(telephone)),
			zap.Error(err),
		)
		return constants.SYSTEM_ERROR, -1
	}

	request := &dypnsapi20170525.SendSmsVerifyCodeRequest{
		PhoneNumber:      tea.String(telephone),
		SignName:         tea.String(authCfg.SignName),
		TemplateCode:     tea.String(authCfg.TemplateCode),
		TemplateParam:    tea.String(templateParam),
		CountryCode:      tea.String(authCfg.CountryCode),
		CodeLength:       tea.Int64(int64(authCfg.CodeLength)),
		ValidTime:        tea.Int64(int64(authCfg.ValidTime)),
		DuplicatePolicy:  tea.Int64(int64(authCfg.DuplicatePolicy)),
		Interval:         tea.Int64(int64(authCfg.Interval)),
		CodeType:         tea.Int64(int64(authCfg.CodeType)),
		ReturnVerifyCode: tea.Bool(authCfg.ReturnVerifyCode),
		AutoRetry:        tea.Int64(int64(authCfg.AutoRetry)),
	}
	if trimmedSchemeName != "" {
		request.SchemeName = tea.String(trimmedSchemeName)
	}

	runtime := &dara.RuntimeOptions{}
	response, err := client.SendSmsVerifyCodeWithOptions(request, runtime)
	if err != nil {
		zlog.Error(
			"调用阿里云发送验证码接口失败",
			zap.String("telephone", maskTelephone(telephone)),
			zap.String("scheme_name", trimmedSchemeName),
			zap.Error(err),
		)
		return constants.SYSTEM_ERROR, -1
	}
	if message, ret := validateSendResponse(response); ret != 0 {
		return message, ret
	}

	zlog.Info(
		"短信验证码发送成功",
		zap.String("telephone", maskTelephone(telephone)),
		zap.String("scheme_name", trimmedSchemeName),
		zap.String("code", teaString(response.Body.Code)),
		zap.String("message", teaString(response.Body.Message)),
		zap.String("request_id", teaString(response.Body.RequestId)),
		zap.String("biz_id", sendBizID(response)),
	)

	return "验证码发送成功，请及时在对应电话查收短信", 0
}

// CheckVerificationCode 校验短信验证码。校验逻辑由阿里云负责。
func CheckVerificationCode(telephone string, verifyCode string) (string, int) {
	client, err := createClient()
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}

	authCfg := config.GetConfig().AuthCodeConfig
	if message, ret := validateAuthCodeConfig(authCfg); ret != 0 {
		return message, ret
	}

	trimmedSchemeName := strings.TrimSpace(authCfg.SchemeName)
	zlog.Info(
		"开始校验短信验证码",
		zap.String("telephone", maskTelephone(telephone)),
		zap.String("scheme_name", trimmedSchemeName),
		zap.String("country_code", authCfg.CountryCode),
	)

	request := &dypnsapi20170525.CheckSmsVerifyCodeRequest{
		PhoneNumber:    tea.String(telephone),
		VerifyCode:     tea.String(verifyCode),
		CountryCode:    tea.String(authCfg.CountryCode),
		CaseAuthPolicy: tea.Int64(int64(authCfg.CaseAuthPolicy)),
	}
	if trimmedSchemeName != "" {
		request.SchemeName = tea.String(trimmedSchemeName)
	}

	runtime := &dara.RuntimeOptions{}
	response, err := client.CheckSmsVerifyCodeWithOptions(request, runtime)
	if err != nil {
		zlog.Error(
			"调用阿里云校验验证码接口失败",
			zap.String("telephone", maskTelephone(telephone)),
			zap.String("scheme_name", trimmedSchemeName),
			zap.Error(err),
		)
		return constants.SYSTEM_ERROR, -1
	}
	if message, ret := validateCheckResponse(response); ret != 0 {
		return message, ret
	}

	zlog.Info(
		"短信验证码校验成功",
		zap.String("telephone", maskTelephone(telephone)),
		zap.String("scheme_name", trimmedSchemeName),
		zap.String("code", teaString(response.Body.Code)),
		zap.String("message", teaString(response.Body.Message)),
		zap.String("verify_result", checkVerifyResult(response)),
	)

	return "验证码校验成功", 0
}

func validateAuthCodeConfig(authCfg config.AuthCodeConfig) (string, int) {
	if strings.TrimSpace(authCfg.AccessKeyID) == "" || strings.TrimSpace(authCfg.AccessKeySecret) == "" {
		zlog.Warn("短信认证 AccessKey 配置不完整")
		return "短信认证配置不完整，请检查 AccessKey", -2
	}
	if strings.TrimSpace(authCfg.SignName) == "" || strings.TrimSpace(authCfg.TemplateCode) == "" {
		zlog.Warn("短信认证签名或模板配置不完整")
		return "短信认证配置不完整，请检查签名和模板", -2
	}
	return "", 0
}

func buildTemplateParam() (string, error) {
	authConfig := config.GetConfig().AuthCodeConfig
	validMin := authConfig.ValidTime / 60
	if authConfig.ValidTime%60 != 0 {
		validMin++
	}
	templateParam := map[string]string{
		"code": "##code##",
		"min":  strconv.Itoa(validMin),
	}
	templateParamBytes, err := json.Marshal(templateParam)
	if err != nil {
		return "", err
	}
	return string(templateParamBytes), nil
}

func validateSendResponse(response *dypnsapi20170525.SendSmsVerifyCodeResponse) (string, int) {
	if response == nil || response.Body == nil {
		zlog.Error("阿里云发送验证码返回体为空")
		return constants.SYSTEM_ERROR, -1
	}
	if response.Body.Success == nil || !tea.BoolValue(response.Body.Success) {
		zlog.Warn(
			"阿里云发送验证码失败",
			zap.String("code", teaString(response.Body.Code)),
			zap.String("message", teaString(response.Body.Message)),
			zap.String("request_id", teaString(response.Body.RequestId)),
		)
		if response.Body.Message != nil {
			return tea.StringValue(response.Body.Message), -2
		}
		return "验证码发送失败", -2
	}
	if response.Body.Code == nil || tea.StringValue(response.Body.Code) != "OK" {
		zlog.Warn(
			"阿里云发送验证码返回非 OK 状态",
			zap.String("code", teaString(response.Body.Code)),
			zap.String("message", teaString(response.Body.Message)),
			zap.String("request_id", teaString(response.Body.RequestId)),
		)
		if response.Body.Message != nil {
			return tea.StringValue(response.Body.Message), -2
		}
		return "验证码发送失败", -2
	}
	return "", 0
}

func validateCheckResponse(response *dypnsapi20170525.CheckSmsVerifyCodeResponse) (string, int) {
	if response == nil || response.Body == nil {
		zlog.Error("阿里云校验验证码返回体为空")
		return constants.SYSTEM_ERROR, -1
	}
	if response.Body.Success == nil || !tea.BoolValue(response.Body.Success) {
		zlog.Warn(
			"阿里云校验验证码失败",
			zap.String("code", teaString(response.Body.Code)),
			zap.String("message", teaString(response.Body.Message)),
			zap.String("verify_result", checkVerifyResult(response)),
		)
		if response.Body.Message != nil {
			return tea.StringValue(response.Body.Message), -2
		}
		return "验证码校验失败", -2
	}
	if response.Body.Code == nil || tea.StringValue(response.Body.Code) != "OK" {
		zlog.Warn(
			"阿里云校验验证码返回非 OK 状态",
			zap.String("code", teaString(response.Body.Code)),
			zap.String("message", teaString(response.Body.Message)),
			zap.String("verify_result", checkVerifyResult(response)),
		)
		if response.Body.Message != nil {
			return tea.StringValue(response.Body.Message), -2
		}
		return "验证码校验失败", -2
	}
	if response.Body.Model == nil || response.Body.Model.VerifyResult == nil || tea.StringValue(response.Body.Model.VerifyResult) != "PASS" {
		zlog.Warn(
			"短信验证码校验未通过",
			zap.String("code", teaString(response.Body.Code)),
			zap.String("message", teaString(response.Body.Message)),
			zap.String("verify_result", checkVerifyResult(response)),
		)
		return "验证码不正确，请重试", -2
	}
	return "", 0
}

func maskTelephone(telephone string) string {
	trimmedTelephone := strings.TrimSpace(telephone)
	if len(trimmedTelephone) <= 7 {
		return trimmedTelephone
	}
	return trimmedTelephone[:3] + "****" + trimmedTelephone[len(trimmedTelephone)-4:]
}

func teaString(value *string) string {
	if value == nil {
		return ""
	}
	return tea.StringValue(value)
}

func sendBizID(response *dypnsapi20170525.SendSmsVerifyCodeResponse) string {
	if response == nil || response.Body == nil || response.Body.Model == nil || response.Body.Model.BizId == nil {
		return ""
	}
	return tea.StringValue(response.Body.Model.BizId)
}

func checkVerifyResult(response *dypnsapi20170525.CheckSmsVerifyCodeResponse) string {
	if response == nil || response.Body == nil || response.Body.Model == nil || response.Body.Model.VerifyResult == nil {
		return ""
	}
	return tea.StringValue(response.Body.Model.VerifyResult)
}
