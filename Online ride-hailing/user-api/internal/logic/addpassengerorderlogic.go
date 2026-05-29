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

type AddPassengerOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPassengerOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPassengerOrderLogic {
	return &AddPassengerOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPassengerOrderLogic) AddPassengerOrder(req *types.AddPassengerOrderReq) (resp *types.Request, err error) {
	// 参数验证
	if req.UserId <= 0 {
		return &types.Request{
			Code:    400,
			Message: "用户ID无效",
			Data:    nil,
		}, nil
	}
	if req.PassengerName == "" {
		return &types.Request{
			Code:    400,
			Message: "乘客姓名不能为空",
			Data:    nil,
		}, nil
	}
	if req.StartAddress == "" || req.EndAddress == "" {
		return &types.Request{
			Code:    400,
			Message: "出发地和目的地不能为空",
			Data:    nil,
		}, nil
	}

	data, err := l.svcCtx.UserRpc.AddPassengerOrder(l.ctx, &usermodel.AddPassengerOrderReq{
		UserId:         req.UserId,
		PassengerName:  req.PassengerName,
		PassengerPhone: req.PassengerPhone,
		StartAddress:   req.StartAddress,
		StartLng:       float32(req.StartLng),
		StartLat:       float32(req.StartLat),
		EndAddress:     req.EndAddress,
		EndLng:         float32(req.EndLng),
		EndLat:         float32(req.EndLat),
		CarType:        req.CarType,
		CouponId:       req.CouponId,
		Remark:         req.Remark,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "创建订单失败",
			Data:    nil,
		}, nil
	}
	// 检查 RPC 返回的业务状态码
	if data.Code != 200 {
		return &types.Request{
			Code:    int32(data.Code),
			Message: data.Message,
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "创建订单成功",
		Data:    data,
	}, nil
}
