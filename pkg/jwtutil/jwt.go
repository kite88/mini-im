// Package jwtutil 封装 JWT 的签发与解析。
package jwtutil

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken 非法或已过期的 token
var ErrInvalidToken = errors.New("invalid token")

// Claims 自定义载荷
type Claims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

// Create 签发 token，返回 token 字符串与过期时间（Unix 秒）
func Create(uid int64, secret string, ttl time.Duration) (string, int64, error) {
	now := time.Now()
	expireAt := now.Add(ttl)

	claims := Claims{
		UID: strconv.FormatInt(uid, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return token, expireAt.Unix(), nil
}

// Parse 解析 token 并返回用户 ID
func Parse(tokenStr, secret string) (int64, error) {
	if tokenStr == "" {
		return 0, ErrInvalidToken
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	uid, err := strconv.ParseInt(claims.UID, 10, 64)
	if err != nil || uid <= 0 {
		return 0, ErrInvalidToken
	}
	return uid, nil
}
