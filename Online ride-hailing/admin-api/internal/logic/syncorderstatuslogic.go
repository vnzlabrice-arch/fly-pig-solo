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

type SyncOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncOrderStatusLogic {
	return &SyncOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncOrderStatusLogic) SyncOrderStatus(req *types.SyncOrderStatusRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.SyncOrderStatus(l.ctx, &admin.SyncOrderStatusRequest{
		OrderId:      req.OrderId,
		Status:       req.Status,
		DriverId:     req.DriverId,
		FinalPrice:   req.FinalPrice,
		CancelReason: req.CancelReason,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "同步订单状态失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
		Data: map[string]interface{}{
			"order_id": rpcResp.OrderId,
			"status":   rpcResp.Status,
		},
	}, nil
}
