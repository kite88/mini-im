// Package handler 实现 HTTP 接口。
package handler

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"mini-im/api/middleware"
	"mini-im/pkg/response"
	"mini-im/service"
)

type loginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

// Login POST /api/login
func Login(c *gin.Context) {
	var form loginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	username := strings.TrimSpace(form.Username)
	password := strings.TrimSpace(form.Password)
	if username == "" || password == "" {
		response.Fail(c, "用户名或密码为空")
		return
	}

	user, token, err := service.Login(username, password)
	if err != nil {
		if errors.Is(err, service.ErrCredential) {
			response.Fail(c, err.Error())
		} else {
			response.Fail(c, "登录失败，请稍后重试")
		}
		return
	}

	response.OKMsg(c, "登录成功", gin.H{
		"user":   user,
		"tokens": token,
	})
}

// Register POST /api/register
func Register(c *gin.Context) {
	var form registerForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	user, err := service.Register(
		strings.TrimSpace(form.Username),
		strings.TrimSpace(form.Password),
		strings.TrimSpace(form.Nickname),
	)
	switch {
	case err == nil:
	case errors.Is(err, service.ErrUsernameInvalid),
		errors.Is(err, service.ErrPasswordTooShort),
		errors.Is(err, service.ErrUsernameExists):
		// 账号名 / 密码不合规则、账号已被占用：直接把原因回给前端
		response.Fail(c, err.Error())
		return
	default:
		response.Fail(c, "注册失败，请稍后重试")
		return
	}

	response.OKMsg(c, "注册成功", gin.H{"user": user})
}

// Token PUT /api/token
func Token(c *gin.Context) {
	token, err := service.RefreshToken(middleware.UID(c))
	if err != nil {
		response.Fail(c, "token 续签失败")
		return
	}
	response.OK(c, token)
}

// Logout POST /api/logout
func Logout(c *gin.Context) {
	if err := service.Logout(middleware.UID(c)); err != nil {
		response.Fail(c, "退出登录失败")
		return
	}
	response.OKMsg(c, "退出登录成功", nil)
}
