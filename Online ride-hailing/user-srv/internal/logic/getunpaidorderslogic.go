package logic

import (
	"context"
	"time"
	"user-srv/model"

	"user-srv/global"
	"user-srv/internal/svc"
	"user-srv/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnpaidOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUnpaidOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnpaidOrdersLogic {
	return &GetUnpaidOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUnpaidOrders 查询未支付订单
func (l *GetUnpaidOrdersLogic) GetUnpaidOrders(in *user.GetUnpaidOrdersReq) (*user.GetUnpaidOrdersResp, error) {
	// 设置默认查询24小时内的订单
	hoursBefore := in.HoursBefore
	if hoursBefore <= 0 {
		hoursBefore = 24
	}

	// 计算时间阈值
	timeThreshold := time.Now().Add(-time.Duration(hoursBefore) * time.Hour)

	// 查询条件：
	// 1. 订单状态为已完成（status=5）
	// 2. 支付状态为未支付（pay_status=0）
	// 3. 行程已结束（end_time不为空）
	// 4. 在指定时间范围内
	var orders []model.PassengerOrder
	err := global.DB.Where("status = ?", 5).
		Where("pay_status = ?", 0).
		Where("end_time IS NOT NULL").
		Where("end_time >= ?", timeThreshold).
		Order("end_time DESC").
		Find(&orders).Error

	if err != nil {
		l.Errorf("查询未支付订单失败: %v", err)
		return &user.GetUnpaidOrdersResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 转换为Proto响应格式
	var orderList []*user.PassengerOrderInfo
	for _, order := range orders {
		// 处理时间字段
		var bookTime int64 = order.BookTime.Unix()
		var acceptTime int64
		if order.PickupTime != nil {
			acceptTime = order.PickupTime.Unix()
		}
		var startTime int64
		if order.StartTime != nil {
			startTime = order.StartTime.Unix()
		}
		var finishTime int64
		if order.EndTime != nil {
			finishTime = order.EndTime.Unix()
		}

		orderList = append(orderList, &user.PassengerOrderInfo{
			OrderId:        order.OrderID,
			PassengerId:    order.PassengerID,
			PassengerName:  order.PassengerName,
			PassengerPhone: order.PassengerPhone,
			OrderType:      int32(order.OrderType),
			CarType:        order.CarType,
			Status:         int32(order.Status),
			StartAddress:   order.StartAddress,
			StartLng:       order.StartLng,
			StartLat:       order.StartLat,
			EndAddress:     order.EndAddress,
			EndLng:         order.EndLng,
			EndLat:         order.EndLat,
			CouponName:     order.CouponName,
			EstimatedPrice: order.EstimatedPrice,
			ActualPrice:    order.FinalPrice,
			Remark:         order.PassRemark,
			BookTime:       bookTime,
			AcceptTime:     acceptTime,
			StartTime:      startTime,
			FinishTime:     finishTime,
			CreateTime:     order.CreatedAt.Unix(),
		})
	}

	return &user.GetUnpaidOrdersResp{
		Code:    200,
		Message: "查询成功",
		Orders:  orderList,
	}, nil
}
