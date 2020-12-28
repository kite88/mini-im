package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"miniIM/server/app/model"
	"miniIM/server/app/service"
	"miniIM/server/conf"
	"net/http"
	"strconv"
	"time"
)

type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

type Message struct {
	Sender      string `json:"sender"`
	Recipient   string `json:"recipient"`
	Content     string `json:"content"`
	Time        string `json:"time"`
	ContentType string `json:"content_type"`
	Subjoin     struct {
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
	} `json:"subjoin"`
}

var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
}

// Start is  项目运行前, 协程开启start -> go Manager.Start()
func (manager *ClientManager) Start() {
	for {
		log.Println("<---管道通信--->")
		select {
		case conn := <-Manager.Register:
			log.Printf("新用户加入:%v", conn.ID)
			Manager.Clients[conn.ID] = conn
			jsonMessage, _ := json.Marshal(&Message{Content: "socket服务连接成功"})
			conn.Send <- jsonMessage
			service.LineUser(conn.ID, true)
			offlineMsgList := service.ReadOfflineMsg(conn.ID)
			for _, value := range offlineMsgList {
				time.Sleep(500 * time.Millisecond)
				conn.Send <- []byte(value)
			}
		case conn := <-Manager.Unregister:
			log.Printf("用户离开:%v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				jsonMessage, _ := json.Marshal(&Message{Content: "socket已断开"})
				conn.Send <- jsonMessage
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
			service.LineUser(conn.ID, false)
		case message := <-Manager.Broadcast:
			MessageStruct := Message{}
			json.Unmarshal(message, &MessageStruct)
			for id, conn := range Manager.Clients {
				if id != creatId(MessageStruct.Recipient) {
					continue
				}
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
		}
	}
}
func creatId(uid string) string {
	return conf.ImUserPrefix + uid
}
func (c *Client) Read(uid string) {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		var messageData Message
		json.Unmarshal([]byte(string(message)), &messageData)
		messageData.Sender = uid
		messageData.Time = strconv.FormatInt(time.Now().Unix(), 10)
		jsonMessage, _ := json.Marshal(&messageData)
		log.Printf("读取到客户端的信息:%s", string(jsonMessage))
		Manager.Broadcast <- jsonMessage
		//存入队列里
		uidInt64, _ := strconv.ParseInt(messageData.Sender, 10, 64)
		fidInt64, _ := strconv.ParseInt(messageData.Recipient, 10, 64)
		timeInt64, _ := strconv.ParseInt(messageData.Time, 10, 64)
		contentTypeInt64, _ := strconv.ParseInt(messageData.ContentType, 10, 64)
		mModel := model.MessageModel{
			ContentType: contentTypeInt64,
			Content:     messageData.Content,
			CreateTime:  timeInt64,
			UserId:      uidInt64,
			FriendId:    fidInt64,
		}
		service.SaveToList(mModel)
		//如果离线 存到离线列表
		service.SaveOfflineMsgList(creatId(messageData.Recipient), string(jsonMessage))
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("发送到到客户端的信息:%s", string(message))
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

//TestHandler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func WsHandler(c *gin.Context) {
	uidInt64, errMsg := service.TokenVerify(c.Query("token"))
	if uidInt64 < 1 {
		log.Println("ws token errMsg", errMsg)
		return
	}
	uid := strconv.FormatInt(uidInt64, 10)
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	client := &Client{
		ID:     creatId(uid),
		Socket: conn,
		Send:   make(chan []byte),
	}
	Manager.Register <- client
	go client.Read(uid)
	go client.Write()
}
