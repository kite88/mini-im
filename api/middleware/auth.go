// Package middleware 提供 gin 中间件。
package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"

	"mini-im/pkg/httpx"
	"mini-im/pkg/response"
	"mini-im/service"
)

// CtxUIDKey gin.Context 中存放当前登录用户 ID 的键
const CtxUIDKey = "uid"

// Auth 登录态校验：解析 Authorization 头中的 token，校验通过后写入上下文
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := httpx.Bearer(c.GetHeader("Authorization"))
		if token == "" {
			response.Unauthorized(c, response.UnauthorizedMsg)
			return
		}

		uid, err := service.VerifyToken(token)
		if err != nil {
			msg := response.UnauthorizedMsg
			if errors.Is(err, service.ErrTokenReplaced) {
				msg = response.OtherLoginMsg
			}
			response.Unauthorized(c, msg)
			return
		}

		c.Set(CtxUIDKey, uid)
		c.Next()
	}
}

// UID 从上下文中取出当前登录用户 ID
func UID(c *gin.Context) int64 {
	value, ok := c.Get(CtxUIDKey)
	if !ok {
		return 0
	}
	uid, _ := value.(int64)
	return uid
}
