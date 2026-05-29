package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"

	"github.com/zeromicro/go-zero/core/logx"
)

type DriverOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDriverOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverOrderListLogic {
	return &DriverOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取进行中订单
func (l *DriverOrderListLogic) DriverOrderList(in *driver.DriverOrderListReq) (*driver.DriverOrderListResp, error) {
	var driverOrder user3.DriverOrderRecommend
	err := global.DB.Where("driver_id = ? AND (status = 2 OR status = 3 OR status = 4)", in.DriverId).
		Order("created_at desc").First(&driverOrder).Error
	if err != nil {
		return &driver.DriverOrderListResp{}, nil
	}

	var passengerOrder user3.PassengerOrder
	err = global.DB.Where("order_id = ?", driverOrder.OrderID).First(&passengerOrder).Error
	if err != nil {
		return nil, err
	}

	if driverOrder.Status != passengerOrder.Status {
		_ = global.DB.Model(&driverOrder).Update("status", passengerOrder.Status).Error
	}

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

	estimatedArrivalTime := ""
	if passengerOrder.EstimatedArrivalTime != nil && !passengerOrder.EstimatedArrivalTime.IsZero() {
		estimatedArrivalTime = passengerOrder.EstimatedArrivalTime.Format("2006-01-02T15:04:05Z")
	}

	return &driver.DriverOrderListResp{
		OrderId:              passengerOrder.OrderID,
		OrderType:            int64(passengerOrder.OrderType),
		CarType:              passengerOrder.CarType,
		Status:               int64(passengerOrder.Status),
		PassengerName:        passengerOrder.PassengerName,
		PassengerPhone:       passengerOrder.PassengerPhone,
		PassengerAvatar:      "",
		StartAddress:         passengerOrder.StartAddress,
		StartLng:             float32(passengerOrder.StartLng),
		StartLat:             float32(passengerOrder.StartLat),
		EndAddress:           passengerOrder.EndAddress,
		EndLng:               float32(passengerOrder.EndLng),
		EndLat:               float32(passengerOrder.EndLat),
		PassRemark:           passengerOrder.PassRemark,
		BookTime:             bookTime,
		PickupTime:           pickupTime,
		StartTime:            startTime,
		DriverRemark:         passengerOrder.DriverRemark,
		PassengerRiskLevel:   int64(passengerOrder.PassengerRiskLevel),
		EstimatedArrivalTime: estimatedArrivalTime,
		SurgePrice:           float32(passengerOrder.SurgePrice),
	}, nil
}
