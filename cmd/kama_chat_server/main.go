package main

import (
	"fmt"
	"kama_chat_server/internal/config"
	"kama_chat_server/internal/https_server"
	"kama_chat_server/internal/service/chat"
	"kama_chat_server/internal/service/kafka"
	myredis "kama_chat_server/internal/service/redis"
	"kama_chat_server/pkg/zlog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port
	kafkaConfig := conf.KafkaConfig

	if kafkaConfig.MessageMode == "kafka" {
		kafka.KafkaService.KafkaInit()
	}

	if kafkaConfig.MessageMode == "channel" {
		go chat.ChatServer.Start()
	} else {
		go chat.KafkaChatServer.Start()
	}

	go func() {
		// 根据配置决定使用 HTTP 还是 HTTPS
		if conf.MainConfig.EnableHTTPS {
			zlog.Info(fmt.Sprintf("启动 HTTPS 服务器在 %s:%d", host, port))
			if err := https_server.GE.RunTLS(fmt.Sprintf("%s:%d", host, port), conf.MainConfig.CertFile, conf.MainConfig.KeyFile); err != nil {
				zlog.Fatal("HTTPS server running fault: " + err.Error())
				return
			}
		} else {
			zlog.Info(fmt.Sprintf("启动 HTTP 服务器在 %s:%d", host, port))
			if err := https_server.GE.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
				zlog.Fatal("HTTP server running fault: " + err.Error())
				return
			}
		}
	}()

	// 设置信号监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-quit

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
