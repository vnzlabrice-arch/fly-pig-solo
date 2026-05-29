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

type OutLoginReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOutLoginReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OutLoginReqLogic {
	return &OutLoginReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OutLoginReqLogic) OutLoginReq(req *types.OutLoginReq) (resp *types.Response, err error) {
	res, err := l.svcCtx.DriverSrv.OutLogin(l.ctx, &driver.OutLoginReq{})
	if err != nil {
		return &types.Response{
			Code: 0,
			Msg:  "退出登录失败",
			Data: nil,
		}, nil
	}

	return &types.Response{
		Code: 200,
		Msg:  "退出登录成功",
		Data: res,
	}, nil

}
