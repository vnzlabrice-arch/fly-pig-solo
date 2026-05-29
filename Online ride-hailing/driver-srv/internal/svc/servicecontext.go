package svc

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/config"
	"errors"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     global.DB,  // 这里改成你实际的全局DB
		Redis:  global.RDB, // 这里改成你实际的全局Redis
	}
}

// 从RPC上下文获取 token
func (svc *ServiceContext) GetRpcToken(ctx context.Context) (string, error) {
	//// 关键修复：用正确的key "token"
	//token, ok := ctx.Value("token").(string)
	//if !ok || token == "" {
	//	return "", errors.New("token不存在")
	//}
	//return token, nil

	// 1. 从gRPC上下文中获取Metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("未携带token")
	}

	// 2. 获取token（注意：gRPC Metadata的key会被自动转为小写）
	tokens := md.Get("token")
	if len(tokens) == 0 || tokens[0] == "" {
		return "", errors.New("token不存在")
	}

	return tokens[0], nil
}
