package logic

import (
	"context"
	"errors"

	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消订单
func (l *CancelOrderLogic) CancelOrder(in *driver.CancelOrderReq) (*driver.CancelOrderResp, error) {
	if in.DriverId <= 0 || in.OrderId == "" {
		return nil, errors.New("司机ID和订单号不能为空")
	}
	if in.Reason == "" {
		return nil, errors.New("取消原因不能为空")
	}

	cancelFee := float32(0)
	cancelResponsibility := "司机"
	appealSupport := true

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
		if driverOrder.Status == 5 || driverOrder.Status == 6 {
			return errors.New("订单已结束，不能重复取消")
		}

		switch driverOrder.Status {
		case 2:
			cancelFee = 8
		case 3:
			cancelFee = 12
		case 4:
			cancelFee = 20
			cancelResponsibility = "平台"
			appealSupport = false
		}

		passengerOrder.Status = 6
		passengerOrder.CancelReason = in.Reason
		passengerOrder.CancelFee = float64(cancelFee)
		passengerOrder.CancelResponsibility = cancelResponsibility
		passengerOrder.EvidencePhoto = in.EvidencePhoto
		passengerOrder.AppealSupport = appealSupport
		err = tx.Model(&passengerOrder).Updates(map[string]interface{}{
			"status":                passengerOrder.Status,
			"cancel_reason":         passengerOrder.CancelReason,
			"cancel_fee":            passengerOrder.CancelFee,
			"cancel_responsibility": passengerOrder.CancelResponsibility,
			"evidence_photo":        passengerOrder.EvidencePhoto,
			"appeal_support":        passengerOrder.AppealSupport,
		}).Error
		if err != nil {
			return err
		}

		driverOrder.Status = 6
		err = tx.Model(&driverOrder).Update("status", driverOrder.Status).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &driver.CancelOrderResp{
		Success:              true,
		CancelFee:            cancelFee,
		CancelResponsibility: cancelResponsibility,
		AppealSupport:        appealSupport,
	}, nil
}
