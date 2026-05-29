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

type DriverOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDriverOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverOnlineLogic {
	return &DriverOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 司机上线
func (l *DriverOnlineLogic) DriverOnline(in *driver.DriverOnlineReq) (*driver.DriverOnlineResp, error) {
	if in.DriverId <= 0 {
		return nil, errors.New("司机ID不能为空")
	}
	if in.DriverLocation == nil {
		return nil, errors.New("司机位置不能为空")
	}

	var driverOnline user3.DriverOnlineStatus
	err := global.DB.Where("driver_id = ?", in.DriverId).First(&driverOnline).Error
	if err != nil {
		driverOnline = user3.DriverOnlineStatus{DriverID: in.DriverId}
	}

	driverOnline.OnlineStatus = 1
	driverOnline.Lng = float64(in.DriverLocation.Lng)
	driverOnline.Lat = float64(in.DriverLocation.Lat)
	if in.CarStatus != 0 {
		driverOnline.CarStatus = int8(in.CarStatus)
	} else if driverOnline.CarStatus == 0 {
		driverOnline.CarStatus = 1
	}
	if in.OnlineType != 0 {
		driverOnline.OnlineType = int8(in.OnlineType)
	} else if driverOnline.OnlineType == 0 {
		driverOnline.OnlineType = 1
	}
	if time.Since(driverOnline.UpdatedAt) >= 4*time.Hour {
		driverOnline.RestReminder = "您已持续在线较久，请注意休息"
	}

	if driverOnline.ID == 0 {
		err = global.DB.Create(&driverOnline).Error
	} else {
		err = global.DB.Model(&driverOnline).Updates(map[string]interface{}{
			"online_status": driverOnline.OnlineStatus,
			"lng":           driverOnline.Lng,
			"lat":           driverOnline.Lat,
			"car_status":    driverOnline.CarStatus,
			"online_type":   driverOnline.OnlineType,
			"rest_reminder": driverOnline.RestReminder,
		}).Error
	}
	if err != nil {
		return nil, err
	}

	return &driver.DriverOnlineResp{
		Success:      true,
		RestReminder: driverOnline.RestReminder,
	}, nil
}
