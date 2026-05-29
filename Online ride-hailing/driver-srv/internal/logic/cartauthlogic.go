package logic

import (
	"context"
	"driver-srv/global"
	user3 "driver-srv/model/driver"
	"driver-srv/pkg"
	"errors"
	"strconv"

	"driver-srv/internal/svc"
	"driver-srv/pb/driver"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CartAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartAuthLogic {
	return &CartAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CartAuthLogic) CartAuth(in *driver.CartAuthReq) (*driver.CartAuthResp, error) {
	// 1. 校验token
	token := in.Token
	if token == "" {
		return nil, errors.New("token不能为空")
	}

	// 校验token是否在黑名单里
	if pkg.IsBlacklisted(token) {
		return nil, errors.New("token已失效，请重新登录")
	}

	// 2. 解析JWT获取driverId
	claims, err := pkg.GetToken(token)
	if err != nil {
		return nil, errors.New("token解析失败")
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token解析失败")
	}

	driverIdStr, ok := mapClaims["driverId"].(string)
	if !ok {
		return nil, errors.New("token中缺少driverId")
	}

	driverId, err := strconv.Atoi(driverIdStr)
	if err != nil {
		return nil, errors.New("driverId格式错误")
	}

	// 3. 校验参数：车牌号、车型、行驶证不能为空
	if in.CarPlate == "" || in.CarModel == "" || in.DrivingLicense == "" {
		return nil, errors.New("车牌号、车型、行驶证不能为空")
	}

	// 4. 验证司机是否存在
	var driverUser user3.DriverUser
	err = global.DB.Where("id = ?", driverId).First(&driverUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("司机不存在")
		}
		l.Logger.Errorf("查询司机失败: %v", err)
		return nil, errors.New("系统异常，请稍后重试")
	}

	// 5. 检查是否已有车辆认证记录
	var existingCar user3.DriverCar
	hasExistingRecord := false
	err = global.DB.Where("driver_id = ?", driverId).First(&existingCar).Error
	if err == nil {
		hasExistingRecord = true
		switch existingCar.Status {
		case 2: // 审核中
			return nil, errors.New("您的车辆认证正在审核中，请耐心等待")
		case 3: // 已通过
			return &driver.CartAuthResp{Msg: "您的车辆已完成认证"}, nil
		case 4: // 已驳回，允许重新提交
			break
		default:
			break
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Errorf("查询车辆认证记录失败: %v", err)
		return nil, errors.New("系统异常，请稍后重试")
	}

	// 6. OCR行驶证识别（mock实现，传入用户输入的车牌用于比对）
	license, err := pkg.OcrLicense(in.DrivingLicense, in.CarPlate)
	if err != nil {
		l.Logger.Errorf("行驶证识别失败: %v", err)
		return nil, errors.New("行驶证识别失败，请重新上传")
	}

	if !license.IsValid {
		return nil, errors.New("行驶证识别不通过，请上传清晰的行驶证照片")
	}

	// 7. 保存或更新车辆认证记录
	if hasExistingRecord {
		// 之前被驳回，更新后重新提交
		existingCar.CarPlate = in.CarPlate
		existingCar.CarModel = in.CarModel
		existingCar.DrivingLicense = in.DrivingLicense
		existingCar.Status = 2 // 审核中
		existingCar.RejectReason = ""
		err = global.DB.Save(&existingCar).Error
		if err != nil {
			l.Logger.Errorf("更新车辆认证记录失败: %v", err)
			return nil, errors.New("车辆认证提交失败")
		}
		l.Logger.Infof("司机 %d 重新提交车辆认证", driverId)
		return &driver.CartAuthResp{Msg: "车辆认证信息已重新提交，请等待审核"}, nil
	}

	// 首次提交车辆认证
	car := user3.DriverCar{
		DriverID:       int64(driverId),
		CarPlate:       in.CarPlate,
		CarModel:       in.CarModel,
		CarColor:       "未填写",           // 默认值
		DrivingLicense: in.DrivingLicense,
		CarImg:         in.DrivingLicense, // 默认用行驶证图片
		Status:         2,                 // 2-审核中
		RejectReason:   "",
	}

	err = car.CreateData(global.DB)
	if err != nil {
		l.Logger.Errorf("创建车辆认证记录失败: %v", err)
		return nil, errors.New("车辆认证提交失败")
	}

	l.Logger.Infof("司机 %d 首次提交车辆认证", driverId)
	return &driver.CartAuthResp{Msg: "车辆认证提交成功，请等待审核"}, nil
}
