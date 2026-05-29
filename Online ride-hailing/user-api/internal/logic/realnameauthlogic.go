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

type RealNameAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRealNameAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RealNameAuthLogic {
	return &RealNameAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RealNameAuthLogic) RealNameAuth(req *types.RealNameAuthReq) (resp *types.Request, err error) {
	data, err := l.svcCtx.UserRpc.RealNameAuth(l.ctx, &usermodel.RealNameAuthReq{
		UserId:   req.UserId,
		RealName: req.RealName,
		IdCard:   req.IdCard,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "实名验证失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "实名验证成功",
		Data:    data,
	}, nil
}
