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

// addFriendForm 发起好友申请的请求体。
// to_id 是雪花 ID，按字符串传递，避免前端精度丢失（与接口出参保持一致）。
type addFriendForm struct {
	ToID    int64  `json:"to_id,string" binding:"required"`
	Message string `json:"message"`
}

// SearchUsers GET /api/users?keyword=xxx
// 按账号 / 昵称搜索用户，并标注与当前用户的关系（可添加 / 已是好友 / 申请中）
func SearchUsers(c *gin.Context) {
	list, err := service.SearchUsers(middleware.UID(c), c.Query("keyword"))
	if err != nil {
		response.Fail(c, "搜索用户失败")
		return
	}
	if list == nil {
		response.OK(c, response.EmptyList{})
		return
	}
	response.OK(c, list)
}

// FriendRequests GET /api/friend-requests
// 待我处理的好友申请
func FriendRequests(c *gin.Context) {
	list, err := service.FriendRequestList(middleware.UID(c))
	if err != nil {
		response.Fail(c, "获取好友申请失败")
		return
	}
	if list == nil {
		response.OK(c, response.EmptyList{})
		return
	}
	response.OK(c, list)
}

// AddFriend POST /api/friend-requests
// 发起好友申请；若对方此前已申请我，则直接互为好友
func AddFriend(c *gin.Context) {
	var form addFriendForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	becameFriend, err := service.SendFriendRequest(middleware.UID(c), form.ToID, strings.TrimSpace(form.Message))
	switch {
	case err == nil && becameFriend:
		// became_friend 让前端知道无需等待确认，直接刷新好友列表
		response.OKMsg(c, "你们已成为好友", gin.H{"became_friend": true})
	case err == nil:
		response.OKMsg(c, "好友申请已发送，等待对方同意", gin.H{"became_friend": false})
	case errors.Is(err, service.ErrUserNotFound),
		errors.Is(err, service.ErrFriendSelf),
		errors.Is(err, service.ErrAlreadyFriend),
		errors.Is(err, service.ErrBlacklistBlockedPeer),
		errors.Is(err, service.ErrBlacklistRejected):
		response.Fail(c, err.Error())
	default:
		response.Fail(c, "好友申请发送失败，请稍后重试")
	}
}

// DeleteFriend DELETE /api/friends/:FID
// 删除好友（软删除）：双方的好友列表都会移除该好友，聊天记录保留
func DeleteFriend(c *gin.Context) {
	fid, err := strconv.ParseInt(strings.TrimSpace(c.Param("FID")), 10, 64)
	if err != nil || fid <= 0 {
		response.Fail(c, "参数错误")
		return
	}

	err = service.DeleteFriend(middleware.UID(c), fid)
	switch {
	case err == nil:
		response.OKMsg(c, "已删除好友", nil)
	case errors.Is(err, service.ErrFriendNotFound):
		response.Fail(c, err.Error())
	default:
		response.Fail(c, "删除好友失败，请稍后重试")
	}
}

// AcceptFriendRequest POST /api/friend-requests/:ID/accept
func AcceptFriendRequest(c *gin.Context) { handleFriendRequest(c, true) }

// RejectFriendRequest POST /api/friend-requests/:ID/reject
func RejectFriendRequest(c *gin.Context) { handleFriendRequest(c, false) }

// handleFriendRequest 同意 / 拒绝申请，两者只是状态不同
func handleFriendRequest(c *gin.Context, accept bool) {
	id, err := strconv.ParseInt(strings.TrimSpace(c.Param("ID")), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, "参数错误")
		return
	}

	err = service.HandleFriendRequest(middleware.UID(c), id, accept)
	switch {
	case err == nil && accept:
		response.OKMsg(c, "已同意好友申请", nil)
	case err == nil:
		response.OKMsg(c, "已拒绝好友申请", nil)
	case errors.Is(err, service.ErrFriendRequestNotFound),
		errors.Is(err, service.ErrFriendRequestHandled):
		response.Fail(c, err.Error())
	default:
		response.Fail(c, "操作失败，请稍后重试")
	}
}
