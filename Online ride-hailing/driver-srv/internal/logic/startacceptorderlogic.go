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

type StartAcceptOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartAcceptOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartAcceptOrderLogic {
	return &StartAcceptOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 开始接单
func (l *StartAcceptOrderLogic) StartAcceptOrder(in *driver.StartAcceptOrderReq) (*driver.StartAcceptOrderResp, error) {
	if in.DriverId <= 0 {
		return nil, errors.New("司机ID不能为空")
	}

	var driverOnline user3.DriverOnlineStatus
	err := global.DB.Where("driver_id = ?", in.DriverId).First(&driverOnline).Error
	if err != nil {
		return nil, errors.New("司机状态不存在")
	}
	if driverOnline.OnlineStatus != 1 {
		return nil, errors.New("司机未上线，不能开始接单")
	}

	err = global.DB.Model(&driverOnline).Update("accept_order", 1).Error
	if err != nil {
		return nil, err
	}

	return &driver.StartAcceptOrderResp{Success: true}, nil
}
