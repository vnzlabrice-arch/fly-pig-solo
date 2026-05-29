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

type GetUserCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCouponLogic {
	return &GetUserCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserCouponLogic) GetUserCoupon(req *types.GetUserCouponReq) (resp *types.Request, err error) {
	// todo: add your logic here and delete this line
	data, err := l.svcCtx.UserRpc.GetUserCouponList(l.ctx, &usermodel.UserCouponListReq{
		UserId:   req.UserId,
		Status:   int32(req.Status),
		Page:     int32(req.Page),
		PageSize: req.PageSize,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "查询用户优惠券失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "查询用户优惠券成功",
		Data:    data,
	}, nil
}
