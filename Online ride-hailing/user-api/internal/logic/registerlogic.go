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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.Request, err error) {

	data, err := l.svcCtx.UserRpc.Register(l.ctx, &usermodel.RegisterReq{

		Phone:    req.Phone,
		Code:     req.Code,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "注册失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    200,
		Message: "注册成功",
		Data:    data,
	}, nil
}
