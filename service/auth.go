// Package service 承载业务逻辑。
package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"mini-im/conf"
	"mini-im/dao"
	"mini-im/db"
	"mini-im/model"
	"mini-im/pkg/jwtutil"
	"mini-im/pkg/response"
)

// 注册限制
const (
	// MinPasswordLen 密码最小长度
	MinPasswordLen = 6
)

// UsernamePattern 用户名规则：6~20 位字母、数字、下划线或减号
var UsernamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{6,20}$`)

var (
	// ErrCredential 用户名或密码错误
	ErrCredential = errors.New("用户名或密码错误")
	// ErrUsernameExists 用户名已存在
	ErrUsernameExists = errors.New("用户名已存在")
	// ErrUsernameInvalid 用户名不符合命名规则
	ErrUsernameInvalid = errors.New("用户名需为 6-20 位字母、数字、下划线或减号组合")
	// ErrPasswordTooShort 密码过短
	ErrPasswordTooShort = errors.New("密码至少 6 位")
	// ErrTokenInvalid token 无效或已过期
	ErrTokenInvalid = errors.New(response.UnauthorizedMsg)
	// ErrTokenReplaced 账号在别处登录（单点登录被挤下线）
	ErrTokenReplaced = errors.New(response.OtherLoginMsg)
)

// Token 登录凭证
type Token struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

// tokenTTL token 有效期
func tokenTTL() time.Duration { return time.Duration(conf.C.JWT.Expire) * time.Second }

// Register 注册新账号
//
// 账号名与密码的规则在这里校验：业务层是唯一入口，HTTP 之外（脚本 / 任务）调用也绕不过。
func Register(username, password, nickname string) (*model.User, error) {
	if !UsernamePattern.MatchString(username) {
		return nil, ErrUsernameInvalid
	}
	if len(password) < MinPasswordLen {
		return nil, ErrPasswordTooShort
	}

	userDao := dao.NewUserDao()

	if _, err := userDao.FindByUsername(username); err == nil {
		return nil, ErrUsernameExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	if nickname == "" {
		nickname = username
	}
	user := &model.User{
		Username:   username,
		Nickname:   nickname,
		Avatar:     db.AvatarURL(username, nickname),
		Password:   string(hash),
		CreateTime: now,
		UpdateTime: now,
	}
	if err = userDao.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 校验账号密码并签发 token
func Login(username, password string) (*model.User, Token, error) {
	user, err := dao.NewUserDao().FindByUsername(username)
	if err != nil {
		// 不区分「用户不存在」与「密码错误」，避免账号枚举
		return nil, Token{}, ErrCredential
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, Token{}, ErrCredential
	}

	token, err := issueToken(user.ID)
	if err != nil {
		return nil, Token{}, err
	}
	return user, token, nil
}

// issueToken 签发 token 并写入 Redis（覆盖旧值即实现单点登录）
func issueToken(uid int64) (Token, error) {
	token, exp, err := jwtutil.Create(uid, conf.C.JWT.Secret, tokenTTL())
	if err != nil {
		return Token{}, err
	}
	if err = db.Redis().Set(db.Ctx, conf.C.AuthTokenKey(uid), token, tokenTTL()).Err(); err != nil {
		return Token{}, err
	}
	return Token{Token: token, Exp: exp}, nil
}

// RefreshToken 续签 token
func RefreshToken(uid int64) (Token, error) { return issueToken(uid) }

// VerifyToken 校验 token 并返回用户 ID
func VerifyToken(token string) (int64, error) {
	uid, err := jwtutil.Parse(token, conf.C.JWT.Secret)
	if err != nil {
		return 0, ErrTokenInvalid
	}

	saved, err := db.Redis().Get(db.Ctx, conf.C.AuthTokenKey(uid)).Result()
	if errors.Is(err, redis.Nil) {
		return 0, ErrTokenInvalid
	}
	if err != nil {
		return 0, ErrTokenInvalid
	}
	if saved != token {
		return 0, ErrTokenReplaced
	}
	return uid, nil
}

// Logout 注销登录态
func Logout(uid int64) error {
	return db.Redis().Del(db.Ctx, conf.C.AuthTokenKey(uid)).Err()
}
