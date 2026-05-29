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

type AddCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCouponLogic {
	return &AddCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCouponLogic) AddCoupon(req *types.AddcouponReq) (resp *types.Request, err error) {
	// todo: add your logic here and delete this line
	data, err := l.svcCtx.UserRpc.AddCoupon(l.ctx, &usermodel.AddCouponReq{
		UserId:     req.UserId,
		TemplateId: req.TemplateId,
		Count:      int32(req.Count),
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "添加优惠券失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "添加优惠券成功",
		Data:    data,
	}, nil
}
