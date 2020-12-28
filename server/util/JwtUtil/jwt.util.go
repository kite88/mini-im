package JwtUtil

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func CreateToken(uid string, timeoutSecond int64, secret string) (string, int64) {
	exp := time.Now().Add(time.Second * time.Duration(timeoutSecond)).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": exp,
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", 0
	}
	return token, exp
}

func ParseToken(token string, secret string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return claim.Claims.(jwt.MapClaims)["uid"].(string), nil
}
