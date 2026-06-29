package chat

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"gochat/internal/config"
	"gochat/internal/dao"
	"gochat/internal/dto/request"
	"gochat/internal/model"
	myKafka "gochat/internal/service/kafka"
	"gochat/pkg/constants"
	"gochat/pkg/enum/message/message_status_enum"
	"gochat/pkg/zlog"
	"log"
	"net/http"
	"strconv"
)

type MessageBack struct {
	Message []byte
	Uuid    string
}

type Client struct {
	Conn     *websocket.Conn
	Uuid     string
	SendTo   chan []byte       // 给server端
	SendBack chan *MessageBack // 给前端
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	// 检查连接的Origin头
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ctx = context.Background()

var messageMode = config.GetConfig().KafkaConfig.MessageMode

// 读取websocket消息并发送给send通道
func (c *Client) Read() {
	zlog.Info("ws read goroutine start")
	for {
		// 阻塞有一定隐患，因为下面要处理缓冲的逻辑，但是可以先不做优化，问题不大
		_, jsonMessage, err := c.Conn.ReadMessage() // 阻塞状态
		if err != nil {
			zlog.Error(err.Error())
			return // 直接断开websocket
		} else {
			var message = request.ChatMessageRequest{}
			if err := json.Unmarshal(jsonMessage, &message); err != nil {
				zlog.Error(err.Error())
			}
			log.Println("接受到消息为: ", jsonMessage)
			if messageMode == "channel" {
				// 如果server的转发channel没满，先把sendto中的给transmit
				for len(ChatServer.Transmit) < constants.CHANNEL_SIZE && len(c.SendTo) > 0 {
					sendToMessage := <-c.SendTo
					ChatServer.SendMessageToTransmit(sendToMessage)
				}
				// 如果server没满，sendto空了，直接给server的transmit
				if len(ChatServer.Transmit) < constants.CHANNEL_SIZE {
					ChatServer.SendMessageToTransmit(jsonMessage)
				} else if len(c.SendTo) < constants.CHANNEL_SIZE {
					// 如果server满了，直接塞sendto
					c.SendTo <- jsonMessage
				} else {
					// 否则考虑加宽channel size，或者使用kafka
					if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("由于目前同一时间过多用户发送消息，消息发送失败，请稍后重试")); err != nil {
						zlog.Error(err.Error())
					}
				}
			} else {
				if err := myKafka.KafkaService.ChatWriter.WriteMessages(ctx, kafka.Message{
					Key:   []byte(strconv.Itoa(config.GetConfig().KafkaConfig.Partition)),
					Value: jsonMessage,
				}); err != nil {
					zlog.Error(err.Error())
				}
				zlog.Info("已发送消息：" + string(jsonMessage))
			}
		}
	}
}

// 从send通道读取消息发送给websocket
func (c *Client) Write() {
	zlog.Info("ws write goroutine start")
	for messageBack := range c.SendBack { // 阻塞状态
		// 通过 WebSocket 发送消息
		err := c.Conn.WriteMessage(websocket.TextMessage, messageBack.Message)
		if err != nil {
			zlog.Error(err.Error())
			return // 直接断开websocket
		}
		// log.Println("已发送消息：", messageBack.Message)
		// 说明顺利发送，修改状态为已发送
		if res := dao.GormDB.Model(&model.Message{}).Where("uuid = ?", messageBack.Uuid).Update("status", message_status_enum.Sent); res.Error != nil {
			zlog.Error(res.Error.Error())
		}
	}
}

// NewClientInit 当接受到前端有登录消息时，会调用该函数
func NewClientInit(c *gin.Context, clientId string) {
	kafkaConfig := config.GetConfig().KafkaConfig
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.Error(err.Error())
	}
	client := &Client{
		Conn:     conn,
		Uuid:     clientId,
		SendTo:   make(chan []byte, constants.CHANNEL_SIZE),
		SendBack: make(chan *MessageBack, constants.CHANNEL_SIZE),
	}
	if kafkaConfig.MessageMode == "channel" {
		ChatServer.SendClientToLogin(client)
	} else {
		KafkaChatServer.SendClientToLogin(client)
	}
	go client.Read()
	go client.Write()
	zlog.Info("ws连接成功")
}

// ClientLogout 当接受到前端有登出消息时，会调用该函数
func ClientLogout(clientId string) (string, int) {
	kafkaConfig := config.GetConfig().KafkaConfig
	client := ChatServer.Clients[clientId]
	if client != nil {
		if kafkaConfig.MessageMode == "channel" {
			ChatServer.SendClientToLogout(client)
		} else {
			KafkaChatServer.SendClientToLogout(client)
		}
		if err := client.Conn.Close(); err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}
		close(client.SendTo)
		close(client.SendBack)
	}
	return "退出成功", 0
}
