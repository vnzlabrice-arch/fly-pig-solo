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

type StartTripLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartTripLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartTripLogic {
	return &StartTripLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 开始行程
func (l *StartTripLogic) StartTrip(in *driver.StartTripReq) (*driver.StartTripResp, error) {
	if in.DriverId <= 0 || in.OrderId == "" {
		return nil, errors.New("司机ID和订单号不能为空")
	}
	if in.StartOdometer <= 0 {
		return nil, errors.New("起始里程必须大于0")
	}
	if in.PassengerConfirmCode == "" {
		return nil, errors.New("乘客确认码不能为空")
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
		if passengerOrder.Status != 3 || driverOrder.Status != 3 {
			return errors.New("当前订单状态不允许开始行程")
		}

		expectedCode := passengerOrder.PassengerConfirmCode
		if expectedCode == "" {
			if len(passengerOrder.OrderID) >= 4 {
				expectedCode = passengerOrder.OrderID[len(passengerOrder.OrderID)-4:]
			} else {
				expectedCode = "0000"
			}
		}
		if in.PassengerConfirmCode != expectedCode {
			return errors.New("乘客确认码校验失败")
		}

		passengerOrder.Status = 4
		passengerOrder.StartTime = &now
		passengerOrder.PickupPhoto = in.PickupPhoto
		passengerOrder.StartOdometer = float64(in.StartOdometer)
		passengerOrder.PassengerConfirmCode = expectedCode
		err = tx.Model(&passengerOrder).Updates(map[string]interface{}{
			"status":                 passengerOrder.Status,
			"start_time":             passengerOrder.StartTime,
			"pickup_photo":           passengerOrder.PickupPhoto,
			"start_odometer":         passengerOrder.StartOdometer,
			"passenger_confirm_code": passengerOrder.PassengerConfirmCode,
		}).Error
		if err != nil {
			return err
		}

		driverOrder.Status = 4
		err = tx.Model(&driverOrder).Update("status", driverOrder.Status).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &driver.StartTripResp{Success: true}, nil
}
