package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"mini-im/pkg/httpx"
	"mini-im/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handler 返回 WebSocket 升级处理器
//
// 鉴权：/ws?token=xxx，也兼容 Authorization 头
func Handler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = httpx.Bearer(c.GetHeader("Authorization"))
		}

		uid, err := service.VerifyToken(token)
		if err != nil {
			log.Printf("[ws] 鉴权失败: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("[ws] 协议升级失败: %v", err)
			return
		}

		client := newClient(hub, uid, conn)
		hub.add(client)

		go client.writePump()
		client.flushPending()
		client.readPump()
	}
}
