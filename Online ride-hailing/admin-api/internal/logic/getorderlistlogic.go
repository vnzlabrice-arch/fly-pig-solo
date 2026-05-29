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

type GetOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderListLogic) GetOrderList(req *types.GetOrderListRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetOrderList(l.ctx, &admin.GetOrderListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Status:   req.Status,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "查询订单列表失败: " + err.Error(),
		}, nil
	}

	var list []interface{}
	for _, item := range rpcResp.List {
		list = append(list, map[string]interface{}{
			"order_id":        item.OrderId,
			"passenger_name":  item.PassengerName,
			"passenger_phone": item.PassengerPhone,
			"driver_name":     item.DriverName,
			"driver_phone":    item.DriverPhone,
			"car_type":        item.CarType,
			"start_address":   item.StartAddress,
			"end_address":     item.EndAddress,
			"status":          item.Status,
			"estimated_price": item.EstimatedPrice,
			"final_price":     item.FinalPrice,
			"created_at":      item.CreatedAt,
		})
	}

	return &types.Request{
		Code:    200,
		Message: "查询成功",
		Data: map[string]interface{}{
			"total": rpcResp.Total,
			"list":  list,
		},
	}, nil
}
