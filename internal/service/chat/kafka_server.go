package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gochat/internal/dao"
	"gochat/internal/dto/request"
	"gochat/internal/dto/respond"
	"gochat/internal/model"
	"gochat/internal/service/kafka"
	myredis "gochat/internal/service/redis"
	"gochat/pkg/constants"
	"gochat/pkg/enum/message/message_status_enum"
	"gochat/pkg/enum/message/message_type_enum"
	"gochat/pkg/util/random"
	"gochat/pkg/zlog"
	"log"
	"os"
	"sync"
	"time"
)

type KafkaServer struct {
	Clients map[string]*Client
	mutex   *sync.Mutex
	Login   chan *Client // 登录通道
	Logout  chan *Client // 退出登录通道
}

var KafkaChatServer *KafkaServer

var kafkaQuit = make(chan os.Signal, 1)

func init() {
	if KafkaChatServer == nil {
		KafkaChatServer = &KafkaServer{
			Clients: make(map[string]*Client),
			mutex:   &sync.Mutex{},
			Login:   make(chan *Client),
			Logout:  make(chan *Client),
		}
	}
	//signal.Notify(kafkaQuit, syscall.SIGINT, syscall.SIGTERM)
}

func (k *KafkaServer) Start() {
	defer func() {
		if r := recover(); r != nil {
			zlog.Error(fmt.Sprintf("kafka server panic: %v", r))
		}
		close(k.Login)
		close(k.Logout)
	}()

	// read chat message
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zlog.Error(fmt.Sprintf("kafka server panic: %v", r))
			}
		}()
		for {
			kafkaMessage, err := kafka.KafkaService.ChatReader.ReadMessage(ctx)
			if err != nil {
				zlog.Error(err.Error())
			}
			log.Printf("topic=%s, partition=%d, offset=%d, key=%s, value=%s", kafkaMessage.Topic, kafkaMessage.Partition, kafkaMessage.Offset, kafkaMessage.Key, kafkaMessage.Value)
			zlog.Info(fmt.Sprintf("topic=%s, partition=%d, offset=%d, key=%s, value=%s", kafkaMessage.Topic, kafkaMessage.Partition, kafkaMessage.Offset, kafkaMessage.Key, kafkaMessage.Value))
			data := kafkaMessage.Value
			var chatMessageReq request.ChatMessageRequest
			if err := json.Unmarshal(data, &chatMessageReq); err != nil {
				zlog.Error(err.Error())
			}
			log.Println("原消息为：", data, "反序列化后为：", chatMessageReq)
			if chatMessageReq.Type == message_type_enum.Text {
				// 存message
				message := model.Message{
					Uuid:       fmt.Sprintf("M%s", random.GetNowAndLenRandomString(11)),
					SessionId:  chatMessageReq.SessionId,
					Type:       chatMessageReq.Type,
					Content:    chatMessageReq.Content,
					Url:        "",
					SendId:     chatMessageReq.SendId,
					SendName:   chatMessageReq.SendName,
					SendAvatar: chatMessageReq.SendAvatar,
					ReceiveId:  chatMessageReq.ReceiveId,
					FileSize:   "0B",
					FileType:   "",
					FileName:   "",
					Status:     message_status_enum.Unsent,
					CreatedAt:  time.Now(),
					AVdata:     "",
				}
				// 对SendAvatar去除前面/static之前的所有内容，防止ip前缀引入
				message.SendAvatar = normalizePath(message.SendAvatar)
				if res := dao.GormDB.Create(&message); res.Error != nil {
					zlog.Error(res.Error.Error())
				}
				if message.ReceiveId[0] == 'U' { // 发送给User
					// 如果能找到ReceiveId，说明在线，可以发送，否则存表后跳过
					// 因为在线的时候是通过websocket更新消息记录的，离线后通过存表，登录时只调用一次数据库操作
					// 切换chat对象后，前端的messageList也会改变，获取messageList从第二次就是从redis中获取
					messageRsp := respond.GetMessageListRespond{
						SendId:     message.SendId,
						SendName:   message.SendName,
						SendAvatar: chatMessageReq.SendAvatar,
						ReceiveId:  message.ReceiveId,
						Type:       message.Type,
						Content:    message.Content,
						Url:        message.Url,
						FileSize:   message.FileSize,
						FileName:   message.FileName,
						FileType:   message.FileType,
						CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
					}
					jsonMessage, err := json.Marshal(messageRsp)
					if err != nil {
						zlog.Error(err.Error())
					}
					log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
					var messageBack = &MessageBack{
						Message: jsonMessage,
						Uuid:    message.Uuid,
					}
					k.mutex.Lock()
					if receiveClient, ok := k.Clients[message.ReceiveId]; ok {
						//messageBack.Message = jsonMessage
						//messageBack.Uuid = message.Uuid
						receiveClient.SendBack <- messageBack // 向client.Send发送
					}
					// 因为send_id肯定在线，所以这里在后端进行在线回显message，其实优化的话前端可以直接回显
					// 问题在于前后端的req和rsp结构不同，前端存储message的messageList不能存req，只能存rsp
					// 所以这里后端进行回显，前端不回显
					sendClient := k.Clients[message.SendId]
					sendClient.SendBack <- messageBack
					k.mutex.Unlock()

					// redis
					var rspString string
					rspString, err = myredis.GetKeyNilIsErr("message_list_" + message.SendId + "_" + message.ReceiveId)
					if err == nil {
						var rsp []respond.GetMessageListRespond
						if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
							zlog.Error(err.Error())
						}
						rsp = append(rsp, messageRsp)
						rspByte, err := json.Marshal(rsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						if err := myredis.SetKeyEx("message_list_"+message.SendId+"_"+message.ReceiveId, string(rspByte), time.Minute*constants.REDIS_TIMEOUT); err != nil {
							zlog.Error(err.Error())
						}
					} else {
						if !errors.Is(err, redis.Nil) {
							zlog.Error(err.Error())
						}
					}

				} else if message.ReceiveId[0] == 'G' { // 发送给Group
					messageRsp := respond.GetGroupMessageListRespond{
						SendId:     message.SendId,
						SendName:   message.SendName,
						SendAvatar: chatMessageReq.SendAvatar,
						ReceiveId:  message.ReceiveId,
						Type:       message.Type,
						Content:    message.Content,
						Url:        message.Url,
						FileSize:   message.FileSize,
						FileName:   message.FileName,
						FileType:   message.FileType,
						CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
					}
					jsonMessage, err := json.Marshal(messageRsp)
					if err != nil {
						zlog.Error(err.Error())
					}
					log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
					var messageBack = &MessageBack{
						Message: jsonMessage,
						Uuid:    message.Uuid,
					}
					var group model.GroupInfo
					if res := dao.GormDB.Where("uuid = ?", message.ReceiveId).First(&group); res.Error != nil {
						zlog.Error(res.Error.Error())
					}
					var members []string
					if err := json.Unmarshal(group.Members, &members); err != nil {
						zlog.Error(err.Error())
					}
					k.mutex.Lock()
					for _, member := range members {
						if member != message.SendId {
							if receiveClient, ok := k.Clients[member]; ok {
								receiveClient.SendBack <- messageBack
							}
						} else {
							sendClient := k.Clients[message.SendId]
							sendClient.SendBack <- messageBack
						}
					}
					k.mutex.Unlock()

					// redis
					var rspString string
					rspString, err = myredis.GetKeyNilIsErr("group_messagelist_" + message.ReceiveId)
					if err == nil {
						var rsp []respond.GetGroupMessageListRespond
						if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
							zlog.Error(err.Error())
						}
						rsp = append(rsp, messageRsp)
						rspByte, err := json.Marshal(rsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						if err := myredis.SetKeyEx("group_messagelist_"+message.ReceiveId, string(rspByte), time.Minute*constants.REDIS_TIMEOUT); err != nil {
							zlog.Error(err.Error())
						}
					} else {
						if !errors.Is(err, redis.Nil) {
							zlog.Error(err.Error())
						}
					}
				}
			} else if chatMessageReq.Type == message_type_enum.File {
				// 存message
				message := model.Message{
					Uuid:       fmt.Sprintf("M%s", random.GetNowAndLenRandomString(11)),
					SessionId:  chatMessageReq.SessionId,
					Type:       chatMessageReq.Type,
					Content:    "",
					Url:        chatMessageReq.Url,
					SendId:     chatMessageReq.SendId,
					SendName:   chatMessageReq.SendName,
					SendAvatar: chatMessageReq.SendAvatar,
					ReceiveId:  chatMessageReq.ReceiveId,
					FileSize:   chatMessageReq.FileSize,
					FileType:   chatMessageReq.FileType,
					FileName:   chatMessageReq.FileName,
					Status:     message_status_enum.Unsent,
					CreatedAt:  time.Now(),
					AVdata:     "",
				}
				// 对SendAvatar去除前面/static之前的所有内容，防止ip前缀引入
				message.SendAvatar = normalizePath(message.SendAvatar)
				if res := dao.GormDB.Create(&message); res.Error != nil {
					zlog.Error(res.Error.Error())
				}
				if message.ReceiveId[0] == 'U' { // 发送给User
					// 如果能找到ReceiveId，说明在线，可以发送，否则存表后跳过
					// 因为在线的时候是通过websocket更新消息记录的，离线后通过存表，登录时只调用一次数据库操作
					// 切换chat对象后，前端的messageList也会改变，获取messageList从第二次就是从redis中获取
					messageRsp := respond.GetMessageListRespond{
						SendId:     message.SendId,
						SendName:   message.SendName,
						SendAvatar: chatMessageReq.SendAvatar,
						ReceiveId:  message.ReceiveId,
						Type:       message.Type,
						Content:    message.Content,
						Url:        message.Url,
						FileSize:   message.FileSize,
						FileName:   message.FileName,
						FileType:   message.FileType,
						CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
					}
					jsonMessage, err := json.Marshal(messageRsp)
					if err != nil {
						zlog.Error(err.Error())
					}
					log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
					var messageBack = &MessageBack{
						Message: jsonMessage,
						Uuid:    message.Uuid,
					}
					k.mutex.Lock()
					if receiveClient, ok := k.Clients[message.ReceiveId]; ok {
						//messageBack.Message = jsonMessage
						//messageBack.Uuid = message.Uuid
						receiveClient.SendBack <- messageBack // 向client.Send发送
					}
					// 因为send_id肯定在线，所以这里在后端进行在线回显message，其实优化的话前端可以直接回显
					// 问题在于前后端的req和rsp结构不同，前端存储message的messageList不能存req，只能存rsp
					// 所以这里后端进行回显，前端不回显
					sendClient := k.Clients[message.SendId]
					sendClient.SendBack <- messageBack
					k.mutex.Unlock()

					// redis
					var rspString string
					rspString, err = myredis.GetKeyNilIsErr("message_list_" + message.SendId + "_" + message.ReceiveId)
					if err == nil {
						var rsp []respond.GetMessageListRespond
						if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
							zlog.Error(err.Error())
						}
						rsp = append(rsp, messageRsp)
						rspByte, err := json.Marshal(rsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						if err := myredis.SetKeyEx("message_list_"+message.SendId+"_"+message.ReceiveId, string(rspByte), time.Minute*constants.REDIS_TIMEOUT); err != nil {
							zlog.Error(err.Error())
						}
					} else {
						if !errors.Is(err, redis.Nil) {
							zlog.Error(err.Error())
						}
					}
				} else {
					messageRsp := respond.GetGroupMessageListRespond{
						SendId:     message.SendId,
						SendName:   message.SendName,
						SendAvatar: chatMessageReq.SendAvatar,
						ReceiveId:  message.ReceiveId,
						Type:       message.Type,
						Content:    message.Content,
						Url:        message.Url,
						FileSize:   message.FileSize,
						FileName:   message.FileName,
						FileType:   message.FileType,
						CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
					}
					jsonMessage, err := json.Marshal(messageRsp)
					if err != nil {
						zlog.Error(err.Error())
					}
					log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
					var messageBack = &MessageBack{
						Message: jsonMessage,
						Uuid:    message.Uuid,
					}
					var group model.GroupInfo
					if res := dao.GormDB.Where("uuid = ?", message.ReceiveId).First(&group); res.Error != nil {
						zlog.Error(res.Error.Error())
					}
					var members []string
					if err := json.Unmarshal(group.Members, &members); err != nil {
						zlog.Error(err.Error())
					}
					k.mutex.Lock()
					for _, member := range members {
						if member != message.SendId {
							if receiveClient, ok := k.Clients[member]; ok {
								receiveClient.SendBack <- messageBack
							}
						} else {
							sendClient := k.Clients[message.SendId]
							sendClient.SendBack <- messageBack
						}
					}
					k.mutex.Unlock()

					// redis
					var rspString string
					rspString, err = myredis.GetKeyNilIsErr("group_messagelist_" + message.ReceiveId)
					if err == nil {
						var rsp []respond.GetGroupMessageListRespond
						if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
							zlog.Error(err.Error())
						}
						rsp = append(rsp, messageRsp)
						rspByte, err := json.Marshal(rsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						if err := myredis.SetKeyEx("group_messagelist_"+message.ReceiveId, string(rspByte), time.Minute*constants.REDIS_TIMEOUT); err != nil {
							zlog.Error(err.Error())
						}
					} else {
						if !errors.Is(err, redis.Nil) {
							zlog.Error(err.Error())
						}
					}
				}
			} else if chatMessageReq.Type == message_type_enum.AudioOrVideo {
				var avData request.AVData
				if err := json.Unmarshal([]byte(chatMessageReq.AVdata), &avData); err != nil {
					zlog.Error(err.Error())
				}
				//log.Println(avData)
				message := model.Message{
					Uuid:       fmt.Sprintf("M%s", random.GetNowAndLenRandomString(11)),
					SessionId:  chatMessageReq.SessionId,
					Type:       chatMessageReq.Type,
					Content:    "",
					Url:        "",
					SendId:     chatMessageReq.SendId,
					SendName:   chatMessageReq.SendName,
					SendAvatar: chatMessageReq.SendAvatar,
					ReceiveId:  chatMessageReq.ReceiveId,
					FileSize:   "",
					FileType:   "",
					FileName:   "",
					Status:     message_status_enum.Unsent,
					CreatedAt:  time.Now(),
					AVdata:     chatMessageReq.AVdata,
				}
				if avData.MessageId == "PROXY" && (avData.Type == "start_call" || avData.Type == "receive_call" || avData.Type == "reject_call") {
					// 存message
					// 对SendAvatar去除前面/static之前的所有内容，防止ip前缀引入
					message.SendAvatar = normalizePath(message.SendAvatar)
					if res := dao.GormDB.Create(&message); res.Error != nil {
						zlog.Error(res.Error.Error())
					}
				}

				if chatMessageReq.ReceiveId[0] == 'U' { // 发送给User
					// 如果能找到ReceiveId，说明在线，可以发送，否则存表后跳过
					// 因为在线的时候是通过websocket更新消息记录的，离线后通过存表，登录时只调用一次数据库操作
					// 切换chat对象后，前端的messageList也会改变，获取messageList从第二次就是从redis中获取
					messageRsp := respond.AVMessageRespond{
						SendId:     message.SendId,
						SendName:   message.SendName,
						SendAvatar: message.SendAvatar,
						ReceiveId:  message.ReceiveId,
						Type:       message.Type,
						Content:    message.Content,
						Url:        message.Url,
						FileSize:   message.FileSize,
						FileName:   message.FileName,
						FileType:   message.FileType,
						CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
						AVdata:     message.AVdata,
					}
					jsonMessage, err := json.Marshal(messageRsp)
					if err != nil {
						zlog.Error(err.Error())
					}
					// log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
					log.Println("返回的消息为：", messageRsp)
					var messageBack = &MessageBack{
						Message: jsonMessage,
						Uuid:    message.Uuid,
					}
					k.mutex.Lock()
					if receiveClient, ok := k.Clients[message.ReceiveId]; ok {
						//messageBack.Message = jsonMessage
						//messageBack.Uuid = message.Uuid
						receiveClient.SendBack <- messageBack // 向client.Send发送
					}
					// 通话这不能回显，发回去的话就会出现两个start_call。
					//sendClient := s.Clients[message.SendId]
					//sendClient.SendBack <- messageBack
					k.mutex.Unlock()
				}
			}
		}
	}()

	// login, logout message
	for {
		select {
		case client := <-k.Login:
			{
				k.mutex.Lock()
				k.Clients[client.Uuid] = client
				k.mutex.Unlock()
				zlog.Debug(fmt.Sprintf("欢迎来到GoChat聊天服务器，亲爱的用户%s\n", client.Uuid))
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte("欢迎来到GoChat聊天服务器"))
				if err != nil {
					zlog.Error(err.Error())
				}
			}

		case client := <-k.Logout:
			{
				k.mutex.Lock()
				delete(k.Clients, client.Uuid)
				k.mutex.Unlock()
				zlog.Info(fmt.Sprintf("用户%s退出登录\n", client.Uuid))
				if err := client.Conn.WriteMessage(websocket.TextMessage, []byte("已退出登录")); err != nil {
					zlog.Error(err.Error())
				}
			}
		}
	}
}

func (k *KafkaServer) Close() {
	close(k.Login)
	close(k.Logout)
}

func (k *KafkaServer) SendClientToLogin(client *Client) {
	k.mutex.Lock()
	k.Login <- client
	k.mutex.Unlock()
}

func (k *KafkaServer) SendClientToLogout(client *Client) {
	k.mutex.Lock()
	k.Logout <- client
	k.mutex.Unlock()
}

func (k *KafkaServer) RemoveClient(uuid string) {
	k.mutex.Lock()
	delete(k.Clients, uuid)
	k.mutex.Unlock()
}
