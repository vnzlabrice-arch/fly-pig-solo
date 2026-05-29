// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"user-api/internal/config"
	"user-api/internal/ws"
	"user-srv/pb/usermodel"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc usermodel.UserClient
	WSHub   *ws.Hub
}

func NewServiceContext(c config.Config, hub *ws.Hub) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: usermodel.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
		WSHub:   hub,
	}
}
