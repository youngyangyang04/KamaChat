package config

import (
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type MainConfig struct {
	AppName     string `toml:"appName"`
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	EnableHTTPS bool   `toml:"enableHTTPS"`
	CertFile    string `toml:"certFile"`
	KeyFile     string `toml:"keyFile"`
}

type PostgresqlConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"databaseName"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`
}

type LogConfig struct {
	LogPath string `toml:"logPath"`
}

type KafkaConfig struct {
	MessageMode string        `toml:"messageMode"`
	HostPort    string        `toml:"hostPort"`
	LoginTopic  string        `toml:"loginTopic"`
	LogoutTopic string        `toml:"logoutTopic"`
	ChatTopic   string        `toml:"chatTopic"`
	Partition   int           `toml:"partition"`
	Timeout     time.Duration `toml:"timeout"`
}

type StaticSrcConfig struct {
	StaticAvatarPath string `toml:"staticAvatarPath"`
	StaticFilePath   string `toml:"staticFilePath"`
}

type OSSConfig struct {
	AccessKeyID           string `toml:"accessKeyID"`
	AccessKeySecret       string `toml:"accessKeySecret"`
	STSRoleArn            string `toml:"stsRoleArn"`
	Endpoint              string `toml:"endpoint"`
	Bucket                string `toml:"bucket"`
	Region                string `toml:"region"`
	UploadExpireSeconds   int64  `toml:"uploadExpireSeconds"`
	DownloadExpireSeconds int64  `toml:"downloadExpireSeconds"`
	MaxFileSize           int64  `toml:"maxFileSize"`
	EncryptedFilePrefix   string `toml:"encryptedFilePrefix"`
}

type JWTConfig struct {
	SecretKey    string `toml:"secretKey"`
	ExpireHours  int    `toml:"expireHours"`
	RefreshHours int    `toml:"refreshHours"`
}

type Config struct {
	MainConfig       `toml:"mainConfig"`
	PostgresqlConfig `toml:"postgresqlConfig"`
	RedisConfig      `toml:"redisConfig"`
	LogConfig        `toml:"logConfig"`
	KafkaConfig      `toml:"kafkaConfig"`
	StaticSrcConfig  `toml:"staticSrcConfig"`
	JWTConfig        `toml:"jwtConfig"`
	OSSConfig        `toml:"ossConfig"`
}

var config *Config

func LoadConfig() error {
	// 本地部署
	if _, err := toml.DecodeFile("configs/config.toml", config); err != nil {
		log.Fatal(err.Error())
		return err
	}
	// Ubuntu22.04云服务器部署
	// if _, err := toml.DecodeFile("/root/project/KamaChat/configs/config_local.toml", config); err != nil {
	// 	log.Fatal(err.Error())
	// 	return err
	// }

	// 从环境变量加载敏感配置（优先级高于配置文件）
	loadSensitiveConfigFromEnv()

	return nil
}

// loadSensitiveConfigFromEnv 从环境变量加载敏感配置
// 环境变量优先级高于配置文件，适用于开源项目避免泄露敏感信息
func loadSensitiveConfigFromEnv() {
	// OSS 配置
	if v := os.Getenv("OSS_ACCESS_KEY_ID"); v != "" {
		config.OSSConfig.AccessKeyID = v
	}
	if v := os.Getenv("OSS_ACCESS_KEY_SECRET"); v != "" {
		config.OSSConfig.AccessKeySecret = v
	}
	if v := os.Getenv("OSS_STS_ROLE_ARN"); v != "" {
		config.OSSConfig.STSRoleArn = v
	}
	if v := os.Getenv("OSS_ENDPOINT"); v != "" {
		config.OSSConfig.Endpoint = v
	}
	if v := os.Getenv("OSS_BUCKET"); v != "" {
		config.OSSConfig.Bucket = v
	}
	if v := os.Getenv("OSS_REGION"); v != "" {
		config.OSSConfig.Region = v
	}
}

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		_ = LoadConfig()
	}
	return config
}
