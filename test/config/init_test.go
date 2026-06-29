package config

import (
	"fmt"
	"gochat/internal/config"
	"testing"
)

func TestInit(t *testing.T) {
	conf := config.GetConfig()
	fmt.Println(conf.ServerConfig)
	fmt.Println(conf.MysqlConfig)
	fmt.Println(conf.RedisConfig)
	fmt.Println(conf.AuthCodeConfig)
}
