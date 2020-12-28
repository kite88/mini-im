package apiImpl

import (
	"github.com/gin-gonic/gin"
	"miniIM/server/app/service"
	"miniIM/server/util/ApiUtil"
	"strings"
)

type UserApi struct {
	BaseApi
}

type loginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	form := &loginForm{}
	if c.BindJSON(&form) != nil {
		ApiUtil.JSON(nil, "参数错误", ApiUtil.FailCode, c)
	} else if strings.Trim(form.Username, " ") == "" || strings.Trim(form.Password, " ") == "" {
		ApiUtil.JSON(nil, "用户名或密码为空", ApiUtil.FailCode, c)
	} else {
		r, data, te := service.LoginVerify(strings.Trim(form.Username, " "), strings.Trim(form.Password, " "))
		if r == false {
			ApiUtil.JSON(nil, "用户名或密码错误", ApiUtil.FailCode, c)
		} else {
			ApiUtil.JSON(map[string]interface{}{
				"user":   data,
				"tokens": te,
			}, "登录成功", ApiUtil.SuccessCode, c)
		}
	}
}

func Token(c *gin.Context) {
	uid := c.GetInt64("uid")
	te := service.TokenRefresh(uid)
	ApiUtil.JSON(te, ApiUtil.SuccessMsg, ApiUtil.SuccessCode, c)
}

func Logout(c *gin.Context) {
	uid := c.GetInt64("uid")
	service.Logout(uid)
	ApiUtil.JSON(nil, "退出登录成功", ApiUtil.SuccessCode, c)
}
