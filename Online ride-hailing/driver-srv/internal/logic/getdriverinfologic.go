package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetDriverInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDriverInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDriverInfoLogic {
	return &GetDriverInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 司机信息模块
//func (l *GetDriverInfoLogic) GetDriverInfo(in *driver.GetDriverInfoReq) (*driver.GetDriverInfoResp, error) {
//
//	//从token里面获取的id
//	//token := l.svcCtx.Request.Header.Get("token")
//	//if token == "" {
//	//	return nil, errors.New("token不能为空")
//	//}
//	//
//	//// 2. 解析 token 获取 driverId
//	//driverId, err := jwt.GetDriverId(token)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	var driverModel user3.DriverUser
//	err := global.DB.Where("id = ?", in.DriverID).Find(&driverModel).Error
//	if err != nil {
//		return nil, errors.New("查询数据失败")
//	}
//
//	var driverProfile user3.DriverProfile
//	err = global.DB.Where("driver_id = ?", in.DriverID).Find(&driverProfile).Error
//	if err != nil {
//		return nil, errors.New("资料查询失败")
//	}
//
//	//var driverAuth user3.DriverAuth
//	//err = global.DB.Where("driver_id = ? AND auth_type = 1", in.DriverID, driverAuth.AuthType == 1).
//	//	Find(&driverAuth).Error
//	//if err != nil {
//	//	return nil, errors.New("实名查询失败")
//	//}
//
//	var realAuth int32 = 0
//	if err = global.DB.Where("driver_id = ? AND auth_type = ? AND audit_status = ?",
//		in.DriverID, 1, 1).Find(&user3.DriverAuth{}).Error; err == nil {
//		realAuth = 1
//	}
//
//	// 3. 查驾照认证状态（auth_type=2，audit_status=1 为通过）
//	var licenseAuth int32 = 0
//	if err = global.DB.Where("driver_id = ? AND auth_type = ? AND audit_status = ?",
//		in.DriverID, 2, 1).Find(&user3.DriverAuth{}).Error; err == nil {
//		licenseAuth = 1
//	}
//
//	return &driver.GetDriverInfoResp{
//		DriverID:      driverModel.ID,
//		Phone:         driverModel.Phone,
//		RealName:      driverModel.RealName,
//		AuthType:      int64(realAuth),
//		AuditStatus:   int64(licenseAuth),
//		LicenseNo:     driverProfile.LicenseNo,
//		LicenseExpire: driverProfile.LicenseExpire.Format("2006-04-02"),
//		DriveYears:    int64(driverProfile.DriveYears),
//	}, nil
//}

func (l *GetDriverInfoLogic) GetDriverInfo(in *driver.GetDriverInfoReq) (*driver.GetDriverInfoResp, error) {
	// 1. 查询司机主信息（处理ErrRecordNotFound）
	var driverModel user3.DriverUser
	err := global.DB.Where("id = ?", in.DriverID).First(&driverModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("司机不存在")
		}
		return nil, errors.New("查询数据失败")
	}

	// 2. 查询司机资料（初始化默认值，处理ErrRecordNotFound）
	var driverProfile user3.DriverProfile
	driverProfile = user3.DriverProfile{
		LicenseNo:     "",
		LicenseExpire: nil,
		DriveYears:    0,
	}
	err = global.DB.Where("driver_id = ?", in.DriverID).First(&driverProfile).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("资料查询失败")
	}

	// 3. 查实名认证状态
	var realAuth int32 = 0
	if err := global.DB.Where("driver_id = ? AND auth_type = ? AND audit_status = ?",
		in.DriverID, 1, 1).First(&user3.DriverAuth{}).Error; err == nil {
		realAuth = 1
	}

	// 4. 查驾照认证状态
	var licenseAuth int32 = 0
	if err := global.DB.Where("driver_id = ? AND auth_type = ? AND audit_status = ?",
		in.DriverID, 2, 1).First(&user3.DriverAuth{}).Error; err == nil {
		licenseAuth = 1
	}

	// 5. 关键修复：处理LicenseExpire空指针
	licenseExpireStr := ""
	if driverProfile.LicenseExpire != nil {
		licenseExpireStr = driverProfile.LicenseExpire.Format("2006-01-02")
	}

	// 6. 从 DriverAuth 获取真实姓名
	realName := ""
	var authRecord user3.DriverAuth
	if err := global.DB.Where("driver_id = ? AND audit_status = ?",
		in.DriverID, 1).First(&authRecord).Error; err == nil {
		realName = authRecord.RealName
	}

	// 7. 返回数据
	return &driver.GetDriverInfoResp{
		DriverID:      int64(driverModel.ID),
		Phone:         driverModel.Phone,
		RealName:      realName,
		AuthType:      int64(realAuth),
		AuditStatus:   int64(licenseAuth),
		LicenseNo:     driverProfile.LicenseNo,
		LicenseExpire: licenseExpireStr,
		DriveYears:    int64(driverProfile.DriveYears),
	}, nil
}
