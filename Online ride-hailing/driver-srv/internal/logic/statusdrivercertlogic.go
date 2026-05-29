package logic

import (
	"context"
	"driver-srv/global"
	user3 "driver-srv/model/driver"
	"errors"

	"driver-srv/internal/svc"
	"driver-srv/pb/driver"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatusDriverCertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatusDriverCertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusDriverCertLogic {
	return &StatusDriverCertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StatusDriverCertLogic) StatusDriverCert(in *driver.StatusDriverCertReq) (*driver.StatusDriverCertResp, error) {
	//1. 从token拿司机ID
	//driverId, ok := l.ctx.Value("driver_id").(int64)
	//if !ok || driverId <= 0 {
	//	return nil, errors.New("未登录")
	//}

	var driverUser user3.DriverUser
	err := global.DB.Where("id = ?", in.DriverId).Find(&driverUser).Error
	if err != nil {
		return nil, errors.New("数据查询失败1")
	}

	var driverProfile user3.DriverProfile
	err = global.DB.Where("driver_id = ?", in.DriverId).Find(&driverProfile).Error
	if err != nil {
		return nil, errors.New("数据查询失败2")
	}

	var driverAuth user3.DriverAuth
	err = global.DB.Where("driver_id = ?", in.DriverId).Order("id desc").Find(&driverAuth).Error
	if err != nil {
		return nil, errors.New("数据查询失败3")
	}

	return &driver.StatusDriverCertResp{
		AuthType:    0, // AuthType 字段当前未在 DriverAuth 模型中启用
		AuditStatus: int64(driverAuth.AuditStatus),
		Reason:      driverAuth.Reason,
		LicenseNo:   driverProfile.LicenseNo,
		RealName:    driverAuth.RealName,
	}, nil
}
