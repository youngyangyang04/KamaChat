package config

import (
	"log"
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

type AuthCodeConfig struct {
	AccessKeyID     string `toml:"accessKeyID"`
	AccessKeySecret string `toml:"accessKeySecret"`
	SignName        string `toml:"signName"`
	TemplateCode    string `toml:"templateCode"`
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

type JWTConfig struct {
	SecretKey    string `toml:"secretKey"`
	ExpireHours  int    `toml:"expireHours"`
	RefreshHours int    `toml:"refreshHours"`
}

type Config struct {
	MainConfig       `toml:"mainConfig"`
	PostgresqlConfig `toml:"postgresqlConfig"`
	RedisConfig      `toml:"redisConfig"`
	AuthCodeConfig   `toml:"authCodeConfig"`
	LogConfig        `toml:"logConfig"`
	KafkaConfig      `toml:"kafkaConfig"`
	StaticSrcConfig  `toml:"staticSrcConfig"`
	JWTConfig        `toml:"jwtConfig"`
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
	return nil
}

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		_ = LoadConfig()
	}
	return config
}
