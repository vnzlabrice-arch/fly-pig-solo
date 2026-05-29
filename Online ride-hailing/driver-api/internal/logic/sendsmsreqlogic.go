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

type SendSmsReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendSmsReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsReqLogic {
	return &SendSmsReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendSmsReqLogic) SendSmsReq(req *types.SendSmsReq) (resp *types.Response, err error) {
	res, err := l.svcCtx.DriverSrv.SendSms(l.ctx, &driver.SendSmsReq{
		Phone: req.Phone})
	if err != nil {
		return &types.Response{
			Code: 0,
			Msg:  "短信发送失败",
			Data: nil,
		}, nil
	}

	return &types.Response{
		Code: 200,
		Msg:  "短信发送成功",
		Data: res.Success,
	}, nil
}
