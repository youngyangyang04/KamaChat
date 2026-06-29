package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	// 这一组字段控制 HTTP/HTTPS 服务本身怎么启动。
	AppName     string `toml:"appName"`
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	TLSEnabled  bool   `toml:"tlsEnabled"`
	SSLRedirect bool   `toml:"sslRedirect"`
	TLSCertFile string `toml:"tlsCertFile"`
	TLSKeyFile  string `toml:"tlsKeyFile"`
}

type MysqlConfig struct {
	// network=tcp 时走 host:port，network=unix 时走 socketPath。
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"databaseName"`
	Network      string `toml:"network"`
	SocketPath   string `toml:"socketPath"`
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

type Config struct {
	// Config 是运行时配置的总入口，TOML 文件里的各个 section
	// 会分别映射到下面这些子结构体里。
	ServerConfig    ServerConfig    `toml:"serverConfig"`
	MysqlConfig     MysqlConfig     `toml:"mysqlConfig"`
	RedisConfig     RedisConfig     `toml:"redisConfig"`
	AuthCodeConfig  AuthCodeConfig  `toml:"authCodeConfig"`
	LogConfig       LogConfig       `toml:"logConfig"`
	KafkaConfig     KafkaConfig     `toml:"kafkaConfig"`
	StaticSrcConfig StaticSrcConfig `toml:"staticSrcConfig"`
}

var (
	// config 保存“已经初始化好的全局配置单例”。
	// configMu 用来保护首次初始化，避免并发重复加载。
	config   *Config
	configMu sync.Mutex
)

func LoadConfig() (*Config, error) {
	// APP_ENV 决定读取哪份环境覆盖文件。
	// 例如 dev -> configs/config.dev.toml。
	appEnv := strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv == "" {
		appEnv = "dev"
	}

	// base 配置是必须存在的，应用所有环境都会先读它。
	basePath, err := resolveRequiredConfigPath(filepath.Join("configs", "config.toml"))
	if err != nil {
		return nil, err
	}

	// 环境配置是可选的；没有 config.<APP_ENV>.toml 也不报错。
	envPath, err := resolveOptionalConfigPath(filepath.Join("configs", fmt.Sprintf("config.%s.toml", appEnv)))
	if err != nil {
		return nil, err
	}

	// CONFIG_FILE 允许调用方再显式指定一份更高优先级的配置文件。
	explicitConfig := strings.TrimSpace(os.Getenv("CONFIG_FILE"))
	explicitPath := ""
	if explicitConfig != "" {
		explicitPath, err = resolveRequiredConfigPath(explicitConfig)
		if err != nil {
			return nil, err
		}
	}

	return loadConfigFromFiles(basePath, envPath, explicitPath, os.LookupEnv)
}

func GetConfig() *Config {
	// 这里做的是“懒加载的全局单例”：
	// 第一次调用时真正读配置，后面都直接复用同一份结果。
	configMu.Lock()
	defer configMu.Unlock()

	if config == nil {
		loadedConfig, err := LoadConfig()
		if err != nil {
			panic(err)
		}
		config = loadedConfig
	}

	return config
}

func loadConfigFromFiles(basePath string, envPath string, explicitPath string, lookupEnv func(string) (string, bool)) (*Config, error) {
	cfg := &Config{}

	// 同一个 cfg 会被连续多次 decode。
	// 后读到的字段会覆盖先前的值，因此可以实现“base + env + explicit”的层叠覆盖。
	if err := decodeTomlFile(basePath, cfg); err != nil {
		return nil, fmt.Errorf("load base config %q: %w", basePath, err)
	}

	if envPath != "" {
		if err := decodeTomlFile(envPath, cfg); err != nil {
			return nil, fmt.Errorf("load env config %q: %w", envPath, err)
		}
	}

	if explicitPath != "" {
		if err := decodeTomlFile(explicitPath, cfg); err != nil {
			return nil, fmt.Errorf("load explicit config %q: %w", explicitPath, err)
		}
	}

	if lookupEnv != nil {
		// 文件配置读完后，再让环境变量作为最终覆盖层。
		if err := applyEnvOverrides(cfg, lookupEnv); err != nil {
			return nil, err
		}
	}

	// 最后给缺失字段补默认值，保证程序在最少配置下也能启动。
	cfg.normalize()
	return cfg, nil
}

func decodeTomlFile(path string, cfg *Config) error {
	if path == "" {
		return nil
	}
	// toml.DecodeFile 会根据 struct tag（例如 `toml:"serverConfig"`）
	// 把文件内容映射进 cfg 对应的字段。
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return err
	}
	return nil
}

func resolveRequiredConfigPath(path string) (string, error) {
	// 必选配置找不到时直接报错，避免应用带着空配置启动。
	resolvedPath, err := resolveConfigPath(path)
	if err != nil {
		return "", err
	}
	if resolvedPath == "" {
		return "", fmt.Errorf("config file %q not found", path)
	}
	return resolvedPath, nil
}

func resolveOptionalConfigPath(path string) (string, error) {
	// 可选配置找不到时返回空字符串，由调用方自己决定是否跳过。
	return resolveConfigPath(path)
}

func resolveConfigPath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", nil
	}

	if filepath.IsAbs(path) {
		if fileExists(path) {
			return filepath.Clean(path), nil
		}
		return "", nil
	}

	for _, baseDir := range candidateConfigBaseDirs() {
		// 对相对路径，分别尝试“当前工作目录”和“可执行文件目录”
		// 附近的几个候选位置，尽量兼容 go run / 二进制运行 / docker 运行。
		candidatePath := filepath.Clean(filepath.Join(baseDir, path))
		if fileExists(candidatePath) {
			return candidatePath, nil
		}
	}

	return "", nil
}

func candidateConfigBaseDirs() []string {
	// 之所以准备多组目录，是因为开发态和部署态的启动位置不一定一致。
	// 例如 go run 时 cwd 可能在仓库根目录，打包运行时则可能在可执行文件目录。
	dirs := make([]string, 0, 8)

	if workingDir, err := os.Getwd(); err == nil {
		dirs = append(dirs, workingDir, filepath.Join(workingDir, ".."), filepath.Join(workingDir, "..", ".."))
	}

	if executablePath, err := os.Executable(); err == nil {
		executableDir := filepath.Dir(executablePath)
		dirs = append(dirs, executableDir, filepath.Join(executableDir, ".."), filepath.Join(executableDir, "..", ".."))
	}

	seen := make(map[string]struct{}, len(dirs))
	uniqueDirs := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		// 去重，避免重复尝试同一路径。
		cleanDir := filepath.Clean(dir)
		if _, ok := seen[cleanDir]; ok {
			continue
		}
		seen[cleanDir] = struct{}{}
		uniqueDirs = append(uniqueDirs, cleanDir)
	}

	return uniqueDirs
}

func fileExists(path string) bool {
	// 这里只接受“存在且不是目录”的普通文件。
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func applyEnvOverrides(cfg *Config, lookupEnv func(string) (string, bool)) error {
	// 这里不是自动扫描结构体，而是“手写映射关系”：
	// 哪个环境变量覆盖哪个字段，都在下面显式列出来。
	type stringTarget struct {
		envName string
		target  *string
	}
	type intTarget struct {
		envName string
		target  *int
	}
	type boolTarget struct {
		envName string
		target  *bool
	}

	stringTargets := []stringTarget{
		{envName: "GOCHAT_DB_HOST", target: &cfg.MysqlConfig.Host},
		{envName: "GOCHAT_DB_USER", target: &cfg.MysqlConfig.User},
		{envName: "GOCHAT_DB_PASSWORD", target: &cfg.MysqlConfig.Password},
		{envName: "GOCHAT_DB_NAME", target: &cfg.MysqlConfig.DatabaseName},
		{envName: "GOCHAT_REDIS_HOST", target: &cfg.RedisConfig.Host},
		{envName: "GOCHAT_REDIS_PASSWORD", target: &cfg.RedisConfig.Password},
		{envName: "GOCHAT_TLS_CERT_FILE", target: &cfg.ServerConfig.TLSCertFile},
		{envName: "GOCHAT_TLS_KEY_FILE", target: &cfg.ServerConfig.TLSKeyFile},
	}

	for _, item := range stringTargets {
		if value, ok := lookupTrimmedEnv(lookupEnv, item.envName); ok {
			// 字符串字段可直接覆盖。
			*item.target = value
		}
	}

	intTargets := []intTarget{
		{envName: "GOCHAT_DB_PORT", target: &cfg.MysqlConfig.Port},
		{envName: "GOCHAT_REDIS_PORT", target: &cfg.RedisConfig.Port},
	}

	for _, item := range intTargets {
		if value, ok := lookupTrimmedEnv(lookupEnv, item.envName); ok {
			// 环境变量天然是字符串，所以端口这类字段要先做类型转换。
			parsedValue, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("parse %s: %w", item.envName, err)
			}
			*item.target = parsedValue
		}
	}

	boolTargets := []boolTarget{
		{envName: "GOCHAT_TLS_ENABLED", target: &cfg.ServerConfig.TLSEnabled},
		{envName: "GOCHAT_SSL_REDIRECT", target: &cfg.ServerConfig.SSLRedirect},
	}

	for _, item := range boolTargets {
		if value, ok := lookupTrimmedEnv(lookupEnv, item.envName); ok {
			// 布尔字段同理，需要把 "true"/"false" 之类的文本转成 bool。
			parsedValue, err := parseBool(value)
			if err != nil {
				return fmt.Errorf("parse %s: %w", item.envName, err)
			}
			*item.target = parsedValue
		}
	}

	return nil
}

func lookupTrimmedEnv(lookupEnv func(string) (string, bool), key string) (string, bool) {
	if lookupEnv == nil {
		return "", false
	}

	value, ok := lookupEnv(key)
	if !ok {
		return "", false
	}
	// 空字符串在这里按“没有设置”处理，避免把已有配置覆盖成空值。
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	return value, true
}

func parseBool(value string) (bool, error) {
	// 允许多种常见写法，方便在 shell、docker、CI 里传值。
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, errors.New("invalid boolean value")
	}
}

func (cfg *Config) normalize() {
	// normalize 只做一件事：给缺省字段补默认值。
	// 它不会覆盖“已经从 TOML 或环境变量里读到的值”。
	if cfg.ServerConfig.AppName == "" {
		cfg.ServerConfig.AppName = "gochat"
	}
	if cfg.ServerConfig.Host == "" {
		cfg.ServerConfig.Host = "0.0.0.0"
	}
	if cfg.ServerConfig.Port == 0 {
		cfg.ServerConfig.Port = 8000
	}

	if cfg.MysqlConfig.Host == "" {
		cfg.MysqlConfig.Host = "127.0.0.1"
	}
	if cfg.MysqlConfig.Port == 0 {
		cfg.MysqlConfig.Port = 3306
	}
	if cfg.MysqlConfig.User == "" {
		cfg.MysqlConfig.User = "root"
	}
	if cfg.MysqlConfig.DatabaseName == "" {
		cfg.MysqlConfig.DatabaseName = "gochat"
	}
	if cfg.MysqlConfig.Network == "" {
		cfg.MysqlConfig.Network = "tcp"
	}
	if cfg.MysqlConfig.SocketPath == "" {
		cfg.MysqlConfig.SocketPath = "/var/run/mysqld/mysqld.sock"
	}

	if cfg.RedisConfig.Host == "" {
		cfg.RedisConfig.Host = "127.0.0.1"
	}
	if cfg.RedisConfig.Port == 0 {
		cfg.RedisConfig.Port = 6379
	}

	if cfg.LogConfig.LogPath == "" {
		cfg.LogConfig.LogPath = filepath.ToSlash(filepath.Join(".", "logs", "gochat.log"))
	}

	if cfg.KafkaConfig.MessageMode == "" {
		cfg.KafkaConfig.MessageMode = "channel"
	}
	if cfg.KafkaConfig.HostPort == "" {
		cfg.KafkaConfig.HostPort = "127.0.0.1:9092"
	}
	if cfg.KafkaConfig.LoginTopic == "" {
		cfg.KafkaConfig.LoginTopic = "login"
	}
	if cfg.KafkaConfig.LogoutTopic == "" {
		cfg.KafkaConfig.LogoutTopic = "logout"
	}
	if cfg.KafkaConfig.ChatTopic == "" {
		cfg.KafkaConfig.ChatTopic = "chat_message"
	}
	if cfg.KafkaConfig.Timeout == 0 {
		cfg.KafkaConfig.Timeout = 1
	}

	if cfg.StaticSrcConfig.StaticAvatarPath == "" {
		cfg.StaticSrcConfig.StaticAvatarPath = "./static/avatars"
	}
	if cfg.StaticSrcConfig.StaticFilePath == "" {
		cfg.StaticSrcConfig.StaticFilePath = "./static/files"
	}
}
