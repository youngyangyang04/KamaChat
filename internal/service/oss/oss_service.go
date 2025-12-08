package oss

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"kama_chat_server/internal/config"
	"kama_chat_server/pkg/zlog"
	"strings"
	"time"

	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/google/uuid"
)

type ossService struct{}

var OSSService = new(ossService)

const (
	product = "oss"
)

// UploadCredential 上传凭证（使用 STS 临时凭证）
type UploadCredential struct {
	Policy           string `json:"policy"`
	SecurityToken    string `json:"security_token"`
	SignatureVersion string `json:"x_oss_signature_version"`
	Credential       string `json:"x_oss_credential"`
	Date             string `json:"x_oss_date"`
	Signature        string `json:"signature"`
	Host             string `json:"host"`
	Key              string `json:"key"`
	Dir              string `json:"dir"`
	ExpireAt         int64  `json:"expireAt"`
	MaxFileSize      int64  `json:"maxFileSize"`
}

// GenerateUploadCredential 生成上传凭证（使用 STS 临时凭证）
func (s *ossService) GenerateUploadCredential(userId string, fileType string) (*UploadCredential, error) {
	cfg := config.GetConfig().OSSConfig

	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.STSRoleArn == "" {
		return nil, fmt.Errorf("OSS 配置不完整，请检查环境变量 OSS_ACCESS_KEY_ID, OSS_ACCESS_KEY_SECRET, OSS_STS_ROLE_ARN")
	}

	// 生成唯一的文件键
	fileUUID := uuid.New().String()
	// 文件存储目录: encrypted/{userId}/
	dir := fmt.Sprintf("%s%s/", cfg.EncryptedFilePrefix, userId)
	// 完整的文件路径
	key := dir + fileUUID

	// 创建 STS 凭证配置
	credConfig := new(credentials.Config).
		SetType("ram_role_arn").
		SetAccessKeyId(cfg.AccessKeyID).
		SetAccessKeySecret(cfg.AccessKeySecret).
		SetRoleArn(cfg.STSRoleArn).
		SetRoleSessionName("oss-upload-session").
		SetPolicy("").
		SetRoleSessionExpiration(3600) // 1小时

	// 创建凭证提供器
	provider, err := credentials.NewCredential(credConfig)
	if err != nil {
		zlog.Error("创建 STS 凭证提供器失败: " + err.Error())
		return nil, fmt.Errorf("创建 STS 凭证失败: %v", err)
	}

	// 获取临时凭证
	cred, err := provider.GetCredential()
	if err != nil {
		zlog.Error("获取 STS 临时凭证失败: " + err.Error())
		return nil, fmt.Errorf("获取 STS 临时凭证失败: %v", err)
	}

	// 构建 Policy
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(time.Duration(cfg.UploadExpireSeconds) * time.Second)

	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": cfg.Bucket},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, cfg.Region, product)},
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			map[string]string{"x-oss-security-token": *cred.SecurityToken},
			[]any{"starts-with", "$key", dir}, // 限制上传目录
			[]any{"content-length-range", 0, cfg.MaxFileSize},
		},
	}

	// 将 policy 转换为 JSON
	policyBytes, err := json.Marshal(policyMap)
	if err != nil {
		zlog.Error("生成 Policy JSON 失败: " + err.Error())
		return nil, fmt.Errorf("生成 Policy 失败: %v", err)
	}

	// 构造待签名字符串（StringToSign）
	stringToSign := base64.StdEncoding.EncodeToString(policyBytes)

	// 构建 signing key 并生成签名
	signature := s.generateSignatureV4(*cred.AccessKeySecret, date, cfg.Region, stringToSign)

	// OSS 上传地址
	host := fmt.Sprintf("https://%s.%s", cfg.Bucket, cfg.Endpoint)

	credential := &UploadCredential{
		Policy:           stringToSign,
		SecurityToken:    *cred.SecurityToken,
		SignatureVersion: "OSS4-HMAC-SHA256",
		Credential:       fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", *cred.AccessKeyId, date, cfg.Region, product),
		Date:             utcTime.Format("20060102T150405Z"),
		Signature:        signature,
		Host:             host,
		Key:              key,
		Dir:              dir,
		ExpireAt:         expiration.Unix(),
		MaxFileSize:      cfg.MaxFileSize,
	}

	zlog.Info(fmt.Sprintf("生成 STS 上传凭证成功: userId=%s, key=%s", userId, key))
	return credential, nil
}

// generateSignatureV4 使用 OSS4-HMAC-SHA256 算法生成签名
func (s *ossService) generateSignatureV4(accessKeySecret, date, region, stringToSign string) string {
	hmacHash := func() hash.Hash { return sha256.New() }

	// 构建 signing key
	signingKey := "aliyun_v4" + accessKeySecret

	h1 := hmac.New(hmacHash, []byte(signingKey))
	io.WriteString(h1, date)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	io.WriteString(h2, region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	io.WriteString(h3, product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	// 生成签名
	h := hmac.New(hmacHash, h4Key)
	io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}

// GenerateDownloadURL 生成带签名的下载 URL（使用阿里云 OSS SDK）
func (s *ossService) GenerateDownloadURL(ossKey string) (string, error) {
	cfg := config.GetConfig().OSSConfig

	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.STSRoleArn == "" {
		return "", fmt.Errorf("OSS 配置不完整")
	}

	// 创建 STS 凭证
	credConfig := new(credentials.Config).
		SetType("ram_role_arn").
		SetAccessKeyId(cfg.AccessKeyID).
		SetAccessKeySecret(cfg.AccessKeySecret).
		SetRoleArn(cfg.STSRoleArn).
		SetRoleSessionName("oss-download-session").
		SetPolicy("").
		SetRoleSessionExpiration(3600)

	provider, err := credentials.NewCredential(credConfig)
	if err != nil {
		return "", fmt.Errorf("创建 STS 凭证失败: %v", err)
	}

	cred, err := provider.GetCredential()
	if err != nil {
		return "", fmt.Errorf("获取 STS 临时凭证失败: %v", err)
	}

	// 使用官方 SDK 创建 OSS 客户端
	client, err := alioss.New(
		cfg.Endpoint,
		*cred.AccessKeyId,
		*cred.AccessKeySecret,
		alioss.SecurityToken(*cred.SecurityToken),
	)
	if err != nil {
		return "", fmt.Errorf("创建 OSS 客户端失败: %v", err)
	}

	// 获取 bucket
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return "", fmt.Errorf("获取 Bucket 失败: %v", err)
	}

	// 生成签名 URL
	expireSeconds := cfg.DownloadExpireSeconds
	signedURL, err := bucket.SignURL(ossKey, alioss.HTTPGet, int64(expireSeconds))
	if err != nil {
		return "", fmt.Errorf("生成签名 URL 失败: %v", err)
	}

	zlog.Info(fmt.Sprintf("生成下载 URL 成功: key=%s", ossKey))
	return signedURL, nil
}

// ValidateOSSKey 验证 OSS Key 格式
func (s *ossService) ValidateOSSKey(ossKey string) bool {
	cfg := config.GetConfig().OSSConfig
	// 确保 key 以加密文件前缀开头
	return strings.HasPrefix(ossKey, cfg.EncryptedFilePrefix)
}
