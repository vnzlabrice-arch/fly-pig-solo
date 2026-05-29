package logic

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AcceptOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAcceptOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptOrderLogic {
	return &AcceptOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 抢单/接单
func (l *AcceptOrderLogic) AcceptOrder(in *driver.AcceptOrderReq) (*driver.AcceptOrderResp, error) {
	if in.DriverId <= 0 || in.OrderId == "" {
		return nil, errors.New("司机ID和订单号不能为空")
	}
	if in.OrderDistance <= 0 {
		return nil, errors.New("订单距离校验失败")
	}
	if in.DriverCurrentLocation == nil {
		return nil, errors.New("司机实时位置不能为空")
	}

	h := sha1.New()
	h.Write([]byte(in.OrderId))
	h.Write([]byte(strconv.FormatInt(in.DriverId, 10)))
	antiFraudToken := hex.EncodeToString(h.Sum(nil))[:16]
	if in.AntiFraudToken != "" && in.AntiFraudToken != antiFraudToken {
		return nil, errors.New("防重复提交校验失败")
	}

	lockTime := time.Now().Add(90 * time.Second)

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var driverOnline user3.DriverOnlineStatus
		err := tx.Where("driver_id = ?", in.DriverId).First(&driverOnline).Error
		if err != nil {
			return errors.New("司机未上线")
		}
		if driverOnline.OnlineStatus != 1 {
			return errors.New("司机当前未上线")
		}
		if driverOnline.AcceptOrder != 1 {
			return errors.New("司机当前未开启接单")
		}

		var driverOrder user3.DriverOrderRecommend
		err = tx.Where("driver_id = ? AND order_id = ?", in.DriverId, in.OrderId).First(&driverOrder).Error
		if err != nil {
			return errors.New("订单不存在")
		}
		if driverOrder.Status != 1 {
			return errors.New("订单当前不可接单")
		}

		var passengerOrder user3.PassengerOrder
		err = tx.Where("order_id = ?", in.OrderId).First(&passengerOrder).Error
		if err != nil {
			return errors.New("订单不存在")
		}
		if passengerOrder.Status != 1 {
			return errors.New("订单当前不可接单")
		}

		passengerOrder.DriverID = in.DriverId
		passengerOrder.Status = 2
		passengerOrder.DriverLng = float64(in.DriverCurrentLocation.Lng)
		passengerOrder.DriverLat = float64(in.DriverCurrentLocation.Lat)
		passengerOrder.EstimatedArrivalTime = new(time.Now().Add(time.Duration(int(in.OrderDistance/500)+1) * time.Minute))

		err = tx.Model(&passengerOrder).Updates(map[string]interface{}{
			"driver_id":              passengerOrder.DriverID,
			"status":                 passengerOrder.Status,
			"driver_lng":             passengerOrder.DriverLng,
			"driver_lat":             passengerOrder.DriverLat,
			"estimated_arrival_time": passengerOrder.EstimatedArrivalTime,
		}).Error
		if err != nil {
			return err
		}

		driverOrder.Distance = int(in.OrderDistance)
		driverOrder.Status = 2
		driverOrder.OrderLockTime = &lockTime
		driverOrder.AntiFraudToken = antiFraudToken
		err = tx.Model(&driverOrder).Updates(map[string]interface{}{
			"distance":         driverOrder.Distance,
			"status":           driverOrder.Status,
			"order_lock_time":  driverOrder.OrderLockTime,
			"anti_fraud_token": driverOrder.AntiFraudToken,
		}).Error
		if err != nil {
			return err
		}

		driverOnline.Lng = float64(in.DriverCurrentLocation.Lng)
		driverOnline.Lat = float64(in.DriverCurrentLocation.Lat)
		err = tx.Model(&driverOnline).Updates(map[string]interface{}{
			"lng": driverOnline.Lng,
			"lat": driverOnline.Lat,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &driver.AcceptOrderResp{
		Success:       true,
		OrderLockTime: lockTime.Format("2006-01-02T15:04:05Z"),
	}, nil
}
