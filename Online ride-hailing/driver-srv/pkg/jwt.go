package pkg

import (
	"context"
	"driver-srv/global"
	"errors"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	APP_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
)

const TokenExpireDuration = 2 * time.Hour

func TokenHandler(driverId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"driverId": driverId,
		"exp":      time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":      time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(APP_KEY))
	return tokenString, err
}

func GetToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(APP_KEY), nil
	})
	if err != nil {
		return nil, errors.New("token无效 或 解析失败")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims, err
	}
	return nil, nil
}

// 退出登录：把token加入黑名单

const blacklistPrefix = "jwt:blacklist:"

// BlacklistToken 把token加入黑名单
func BlacklistToken(tokenStr string, secret string) error {
	// 1. 解析token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return errors.New("token无效")
	}

	// 2. 获取过期时间
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("token解析失败")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("token无过期时间")
	}
	expTime := time.Unix(int64(exp), 0)
	remaining := time.Until(expTime)
	if remaining <= 0 {
		return errors.New("token已过期")
	}

	// 3. 存入Redis黑名单
	return global.RDB.Set(context.Background(), blacklistPrefix+tokenStr, "1", remaining).Err()
}

// IsBlacklisted 检查token是否被拉黑
func IsBlacklisted(tokenStr string) bool {
	exists, err := global.RDB.Exists(context.Background(), blacklistPrefix+tokenStr).Result()
	return err == nil && exists == 1
}
