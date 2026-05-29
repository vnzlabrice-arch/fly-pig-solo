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

type DriverAuthReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDriverAuthReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverAuthReqLogic {
	return &DriverAuthReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriverAuthReqLogic) DriverAuthReq(req *types.DriverAuthReq) (resp *types.Response, err error) {
	res, err := l.svcCtx.DriverSrv.DriverAuth(l.ctx, &driver.DriverAuthReq{
		RealName:    req.RealName,
		IDCard:      req.IDCard,
		IDCardFront: req.IDCardFront,
		IDCardBack:  req.IDCardBack,
		LicenseImg:  req.LicenseImg,
		Token:       req.Token,
	})
	if err != nil {
		return &types.Response{
			Code: 0,
			Msg:  "实名认证失败",
			Data: nil,
		}, nil
	}

	return &types.Response{
		Code: 200,
		Msg:  "实名认证成功",
		Data: res,
	}, nil

}
