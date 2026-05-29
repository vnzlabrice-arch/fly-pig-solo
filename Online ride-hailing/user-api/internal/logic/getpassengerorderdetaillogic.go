// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"user-srv/pb/usermodel"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPassengerOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPassengerOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPassengerOrderDetailLogic {
	return &GetPassengerOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPassengerOrderDetailLogic) GetPassengerOrderDetail(req *types.PassengerOrderDetailReq) (resp *types.Request, err error) {
	data, err := l.svcCtx.UserRpc.GetPassengerOrderDetail(l.ctx, &usermodel.PassengerOrderDetailReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "用户订单详情查询失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "用户订单详情查询成功",
		Data:    data,
	}, nil
}
