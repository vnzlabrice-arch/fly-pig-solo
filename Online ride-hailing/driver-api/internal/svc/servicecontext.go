// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"driver-api/internal/config"
	"driver-srv/driverclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	DriverSrv driverclient.Driver
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		DriverSrv: driverclient.NewDriver(zrpc.MustNewClient(c.DriverSrv)),
	}
}
