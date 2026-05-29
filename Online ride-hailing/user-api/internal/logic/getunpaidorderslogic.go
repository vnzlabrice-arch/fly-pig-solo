package logic

import (
	"context"

	"user-api/internal/svc"
	"user-api/internal/types"
	"user-srv/pb/usermodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnpaidOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUnpaidOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnpaidOrdersLogic {
	return &GetUnpaidOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnpaidOrdersLogic) GetUnpaidOrders(req *types.GetUnpaidOrdersReq) (resp *types.Request, err error) {
	// 调用RPC服务
	rpcResp, err := l.svcCtx.UserRpc.GetUnpaidOrders(l.ctx, &usermodel.GetUnpaidOrdersReq{
		HoursBefore: int32(req.HoursBefore),
	})

	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	if rpcResp.Code != 200 {
		return &types.Request{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		}, nil
	}

	// 转换订单列表格式
	var orderList []map[string]interface{}
	for _, order := range rpcResp.Orders {
		orderList = append(orderList, map[string]interface{}{
			"order_id":        order.OrderId,
			"passenger_id":    order.PassengerId,
			"passenger_name":  order.PassengerName,
			"passenger_phone": order.PassengerPhone,
			"start_address":   order.StartAddress,
			"end_address":     order.EndAddress,
			"actual_price":    order.ActualPrice,
			"finish_time":     order.FinishTime,
		})
	}

	return &types.Request{
		Code:    0,
		Message: "success",
		Data: map[string]interface{}{
			"orders": orderList,
		},
	}, nil
}
