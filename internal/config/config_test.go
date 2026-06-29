package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromFilesPrecedence(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	baseConfig := `
[serverConfig]
appName = "base-app"
host = "0.0.0.0"
port = 8000

[mysqlConfig]
host = "base-db"
port = 3306
user = "base-user"
password = "base-pass"
databaseName = "base-db-name"
network = "tcp"

[redisConfig]
host = "base-redis"
port = 6379
password = "base-redis-pass"
db = 0

[logConfig]
logPath = "./logs/base.log"
`

	envConfig := `
[serverConfig]
host = "env-host"

[mysqlConfig]
host = "env-db"
databaseName = "env-db-name"
`

	explicitConfig := `
[serverConfig]
port = 9443
tlsEnabled = true

[mysqlConfig]
host = "explicit-db"
`

	basePath := writeTestConfigFile(t, tempDir, "config.toml", baseConfig)
	envPath := writeTestConfigFile(t, tempDir, "config.dev.toml", envConfig)
	explicitPath := writeTestConfigFile(t, tempDir, "config.override.toml", explicitConfig)

	envValues := map[string]string{
		"GOCHAT_DB_HOST":        "env-var-db",
		"GOCHAT_DB_PORT":        "3307",
		"GOCHAT_DB_USER":        "env-var-user",
		"GOCHAT_DB_PASSWORD":    "env-var-pass",
		"GOCHAT_DB_NAME":        "env-var-db-name",
		"GOCHAT_REDIS_HOST":     "env-var-redis",
		"GOCHAT_REDIS_PORT":     "6380",
		"GOCHAT_REDIS_PASSWORD": "env-var-redis-pass",
		"GOCHAT_TLS_ENABLED":    "false",
		"GOCHAT_SSL_REDIRECT":   "true",
		"GOCHAT_TLS_CERT_FILE":  "/tmp/server.crt",
		"GOCHAT_TLS_KEY_FILE":   "/tmp/server.key",
	}

	cfg, err := loadConfigFromFiles(basePath, envPath, explicitPath, func(key string) (string, bool) {
		value, ok := envValues[key]
		return value, ok
	})
	if err != nil {
		t.Fatalf("loadConfigFromFiles() error = %v", err)
	}

	if cfg.ServerConfig.AppName != "base-app" {
		t.Fatalf("expected app name from base config, got %q", cfg.ServerConfig.AppName)
	}
	if cfg.ServerConfig.Host != "env-host" {
		t.Fatalf("expected host from env config, got %q", cfg.ServerConfig.Host)
	}
	if cfg.ServerConfig.Port != 9443 {
		t.Fatalf("expected port from explicit config, got %d", cfg.ServerConfig.Port)
	}
	if cfg.ServerConfig.TLSEnabled {
		t.Fatalf("expected env var to disable tls")
	}
	if !cfg.ServerConfig.SSLRedirect {
		t.Fatalf("expected env var to enable ssl redirect")
	}
	if cfg.ServerConfig.TLSCertFile != "/tmp/server.crt" {
		t.Fatalf("expected env cert path override, got %q", cfg.ServerConfig.TLSCertFile)
	}
	if cfg.ServerConfig.TLSKeyFile != "/tmp/server.key" {
		t.Fatalf("expected env key path override, got %q", cfg.ServerConfig.TLSKeyFile)
	}

	if cfg.MysqlConfig.Host != "env-var-db" {
		t.Fatalf("expected mysql host from env var, got %q", cfg.MysqlConfig.Host)
	}
	if cfg.MysqlConfig.Port != 3307 {
		t.Fatalf("expected mysql port from env var, got %d", cfg.MysqlConfig.Port)
	}
	if cfg.MysqlConfig.User != "env-var-user" {
		t.Fatalf("expected mysql user from env var, got %q", cfg.MysqlConfig.User)
	}
	if cfg.MysqlConfig.Password != "env-var-pass" {
		t.Fatalf("expected mysql password from env var, got %q", cfg.MysqlConfig.Password)
	}
	if cfg.MysqlConfig.DatabaseName != "env-var-db-name" {
		t.Fatalf("expected mysql db name from env var, got %q", cfg.MysqlConfig.DatabaseName)
	}

	if cfg.RedisConfig.Host != "env-var-redis" {
		t.Fatalf("expected redis host from env var, got %q", cfg.RedisConfig.Host)
	}
	if cfg.RedisConfig.Port != 6380 {
		t.Fatalf("expected redis port from env var, got %d", cfg.RedisConfig.Port)
	}
	if cfg.RedisConfig.Password != "env-var-redis-pass" {
		t.Fatalf("expected redis password from env var, got %q", cfg.RedisConfig.Password)
	}
}

func TestLoadConfigRejectsInvalidBooleanEnv(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	basePath := writeTestConfigFile(t, tempDir, "config.toml", `
[serverConfig]
appName = "gochat"
`)

	_, err := loadConfigFromFiles(basePath, "", "", func(key string) (string, bool) {
		if key == "GOCHAT_TLS_ENABLED" {
			return "maybe", true
		}
		return "", false
	})
	if err == nil {
		t.Fatal("expected invalid boolean env to fail")
	}
}

func writeTestConfigFile(t *testing.T, dir string, name string, contents string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write config file %s: %v", path, err)
	}
	return path
}
