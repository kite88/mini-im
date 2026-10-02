package ws

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"mini-im/model"
	"mini-im/service"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait * 9 / 10
	maxMessageSize = 8 << 10 // 8KB
	sendBufferSize = 64
)

// ClientMessage 上行/下行消息体，字段与前端保持一致
type ClientMessage struct {
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

// Client 单个 WebSocket 连接
type Client struct {
	uid  int64
	conn *websocket.Conn
	send chan []byte
	hub  *Hub

	once sync.Once
	done chan struct{}
}

func newClient(hub *Hub, uid int64, conn *websocket.Conn) *Client {
	return &Client{
		uid:  uid,
		conn: conn,
		send: make(chan []byte, sendBufferSize),
		hub:  hub,
		done: make(chan struct{}),
	}
}

// push 非阻塞投递，缓冲写满或连接已关闭时返回 false
func (c *Client) push(payload []byte) bool {
	select {
	case <-c.done:
		return false
	default:
	}

	select {
	case c.send <- payload:
		return true
	default:
		return false
	}
}

// shutdown 幂等关闭：注销连接、通知写协程退出
func (c *Client) shutdown() {
	c.once.Do(func() {
		close(c.done)
		c.hub.remove(c)
		_ = c.conn.Close()
	})
}

// flushPending 连接建立后补发欢迎语与离线消息
func (c *Client) flushPending() {
	welcome, err := json.Marshal(&ClientMessage{Content: "socket服务连接成功"})
	if err == nil {
		c.push(welcome)
	}
	for _, item := range service.DrainOffline(c.uid) {
		c.push([]byte(item))
	}
}

// readPump 读取客户端消息并分发，阻塞运行在请求协程中
func (c *Client) readPump() {
	defer c.shutdown()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var msg ClientMessage
		if err = json.Unmarshal(raw, &msg); err != nil {
			log.Printf("[ws] 消息格式错误: %v", err)
			continue
		}
		c.dispatch(&msg)
	}
}

// writePump 将消息写入连接，并定期发送 ping 保活
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.shutdown()
	}()

	for {
		select {
		case payload := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.done:
			return
		}
	}
}

// dispatch 处理一条上行消息：转发给接收方 + 投递落库队列 + 离线缓存
func (c *Client) dispatch(msg *ClientMessage) {
	recipient, err := strconv.ParseInt(strings.TrimSpace(msg.Recipient), 10, 64)
	if err != nil || recipient <= 0 || recipient == c.uid {
		return
	}

	// 只有好友之间可以发消息：任一方删除好友之后消息都不再可达
	if !service.IsFriend(c.uid, recipient) {
		event := service.EventNotFriend
		if service.RemovedByPeer(c.uid, recipient) {
			// 是对方删除了我：与「我删除了对方」用不同事件，前端提示更准确
			event = service.EventUnfriended
		}
		notifySendRejected(c.uid, recipient, event)
		return
	}

	// 黑名单：接收方已把我拉黑，消息不转发、不离线、不落库，只回执告知发送方
	if !service.CanSendMessage(c.uid, recipient) {
		notifySendRejected(c.uid, recipient, service.EventBlacklisted)
		return
	}

	// 内容类型以服务端归一化结果为准：未知标识回落为文本
	contentType := model.NormalizeContentType(msg.ContentType)
	now := time.Now().Unix()

	// 发送者信息以服务端为准，避免客户端伪造
	msg.Sender = strconv.FormatInt(c.uid, 10)
	msg.Time = strconv.FormatInt(now, 10)
	msg.ContentType = contentType
	if sender, err := service.GetUserBrief(c.uid); err == nil && sender != nil {
		msg.Subjoin.Avatar = sender.Avatar
		msg.Subjoin.Nickname = sender.Nickname
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}

	// 对方无活跃连接则写入离线队列
	if !c.hub.SendTo(recipient, payload) {
		service.PushOffline(recipient, payload)
	}

	// 异步落库
	if err = service.PushMessageQueue(model.Message{
		ContentType: contentType,
		Content:     msg.Content,
		UserID:      c.uid,
		FriendID:    recipient,
		CreateTime:  now,
	}); err != nil {
		log.Printf("[ws] 消息入队失败: %v", err)
	}
}

// notifySendRejected 告知发送方：消息未能送达（双方已不是好友 / 被对方拉黑）。
// 带上对方的资料，前端据此提示是「谁」拒收了消息。
func notifySendRejected(senderID, peerID int64, event string) {
	peer, err := service.GetUserBrief(peerID)
	if err != nil || peer == nil {
		return
	}
	service.Notify(senderID, event, map[string]interface{}{
		"id":       strconv.FormatInt(peer.ID, 10),
		"uid":      strconv.FormatInt(peer.ID, 10),
		"username": peer.Username,
		"nickname": peer.Nickname,
		"avatar":   peer.Avatar,
	})
}
