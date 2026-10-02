// Package route 负责注册所有 HTTP 路由。
package route

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mini-im/api/handler"
	"mini-im/api/middleware"
	"mini-im/conf"
	"mini-im/ws"
)

// Register 注册静态资源、REST 接口与 WebSocket 路由
func Register(r *gin.Engine, hub *ws.Hub) {
	webDir := conf.C.App.WebDir

	// 前端静态资源禁用缓存：否则移动端浏览器（尤其 iOS Safari）会长期沿用旧版
	// css/js，导致改动不生效、需手动清缓存
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/web") {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		c.Next()
	})

	// 静态页面与动态头像
	r.StaticFS("/web", http.Dir(webDir))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/index.html")
	})
	r.GET("/avatar", handler.Avatar)
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"connCount": hub.ConnCount(),
		})
	})

	// WebSocket
	r.GET("/ws", ws.Handler(hub))

	// REST API
	api := r.Group("/api")
	{
		api.POST("/login", handler.Login)
		api.POST("/register", handler.Register)

		authorized := api.Group("")
		authorized.Use(middleware.Auth())
		{
			authorized.PUT("/token", handler.Token)
			authorized.POST("/logout", handler.Logout)
			authorized.GET("/friends", handler.Friends)
			authorized.DELETE("/friends/:FID", handler.DeleteFriend)
			authorized.GET("/messages", handler.Messages)
			authorized.GET("/chats/:FID", handler.Chats)
			authorized.POST("/read/:FID", handler.Read)

			// 加好友：搜索用户 → 发起申请 → 对方同意 / 拒绝
			authorized.GET("/users", handler.SearchUsers)
			authorized.GET("/friend-requests", handler.FriendRequests)
			authorized.POST("/friend-requests", handler.AddFriend)
			authorized.POST("/friend-requests/:ID/accept", handler.AcceptFriendRequest)
			authorized.POST("/friend-requests/:ID/reject", handler.RejectFriendRequest)

			// 黑名单：拉入 / 移出 / 列表
			authorized.GET("/blacklist", handler.Blacklists)
			authorized.POST("/blacklist", handler.AddBlacklist)
			authorized.DELETE("/blacklist/:UID", handler.RemoveBlacklist)
		}
	}
}
