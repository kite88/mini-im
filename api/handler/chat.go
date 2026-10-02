package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mini-im/api/middleware"
	"mini-im/pkg/response"
	"mini-im/service"
)

// Friends GET /api/friends
func Friends(c *gin.Context) {
	list, err := service.FriendList(middleware.UID(c))
	if err != nil {
		response.Fail(c, "获取好友列表失败")
		return
	}
	response.OK(c, list)
}

// Messages GET /api/messages
func Messages(c *gin.Context) {
	list, err := service.MessageList(middleware.UID(c))
	if err != nil {
		response.Fail(c, "获取会话列表失败")
		return
	}
	if list == nil {
		response.OK(c, response.EmptyList{})
		return
	}
	response.OK(c, list)
}

// Chats GET /api/chats/:FID
func Chats(c *gin.Context) {
	fid, err := strconv.ParseInt(c.Param("FID"), 10, 64)
	if err != nil || fid <= 0 {
		response.Fail(c, "参数错误")
		return
	}

	list, err := service.ChatList(middleware.UID(c), fid)
	if err != nil {
		response.Fail(c, "获取聊天记录失败")
		return
	}
	if list == nil {
		response.OK(c, response.EmptyList{})
		return
	}
	response.OK(c, list)
}

// Read POST /api/read/:FID
// 把与指定好友的会话标记为已读（推进服务端已读游标），未读数由此清零
func Read(c *gin.Context) {
	fid, err := strconv.ParseInt(c.Param("FID"), 10, 64)
	if err != nil || fid <= 0 {
		response.Fail(c, "参数错误")
		return
	}

	if err = service.MarkConversationRead(middleware.UID(c), fid); err != nil {
		response.Fail(c, "标记已读失败")
		return
	}
	response.OKMsg(c, "标记已读成功", nil)
}
