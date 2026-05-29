// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"admin-api/internal/config"
	"admin-srv/adminclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	AdminSrv adminclient.Admin
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		AdminSrv: adminclient.NewAdmin(zrpc.MustNewClient(c.AdminSrv)),
	}
}
