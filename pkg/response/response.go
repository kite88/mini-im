// Package response 统一接口返回格式。
//
// 响应体结构与前端约定保持一致：
//
//	{ "StatusCode": 0, "Message": "请求成功", "Data": {} }
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 业务状态码
const (
	SuccessCode      = 0
	SuccessMsg       = "请求成功"
	FailCode         = -1
	FailMsg          = "致命错误"
	UnauthorizedCode = -10
	UnauthorizedMsg  = "暂无权限，请先登录"
	OtherLoginMsg    = "账号在别处登录，尝试重新登录"
)

// Empty 空对象占位，保证 Data 始终是 JSON 对象而不是 null
type Empty struct{}

// EmptyList 空数组占位，保证列表接口返回 [] 而不是 null
type EmptyList []struct{}

// JSON 写出统一格式响应并终止后续处理
func JSON(c *gin.Context, code int, msg string, data interface{}) {
	if data == nil {
		data = Empty{}
	}
	c.JSON(http.StatusOK, gin.H{
		"StatusCode": code,
		"Message":    msg,
		"Data":       data,
	})
	c.Abort()
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	JSON(c, SuccessCode, SuccessMsg, data)
}

// OKMsg 带自定义提示的成功响应
func OKMsg(c *gin.Context, msg string, data interface{}) {
	JSON(c, SuccessCode, msg, data)
}

// Fail 业务失败响应
func Fail(c *gin.Context, msg string) {
	if msg == "" {
		msg = FailMsg
	}
	JSON(c, FailCode, msg, nil)
}

// Unauthorized 未授权 / 登录态失效响应
func Unauthorized(c *gin.Context, msg string) {
	if msg == "" {
		msg = UnauthorizedMsg
	}
	JSON(c, UnauthorizedCode, msg, nil)
}
