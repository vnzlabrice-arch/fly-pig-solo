// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"driver-srv/pb/driver"

	"driver-api/internal/svc"
	"driver-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DriverLoginReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDriverLoginReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverLoginReqLogic {
	return &DriverLoginReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriverLoginReqLogic) DriverLoginReq(req *types.DriverLoginReq) (resp *types.Response, err error) {
	res, err := l.svcCtx.DriverSrv.DriverLogin(l.ctx, &driver.DriverLoginReq{
		Phone:    req.Phone,
		Code:     req.Code,
		DriverId: int64(req.DriverId),
	})
	if err != nil {
		return &types.Response{
			Code: 0,
			Msg:  "登录失败",
			Data: nil,
		}, nil
	}

	return &types.Response{
		Code: 200,
		Msg:  "登录成功",
		Data: res,
	}, nil
}
