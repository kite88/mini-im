package ApiUtil

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SuccessCode      int    = 0
	SuccessMsg       string = "请求成功"
	FailCode         int    = -1
	FailMsg          string = "致命错误"
	UnauthorizedCode int    = -10
	UnauthorizedMsg  string = "暂无权限，请先登录"
	UnauthorizedSMsg string = "账号在别处登录，尝试重新登录"
)

type NoneObject struct{}
type NoneArray []struct{}

func JSON(Data interface{}, Message string, StatusCode int, Context *gin.Context) {
	if Data == nil {
		Data = NoneObject{}
	}
	Context.JSON(http.StatusOK, gin.H{
		"StatusCode": StatusCode,
		"Message":    Message,
		"Data":       Data,
	})
	Context.Abort()
}
