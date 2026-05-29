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

type CartAuthReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartAuthReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartAuthReqLogic {
	return &CartAuthReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartAuthReqLogic) CartAuthReq(req *types.CartAuthReq) (resp *types.Response, err error) {
	res, err := l.svcCtx.DriverSrv.CartAuth(l.ctx, &driver.CartAuthReq{
		CarPlate:       req.CarPlate,
		CarModel:       req.CarModel,
		DrivingLicense: req.DrivingLicense,
		Token:          req.Token,
	})
	if err != nil {
		return &types.Response{
			Code: 0,
			Msg:  "车辆失败",
			Data: nil,
		}, nil
	}

	return &types.Response{
		Code: 200,
		Msg:  "车辆成功",
		Data: res,
	}, nil
}
