package apiImpl

import (
	"github.com/gin-gonic/gin"
	"miniIM/server/app/service"
	"miniIM/server/util/ApiUtil"
	"strconv"
)

type ChatApi struct {
	BaseApi
}

func Friends(c *gin.Context) {
	list := service.FriendList(c.GetInt64("uid"))
	ApiUtil.JSON(list, ApiUtil.SuccessMsg, ApiUtil.SuccessCode, c)
}

func Messages(c *gin.Context) {
	var data interface{}
	var list = service.MessageList(c.GetInt64("uid"))
	if list == nil {
		data = ApiUtil.NoneArray{}
	} else {
		data = list
	}
	ApiUtil.JSON(data, ApiUtil.SuccessMsg, ApiUtil.SuccessCode, c)
}

func Chats(c *gin.Context) {
	FID, _ := strconv.Atoi(c.Param("FID"))
	list := service.ChatList(c.GetInt64("uid"), int64(FID))
	ApiUtil.JSON(list, ApiUtil.SuccessMsg, ApiUtil.SuccessCode, c)
}
