package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	APP_KEY              = "www.topgoer.com"
	TokenExpireDuration = time.Hour * 24
)

// TokenHandler 生成JWT Token
func TokenHandler(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(TokenExpireDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(APP_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
