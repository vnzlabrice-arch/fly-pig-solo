// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"user-api/internal/svc"
	"user-api/internal/types"
	"user-srv/pb/usermodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCodeLogic {
	return &SendCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCodeLogic) SendCode(req *types.SendCodeReq) (resp *types.Request, err error) {
	data, err := l.svcCtx.UserRpc.SendCode(l.ctx, &usermodel.SendCodeReq{

		Phone: req.Phone,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "验证码发送失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "验证码发送成功",
		Data:    data,
	}, nil
}
