// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"admin-api/internal/svc"
	"admin-api/internal/types"
	admin "admin-srv/adminclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDispatchOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchOrderLogic {
	return &DispatchOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DispatchOrderLogic) DispatchOrder(req *types.DispatchOrderRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.DispatchOrder(l.ctx, &admin.DispatchOrderRequest{
		PassengerName:  req.PassengerName,
		PassengerPhone: req.PassengerPhone,
		StartAddress:   req.StartAddress,
		StartLng:       req.StartLng,
		StartLat:       req.StartLat,
		EndAddress:     req.EndAddress,
		EndLng:         req.EndLng,
		EndLat:         req.EndLat,
		CarType:        req.CarType,
		DriverId:       req.DriverId,
		Remark:         req.Remark,
		EstimatedPrice: req.EstimatedPrice,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "派单失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: rpcResp.Message,
		Data: map[string]interface{}{
			"success":  rpcResp.Success,
			"order_id": rpcResp.OrderId,
		},
	}, nil
}
