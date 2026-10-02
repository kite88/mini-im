package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"mini-im/api/middleware"
	"mini-im/pkg/response"
	"mini-im/service"
)

// blacklistForm 拉入黑名单的请求体。
// user_id 是雪花 ID，按字符串传递，避免前端精度丢失（与其它接口一致）。
type blacklistForm struct {
	UserID int64 `json:"user_id,string" binding:"required"`
}

// Blacklists GET /api/blacklist
// 我拉黑的用户列表
func Blacklists(c *gin.Context) {
	list, err := service.Blacklist(middleware.UID(c))
	if err != nil {
		response.Fail(c, "获取黑名单失败")
		return
	}
	if list == nil {
		response.OK(c, response.EmptyList{})
		return
	}
	response.OK(c, list)
}

// AddBlacklist POST /api/blacklist
// 把指定用户拉入黑名单
func AddBlacklist(c *gin.Context) {
	var form blacklistForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	err := service.AddBlacklist(middleware.UID(c), form.UserID)
	switch {
	case err == nil:
		response.OKMsg(c, "已拉入黑名单", nil)
	case errors.Is(err, service.ErrBlacklistSelf),
		errors.Is(err, service.ErrUserNotFound):
		response.Fail(c, err.Error())
	default:
		response.Fail(c, "拉入黑名单失败，请稍后重试")
	}
}

// RemoveBlacklist DELETE /api/blacklist/:UID
// 把指定用户移出黑名单
func RemoveBlacklist(c *gin.Context) {
	targetID, err := strconv.ParseInt(strings.TrimSpace(c.Param("UID")), 10, 64)
	if err != nil || targetID <= 0 {
		response.Fail(c, "参数错误")
		return
	}

	err = service.RemoveBlacklist(middleware.UID(c), targetID)
	switch {
	case err == nil:
		response.OKMsg(c, "已移出黑名单", nil)
	case errors.Is(err, service.ErrBlacklistNotFound):
		response.Fail(c, err.Error())
	default:
		response.Fail(c, "移出黑名单失败，请稍后重试")
	}
}
