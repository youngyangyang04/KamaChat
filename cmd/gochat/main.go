package main

import (
	"fmt"
	"gochat/internal/config"
	"gochat/internal/dao"
	"gochat/internal/https_server"
	"gochat/internal/service/chat"
	"gochat/internal/service/kafka"
	myredis "gochat/internal/service/redis"
	"gochat/pkg/zlog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. 先把配置整体读出来，后面初始化数据库、HTTP 服务、聊天服务都会用到。
	conf := config.GetConfig()
	host := conf.ServerConfig.Host
	port := conf.ServerConfig.Port
	kafkaConfig := conf.KafkaConfig

	// 2. 初始化 Gorm，全局数据库连接会在这里建立，并自动同步表结构。
	dao.MustInit()
	if kafkaConfig.MessageMode == "kafka" {
		kafka.KafkaService.KafkaInit()
	}

	// 3. 启动消息分发服务。
	// channel 模式走内存 channel，kafka 模式走 broker。
	if kafkaConfig.MessageMode == "channel" {
		go chat.ChatServer.Start()
	} else {
		go chat.KafkaChatServer.Start()
	}

	// 4. 启动 Gin HTTP/HTTPS 服务。
	// Run / RunTLS 都是阻塞调用，所以放到 goroutine 里，后面才能继续监听退出信号。
	go func() {
		serverAddr := fmt.Sprintf("%s:%d", host, port)
		var err error
		if conf.ServerConfig.TLSEnabled {
			if conf.ServerConfig.TLSCertFile == "" || conf.ServerConfig.TLSKeyFile == "" {
				zlog.Fatal("tlsEnabled=true requires tlsCertFile and tlsKeyFile")
				return
			}
			err = https_server.GE.RunTLS(serverAddr, conf.ServerConfig.TLSCertFile, conf.ServerConfig.TLSKeyFile)
		} else {
			err = https_server.GE.Run(serverAddr)
		}
		if err != nil {
			zlog.Fatal("server running fault")
			return
		}
	}()

	// 设置信号监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-quit

	// 5. 进入关闭流程，按启动时相反的顺序回收资源。
	if kafkaConfig.MessageMode == "kafka" {
		kafka.KafkaService.KafkaClose()
	}

	chat.ChatServer.Close()

	zlog.Info("关闭服务器...")

	// 删除所有Redis键
	if err := myredis.DeleteAllRedisKeys(); err != nil {
		zlog.Error(err.Error())
	} else {
		zlog.Info("所有Redis键已删除")
	}

	zlog.Info("服务器已关闭")

}
