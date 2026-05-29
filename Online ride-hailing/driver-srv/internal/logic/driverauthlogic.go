package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"
	"driver-srv/pkg"
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DriverAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDriverAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverAuthLogic {
	return &DriverAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DriverAuthLogic) DriverAuth(in *driver.DriverAuthReq) (*driver.DriverAuthResp, error) {
	// 1. 校验token
	token := in.Token
	if token == "" {
		return nil, errors.New("token不能为空")
	}

	// 校验token是否在黑名单里（使用pkg统一的黑名单检查）
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

	// 3. 校验参数
	if in.RealName == "" || in.IDCard == "" || in.IDCardBack == "" || in.IDCardFront == "" || in.LicenseImg == "" {
		return nil, errors.New("姓名、身份证号、证件图片不能为空")
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

	// 5. 检查是否已有认证记录
	var existingAuth user3.DriverAuth
	hasExistingRecord := false
	err = global.DB.Where("driver_id = ?", driverId).First(&existingAuth).Error
	if err == nil {
		hasExistingRecord = true
		// 已有认证记录，根据状态分别处理
		switch existingAuth.AuditStatus {
		case 1:
			return nil, errors.New("您的认证正在审核中，请耐心等待")
		case 2:
			return &driver.DriverAuthResp{Msg: "您已完成实名认证"}, nil
		case 3:
			// 之前被驳回，允许重新提交
			break
		default:
			break
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Errorf("查询认证记录失败: %v", err)
		return nil, errors.New("系统异常，请稍后重试")
	}

	// 6. 校验身份证实名信息（调用第三方API）
	fen, err := pkg.ShenFen(in.RealName, in.IDCard)
	if err != nil {
		l.Logger.Errorf("身份校验调用失败: %v", err)
		return nil, errors.New("身份校验服务异常，请稍后重试")
	}
	if !fen {
		return nil, errors.New("身份信息不匹配，请检查姓名和身份证号")
	}

	// 7. 保存或更新认证记录
	if hasExistingRecord {
		// 之前被驳回的记录，更新后重新提交
		existingAuth.RealName = in.RealName
		existingAuth.IDCard = in.IDCard
		existingAuth.IDCardFront = in.IDCardFront
		existingAuth.IDCardBack = in.IDCardBack
		existingAuth.LicenseImg = in.LicenseImg
		existingAuth.FaceImg = in.IDCardFront // 默认使用身份证正面作为人脸核验照
		existingAuth.AuditStatus = 1          // 重新提交，状态改为待审核
		existingAuth.Reason = ""
		err = existingAuth.UpdateData(global.DB)
		if err != nil {
			l.Logger.Errorf("更新认证记录失败: %v", err)
			return nil, errors.New("认证信息提交失败")
		}
		l.Logger.Infof("司机 %d 重新提交认证", driverId)
		return &driver.DriverAuthResp{Msg: "认证信息已重新提交，请等待审核"}, nil
	}

	// 首次提交认证
	auth := user3.DriverAuth{
		DriverID:    int64(driverId),
		RealName:    in.RealName,
		IDCard:      in.IDCard,
		IDCardFront: in.IDCardFront,
		IDCardBack:  in.IDCardBack,
		LicenseImg:  in.LicenseImg,
		FaceImg:     in.IDCardFront, // 默认使用身份证正面作为人脸核验照
		AuditStatus: 1,              // 1-待审核
		Reason:      "",
	}

	err = auth.CreateAuth(global.DB)
	if err != nil {
		l.Logger.Errorf("创建认证记录失败: %v", err)
		return nil, errors.New("认证信息提交失败")
	}

	l.Logger.Infof("司机 %d 首次提交认证", driverId)
	return &driver.DriverAuthResp{Msg: "认证信息提交成功，请等待审核"}, nil
}
