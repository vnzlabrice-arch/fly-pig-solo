package logic

import (
	"context"
	"errors"

	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopAcceptOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStopAcceptOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopAcceptOrderLogic {
	return &StopAcceptOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 停止接单
func (l *StopAcceptOrderLogic) StopAcceptOrder(in *driver.StopAcceptOrderReq) (*driver.StopAcceptOrderResp, error) {
	if in.DriverId <= 0 {
		return nil, errors.New("司机ID不能为空")
	}

	var driverOnline user3.DriverOnlineStatus
	err := global.DB.Where("driver_id = ?", in.DriverId).First(&driverOnline).Error
	if err != nil {
		return nil, errors.New("司机状态不存在")
	}

	err = global.DB.Model(&driverOnline).Update("accept_order", 0).Error
	if err != nil {
		return nil, err
	}

	return &driver.StopAcceptOrderResp{Success: true}, nil
}
