package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SubmitDriverCertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitDriverCertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitDriverCertLogic {
	return &SubmitDriverCertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubmitDriverCertLogic) SubmitDriverCert(in *driver.SubmitDriverCertReq) (*driver.SubmitDriverCertResp, error) {

	tx := global.DB.Begin()

	// 更新司机主表状态（DriverUser 无 RealName 字段，真实姓名存于 DriverAuth）
	err := tx.Model(&user3.DriverUser{}).Where("id = ?", in.DriverId).
		Update("status", int8(1)).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("修改司机信息失败")
	}

	LicenseExpire, _ := time.Parse("2006-01-02", in.LicenseExpire)

	profile := user3.DriverProfile{
		ID:            in.DriverId,
		LicenseNo:     in.LicenseNo,
		LicenseExpire: &LicenseExpire,
		DriveYears:    int(in.DriveYears),
	}

	err = profile.UpdateData(tx)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("司机资料修改失败")
	}

	auth := user3.DriverAuth{
		Model:       gorm.Model{ID: uint(in.DriverId)},
		DriverID:    in.DriverId,
		RealName:    in.RealName,
		AuditStatus: 1,
	}

	err = auth.UpdateData(tx)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("修改认证表信息失败")
	}

	tx.Commit()
	return &driver.SubmitDriverCertResp{
		Success: true,
	}, nil
}
