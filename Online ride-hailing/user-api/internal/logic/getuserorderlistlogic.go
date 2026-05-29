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

type GetUserOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOrderListLogic {
	return &GetUserOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserOrderListLogic) GetUserOrderList(req *types.UserOrderListReq) (resp *types.Request, err error) {
	// todo: add your logic here and delete this line
	data, err := l.svcCtx.UserRpc.GetUserOrderList(l.ctx, &usermodel.UserOrderListReq{
		UserId:   req.UserId,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "用户订单查询失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "用户订单查询成功",
		Data:    data,
	}, nil
}
