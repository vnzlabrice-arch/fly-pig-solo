package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"admin-srv/global"
	"admin-srv/internal/svc"
	"admin-srv/model/order"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SyncOrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncOrderStatusLogic {
	return &SyncOrderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncOrderStatusLogic) SyncOrderStatus(in *admin.SyncOrderStatusRequest) (*admin.SyncOrderStatusResponse, error) {
	orderID := strings.TrimSpace(in.OrderId)
	if orderID == "" {
		return &admin.SyncOrderStatusResponse{
			Code:    400,
			Message: "订单ID不能为空",
		}, nil
	}

	if global.DB == nil {
		return &admin.SyncOrderStatusResponse{
			Code:    500,
			Message: "数据库未初始化",
		}, nil
	}

	updates, err := buildOrderStatusUpdates(in.Status, in.DriverId, in.FinalPrice, in.CancelReason, time.Now())
	if err != nil {
		return &admin.SyncOrderStatusResponse{
			Code:    400,
			Message: "无效的订单状态",
			OrderId: orderID,
			Status:  in.Status,
		}, nil
	}

	var passengerOrder order.PassengerOrder
	err = global.DB.Select("order_id").
		Where("order_id = ?", orderID).
		First(&passengerOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &admin.SyncOrderStatusResponse{
			Code:    404,
			Message: "订单不存在",
			OrderId: orderID,
			Status:  in.Status,
		}, nil
	}
	if err != nil {
		l.Errorf("查询订单失败: %v", err)
		return &admin.SyncOrderStatusResponse{
			Code:    500,
			Message: "查询订单失败",
			OrderId: orderID,
			Status:  in.Status,
		}, nil
	}

	if err := global.DB.Model(&order.PassengerOrder{}).
		Where("order_id = ?", orderID).
		Updates(updates).Error; err != nil {
		l.Errorf("同步订单状态失败: %v", err)
		return &admin.SyncOrderStatusResponse{
			Code:    500,
			Message: "同步订单状态失败",
			OrderId: orderID,
			Status:  in.Status,
		}, nil
	}

	return &admin.SyncOrderStatusResponse{
		Code:    200,
		Message: "订单状态同步成功",
		OrderId: orderID,
		Status:  in.Status,
	}, nil
}
