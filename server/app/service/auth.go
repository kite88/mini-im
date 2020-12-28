package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"miniIM/server/app/dao"
	"miniIM/server/app/model"
	"miniIM/server/conf"
	"miniIM/server/db"
	"miniIM/server/util/ApiUtil"
	"miniIM/server/util/JwtUtil"
	"strconv"
	"time"
)

var (
	AuthToken string = conf.RedisPrefix + "AUTH_TOKEN_"
)

func Register(username string, password string) bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return false
	}
	encodePWD := string(hash)
	result := dao.UserDao{}.D().Create(&model.UserModel{
		CreateTime: time.Now().Unix(),
		Username:   username,
		Password:   encodePWD,
	})
	if result.RowsAffected > 0 && result.Error == nil {
		return true
	}
	return false
}

type tokenEntity struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

func LoginVerify(username string, password string) (is bool, u model.UserModel, t tokenEntity) {
	r := dao.UserDao{}.D().Where("username = ?", username).First(&u)
	if r.RowsAffected == 0 || r.Error != nil {
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) //验证（对比）
	if err == nil {
		//生成jwt
		token, exp := JwtUtil.CreateToken(strconv.FormatInt(u.ID, 10), conf.TokenExp, conf.JwtSecret)
		t.Token = token
		t.Exp = exp
		//存入redis
		db.Redis().Set(AuthToken+strconv.FormatInt(u.ID, 10), token, time.Second*time.Duration(conf.TokenExp))
		return true, u, t
	}
	return
}

func TokenVerify(token string) (uid int64, errMsg string) {
	uidStr, err := JwtUtil.ParseToken(token, conf.JwtSecret)
	uid, _ = strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		errMsg = err.Error()
		return
	}
	//单点登录
	redisToken, redisErr := db.Redis().Get(AuthToken + uidStr).Result()
	if redisErr != nil {
		fmt.Println("redisErr", redisErr.Error())
	}
	if redisToken == "" {
		uid = 0
		errMsg = ApiUtil.UnauthorizedMsg
		return
	}
	if redisToken != token {
		uid = -1
		errMsg = ApiUtil.UnauthorizedSMsg
		return
	}
	return
}

func TokenRefresh(uid int64) (t tokenEntity) {
	tokenNew, exp := JwtUtil.CreateToken(strconv.FormatInt(uid, 10), conf.TokenExp, conf.JwtSecret)
	t.Token = tokenNew
	t.Exp = exp
	//存入redis
	db.Redis().Set(AuthToken+strconv.FormatInt(uid, 10), tokenNew, time.Second*time.Duration(conf.TokenExp))
	return
}

func Logout(uid int64) bool {
	err := db.Redis().Del(AuthToken + strconv.FormatInt(uid, 10)).Err()
	if err != nil {
		return false
	}
	return true
}
