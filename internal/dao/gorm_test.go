package dao

import (
	"testing"

	"gochat/internal/config"
)

func TestBuildDSNUsesTCPNetwork(t *testing.T) {
	t.Parallel()

	dsn := BuildDSN(config.MysqlConfig{
		Host:         "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "123456",
		DatabaseName: "gochat",
		Network:      "tcp",
	})

	want := "root:123456@tcp(127.0.0.1:3306)/gochat?charset=utf8mb4&parseTime=True&loc=Local"
	if dsn != want {
		t.Fatalf("expected tcp dsn %q, got %q", want, dsn)
	}
}

func TestBuildDSNUsesUnixSocket(t *testing.T) {
	t.Parallel()

	dsn := BuildDSN(config.MysqlConfig{
		User:         "root",
		Password:     "123456",
		DatabaseName: "gochat",
		Network:      "unix",
		SocketPath:   "/var/run/mysqld/mysqld.sock",
	})

	want := "root:123456@unix(/var/run/mysqld/mysqld.sock)/gochat?charset=utf8mb4&parseTime=True&loc=Local"
	if dsn != want {
		t.Fatalf("expected unix dsn %q, got %q", want, dsn)
	}
}
