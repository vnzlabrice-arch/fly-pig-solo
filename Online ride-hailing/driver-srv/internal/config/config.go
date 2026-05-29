package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Jwt struct {
		Secret string `json:"secret"`
	} `json:"jwt"`
}
