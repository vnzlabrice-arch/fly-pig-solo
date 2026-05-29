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
)

type DriverOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDriverOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverOfflineLogic {
	return &DriverOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 司机下线
func (l *DriverOfflineLogic) DriverOffline(in *driver.DriverOfflineReq) (*driver.DriverOfflineResp, error) {
	if in.DriverId <= 0 {
		return nil, errors.New("司机ID不能为空")
	}

	var driverOnline user3.DriverOnlineStatus
	err := global.DB.Where("driver_id = ?", in.DriverId).First(&driverOnline).Error
	if err != nil {
		return nil, errors.New("司机状态不存在")
	}

	driverOnline.OnlineStatus = 0
	driverOnline.AcceptOrder = 0
	if in.DriverLocation != nil {
		driverOnline.Lng = float64(in.DriverLocation.Lng)
		driverOnline.Lat = float64(in.DriverLocation.Lat)
	}
	if in.CarStatus != 0 {
		driverOnline.CarStatus = int8(in.CarStatus)
	}
	if time.Since(driverOnline.UpdatedAt) >= 4*time.Hour {
		driverOnline.RestReminder = "您已持续在线较久，请注意休息"
	}

	err = global.DB.Model(&driverOnline).Updates(map[string]interface{}{
		"online_status": driverOnline.OnlineStatus,
		"accept_order":  driverOnline.AcceptOrder,
		"lng":           driverOnline.Lng,
		"lat":           driverOnline.Lat,
		"car_status":    driverOnline.CarStatus,
		"rest_reminder": driverOnline.RestReminder,
	}).Error
	if err != nil {
		return nil, err
	}

	return &driver.DriverOfflineResp{
		Success:      true,
		RestReminder: driverOnline.RestReminder,
	}, nil
}
