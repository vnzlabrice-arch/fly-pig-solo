package logic

import (
	"context"
	"errors"
	"time"

	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ArrivePassengerLocationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArrivePassengerLocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArrivePassengerLocationLogic {
	return &ArrivePassengerLocationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 到达乘客位置
func (l *ArrivePassengerLocationLogic) ArrivePassengerLocation(in *driver.ArrivePassengerLocationReq) (*driver.ArrivePassengerLocationResp, error) {
	if in.DriverId <= 0 || in.OrderId == "" {
		return nil, errors.New("司机ID和订单号不能为空")
	}
	if in.ArrivalAccuracy <= 0 {
		return nil, errors.New("GPS精度不能为空")
	}

	now := time.Now()
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var driverOrder user3.DriverOrderRecommend
		err := tx.Where("driver_id = ? AND order_id = ?", in.DriverId, in.OrderId).First(&driverOrder).Error
		if err != nil {
			return errors.New("订单不存在")
		}

		var passengerOrder user3.PassengerOrder
		err = tx.Where("order_id = ?", in.OrderId).First(&passengerOrder).Error
		if err != nil {
			return errors.New("订单不存在")
		}
		if passengerOrder.Status != 2 || driverOrder.Status != 2 {
			return errors.New("当前订单状态不允许执行到达操作")
		}

		passengerOrder.Status = 3
		passengerOrder.DriverLng = float64(in.DriverLng)
		passengerOrder.DriverLat = float64(in.DriverLat)
		passengerOrder.ArrivalAccuracy = float64(in.ArrivalAccuracy)
		passengerOrder.ArrivalPhoto = in.ArrivalPhoto
		passengerOrder.PickupTime = &now
		err = tx.Model(&passengerOrder).Updates(map[string]interface{}{
			"status":           passengerOrder.Status,
			"driver_lng":       passengerOrder.DriverLng,
			"driver_lat":       passengerOrder.DriverLat,
			"arrival_accuracy": passengerOrder.ArrivalAccuracy,
			"arrival_photo":    passengerOrder.ArrivalPhoto,
			"pickup_time":      passengerOrder.PickupTime,
		}).Error
		if err != nil {
			return err
		}

		driverOrder.Status = 3
		err = tx.Model(&driverOrder).Update("status", driverOrder.Status).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &driver.ArrivePassengerLocationResp{Success: true}, nil
}
