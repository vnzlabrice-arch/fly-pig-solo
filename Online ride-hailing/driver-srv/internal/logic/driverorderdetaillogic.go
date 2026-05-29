package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type DriverOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDriverOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverOrderDetailLogic {
	return &DriverOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 订单详情
func (l *DriverOrderDetailLogic) DriverOrderDetail(in *driver.DriverOrderDetailReq) (*driver.DriverOrderDetailResp, error) {
	var driverOrder user3.DriverOrderRecommend
	err := global.DB.Where("driver_id = ? AND order_id = ?", in.DriverId, in.OrderId).First(&driverOrder).Error
	if err != nil {
		return nil, errors.New("订单不存在或无权限查看")
	}

	var passengerOrder user3.PassengerOrder
	err = global.DB.Where("order_id = ?", in.OrderId).First(&passengerOrder).Error
	if err != nil {
		return nil, errors.New("订单不存在")
	}

	var income user3.DriverIncome
	_ = global.DB.Where("driver_id = ? AND order_id = ?", in.DriverId, in.OrderId).First(&income).Error

	bookTime := ""
	if !passengerOrder.BookTime.IsZero() {
		bookTime = passengerOrder.BookTime.Format("2006-01-02T15:04:05Z")
	}

	pickupTime := ""
	if passengerOrder.PickupTime != nil && !passengerOrder.PickupTime.IsZero() {
		pickupTime = passengerOrder.PickupTime.Format("2006-01-02T15:04:05Z")
	}

	startTime := ""
	if passengerOrder.StartTime != nil && !passengerOrder.StartTime.IsZero() {
		startTime = passengerOrder.StartTime.Format("2006-01-02T15:04:05Z")
	}

	endTime := ""
	if passengerOrder.EndTime != nil && !passengerOrder.EndTime.IsZero() {
		endTime = passengerOrder.EndTime.Format("2006-01-02T15:04:05Z")
	}

	payTime := ""
	if passengerOrder.PayTime != nil && !passengerOrder.PayTime.IsZero() {
		payTime = passengerOrder.PayTime.Format("2006-01-02T15:04:05Z")
	}

	passengerTags := []string{}
	if passengerOrder.PassengerTags != "" {
		passengerTags = append(passengerTags, passengerOrder.PassengerTags)
	}

	return &driver.DriverOrderDetailResp{
		DriverOrderDetailItem: []*driver.DriverOrderDetailItem{
			{
				OrderId:            passengerOrder.OrderID,
				OrderType:          int64(passengerOrder.OrderType),
				CarType:            passengerOrder.CarType,
				Status:             int64(passengerOrder.Status),
				StartAddress:       passengerOrder.StartAddress,
				StartLng:           float32(passengerOrder.StartLng),
				StartLat:           float32(passengerOrder.StartLat),
				EndAddress:         passengerOrder.EndAddress,
				EndLng:             float32(passengerOrder.EndLng),
				EndLat:             float32(passengerOrder.EndLat),
				PassengerName:      passengerOrder.PassengerName,
				PassengerPhone:     passengerOrder.PassengerPhone,
				PassengerAvatar:    "",
				PassRemark:         passengerOrder.PassRemark,
				EstimatedPrice:     float32(passengerOrder.EstimatedPrice),
				FinalPrice:         float32(passengerOrder.FinalPrice),
				PlatformFee:        float32(income.PlatformFee),
				ActualIncome:       float32(income.ActualIncome),
				BookTime:           bookTime,
				PickupTime:         pickupTime,
				StartTime:          startTime,
				EndTime:            endTime,
				CancelReason:       passengerOrder.CancelReason,
				PayStatus:          int64(passengerOrder.PayStatus),
				PayType:            passengerOrder.PayType,
				PayTime:            payTime,
				CouponDeduction:    float32(passengerOrder.CouponDeduction),
				DriverIncomeDetail: &driver.DriverIncomeDetailData{TotalFee: float32(income.TotalFee), PlatformFee: float32(income.PlatformFee), ActualIncome: float32(income.ActualIncome)},
				ComplaintStatus:    int64(passengerOrder.ComplaintStatus),
				PassengerTags:      passengerTags,
			},
		},
	}, nil
}
