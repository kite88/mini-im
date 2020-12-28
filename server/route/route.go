package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"miniIM/server/app/api/apiImpl"
	"miniIM/server/app/api/middleware"
	"miniIM/server/app/ws"
	"net/http"
	"time"
)

func Web(r *gin.Engine) {
	fmt.Println("加载web静态资源", time.Now().Format("2006-01-02 15:04:05"))
	r.StaticFS("/web", http.Dir("web"))
	r.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/web"
		r.HandleContext(c)
	})
}

func Api(r *gin.Engine) {
	fmt.Println("加载api路由", time.Now().Format("2006-01-02 15:04:05"))
	api := r.Group("/api")
	{
		api.POST("/login", apiImpl.Login)
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware)
		{
			authorized.PUT("/token", apiImpl.Token)
			authorized.POST("/logout", apiImpl.Logout)
			authorized.GET("/friends", apiImpl.Friends)
			authorized.GET("/messages", apiImpl.Messages)
			authorized.GET("/chats/:FID", apiImpl.Chats)
		}
	}
}

func Ws(r *gin.Engine) {
	fmt.Println("websocket启动", time.Now().Format("2006-01-02 15:04:05"))
	r.GET("/ws", ws.WsHandler)
}
