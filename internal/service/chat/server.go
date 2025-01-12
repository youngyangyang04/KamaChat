package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/enum/message/message_status_enum"
	"kama_chat_server/pkg/enum/message/message_type_enum"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
	"log"
	"sync"
	"time"
)

type Server struct {
	Clients  map[string]*Client
	mutex    *sync.Mutex
	Transmit chan []byte  // 转发通道
	Login    chan *Client // 登录通道
	Logout   chan *Client // 退出登录通道
}

var ChatServer *Server

func init() {
	if ChatServer == nil {
		ChatServer = &Server{
			Clients:  make(map[string]*Client),
			mutex:    &sync.Mutex{},
			Transmit: make(chan []byte, constants.CHANNEL_SIZE),
			Login:    make(chan *Client, constants.CHANNEL_SIZE),
			Logout:   make(chan *Client, constants.CHANNEL_SIZE),
		}
	}
}

// Start 启动函数，Server端用主进程起，Client端可以用协程起
func (s *Server) Start() {
	defer func() {
		close(s.Transmit)
		close(s.Logout)
		close(s.Login)
	}()
	for {
		select {
		case client := <-s.Login:
			{
				s.mutex.Lock()
				s.Clients[client.Uuid] = client
				s.mutex.Unlock()
				zlog.Debug(fmt.Sprintf("欢迎来到kama聊天服务器，亲爱的用户%s\n", client.Uuid))
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte("欢迎来到kama聊天服务器"))
				if err != nil {
					zlog.Error(err.Error())
				}
			}

		case client := <-s.Logout:
			{
				s.mutex.Lock()
				delete(s.Clients, client.Uuid)
				s.mutex.Unlock()
				zlog.Info(fmt.Sprintf("用户%s退出登录\n", client.Uuid))
			}

		case data := <-s.Transmit:
			{
				var chatMessageReq request.ChatMessageRequest
				if err := json.Unmarshal(data, &chatMessageReq); err != nil {
					zlog.Error(err.Error())
				}
				log.Println("原消息为：", data, "反序列化后为：", chatMessageReq)
				if chatMessageReq.Type == message_type_enum.Text {
					// 存message
					message := model.Message{
						Uuid:      fmt.Sprintf("M%s", random.GetNowAndLenRandomString(11)),
						SessionId: chatMessageReq.SessionId,
						Type:      chatMessageReq.Type,
						Content:   chatMessageReq.Content,
						Url:       "",
						SendId:    chatMessageReq.SendId,
						SendName:  chatMessageReq.SendName,
						ReceiveId: chatMessageReq.ReceiveId,
						FileSize:  "0B",
						FileType:  "",
						FileName:  "",
						Status:    message_status_enum.Unsent,
						CreatedAt: time.Now(),
					}
					if res := dao.GormDB.Create(&message); res.Error != nil {
						zlog.Error(res.Error.Error())
					}
					if message.ReceiveId[0] == 'U' { // 发送给User
						// 如果能找到ReceiveId，说明在线，可以发送，否则存表后跳过
						// 因为在线的时候是通过websocket更新消息记录的，离线后通过存表，登录时只调用一次数据库操作
						// 切换chat对象后，前端的messageList也会改变，获取messageList从第二次就是从redis中获取
						messageRsp := respond.GetMessageListRespond{
							SendId:    message.SendId,
							SendName:  message.SendName,
							ReceiveId: message.ReceiveId,
							Type:      message.Type,
							Content:   message.Content,
							Url:       message.Url,
							FileSize:  message.FileSize,
							FileName:  message.FileName,
							FileType:  message.FileType,
							CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
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
						if receiveClient, ok := s.Clients[message.ReceiveId]; ok {
							//messageBack.Message = jsonMessage
							//messageBack.Uuid = message.Uuid
							receiveClient.SendBack <- messageBack // 向client.Send发送
						}
						// 因为send_id肯定在线，所以这里在后端进行在线回显message，其实优化的话前端可以直接回显
						// 问题在于前后端的req和rsp结构不同，前端存储message的messageList不能存req，只能存rsp
						// 所以这里后端进行回显，前端不回显
						sendClient := s.Clients[message.SendId]
						sendClient.SendBack <- messageBack
					} else if message.ReceiveId[0] == 'G' { // 发送给Group
						messageRsp := respond.GetMessageListRespond{
							SendId:    message.SendId,
							SendName:  message.SendName,
							ReceiveId: message.ReceiveId,
							Type:      message.Type,
							Content:   message.Content,
							Url:       message.Url,
							FileSize:  message.FileSize,
							FileName:  message.FileName,
							FileType:  message.FileType,
							CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
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
						for _, member := range members {
							if member != message.SendId {
								if receiveClient, ok := s.Clients[member]; ok {
									receiveClient.SendBack <- messageBack
								}
							} else {
								sendClient := s.Clients[message.SendId]
								sendClient.SendBack <- messageBack
							}
						}
					}
				} else if chatMessageReq.Type == message_type_enum.File {
					// 存message
					message := model.Message{
						Uuid:      fmt.Sprintf("M%s", random.GetNowAndLenRandomString(11)),
						SessionId: chatMessageReq.SessionId,
						Type:      chatMessageReq.Type,
						Content:   "",
						Url:       chatMessageReq.Url,
						SendId:    chatMessageReq.SendId,
						SendName:  chatMessageReq.SendName,
						ReceiveId: chatMessageReq.ReceiveId,
						FileSize:  chatMessageReq.FileSize,
						FileType:  chatMessageReq.FileType,
						FileName:  chatMessageReq.FileName,
						Status:    message_status_enum.Unsent,
						CreatedAt: time.Now(),
					}
					if res := dao.GormDB.Create(&message); res.Error != nil {
						zlog.Error(res.Error.Error())
					}
					if message.ReceiveId[0] == 'U' { // 发送给User
						// 如果能找到ReceiveId，说明在线，可以发送，否则存表后跳过
						// 因为在线的时候是通过websocket更新消息记录的，离线后通过存表，登录时只调用一次数据库操作
						// 切换chat对象后，前端的messageList也会改变，获取messageList从第二次就是从redis中获取
						messageRsp := respond.GetMessageListRespond{
							SendId:    message.SendId,
							SendName:  message.SendName,
							ReceiveId: message.ReceiveId,
							Type:      message.Type,
							Content:   message.Content,
							Url:       message.Url,
							FileSize:  message.FileSize,
							FileName:  message.FileName,
							FileType:  message.FileType,
							CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
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
						if receiveClient, ok := s.Clients[message.ReceiveId]; ok {
							//messageBack.Message = jsonMessage
							//messageBack.Uuid = message.Uuid
							receiveClient.SendBack <- messageBack // 向client.Send发送
						}
						// 因为send_id肯定在线，所以这里在后端进行在线回显message，其实优化的话前端可以直接回显
						// 问题在于前后端的req和rsp结构不同，前端存储message的messageList不能存req，只能存rsp
						// 所以这里后端进行回显，前端不回显
						sendClient := s.Clients[message.SendId]
						sendClient.SendBack <- messageBack
					}
				}

			}
		}
	}
}

func (s *Server) SendClientToLogin(client *Client) {
	s.mutex.Lock()
	s.Login <- client
	s.mutex.Unlock()
}

func (s *Server) SendClientToLogout(client *Client) {
	s.mutex.Lock()
	s.Logout <- client
	s.mutex.Unlock()
}

func (s *Server) SendMessageToTransmit(message []byte) {
	s.mutex.Lock()
	s.Transmit <- message
	s.mutex.Unlock()
}

func (s *Server) RemoveClient(uuid string) {
	s.mutex.Lock()
	delete(s.Clients, uuid)
	s.mutex.Unlock()
}
