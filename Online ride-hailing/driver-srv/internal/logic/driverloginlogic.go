package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"
	"driver-srv/pkg"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DriverLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDriverLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverLoginLogic {
	return &DriverLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DriverLoginLogic) DriverLogin(in *driver.DriverLoginReq) (*driver.DriverLoginResp, error) {
	// 1. 参数校验
	if in.Phone == "" {
		return nil, errors.New("手机号不能为空")
	}
	if in.Code == "" {
		return nil, errors.New("验证码不能为空")
	}

	// 2. 防暴力重试（1分钟3次错误锁定）
	limitKey := fmt.Sprintf("login:limit:%s", in.Phone)
	i, _ := global.RDB.Incr(global.Ctx, limitKey).Result()
	global.RDB.Expire(global.Ctx, limitKey, time.Minute*1)
	if i > 3 {
		return nil, errors.New("登录失败次数过多，请1分钟后再试")
	}

	// 3. 验证码校验（一次性消费，用过立即删除）
	result, err := global.RDB.Get(global.Ctx, "send"+in.Phone).Result()
	if err != nil || result != in.Code {
		return nil, errors.New("验证码错误或已过期")
	}
	global.RDB.Del(global.Ctx, "send"+in.Phone)

	// 4. 查询司机 + 自动注册（注册登录一体化）
	var driverModel user3.DriverUser
	err = global.DB.Where("phone = ?", in.Phone).First(&driverModel).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Errorf("查询司机失败: %v", err)
			return nil, errors.New("系统异常，请稍后重试")
		}
		// 司机不存在 → 自动注册
		newDriver := user3.DriverUser{
			Phone:     in.Phone,
			Nickname:  "",
			AvatarURL: "",
			Password:  "",
			Status:    2, // 默认为离线
		}
		err = newDriver.CreateData(global.DB)
		if err != nil {
			l.Logger.Errorf("创建司机失败: %v", err)
			return nil, errors.New("注册失败")
		}
		driverModel = newDriver

		// 新司机自动初始化钱包
		wallet := user3.DriverWallet{
			DriverID:     int64(driverModel.ID),
			Balance:      0,
			Withdrawable: 0,
			Frozen:       0,
			TotalIncome:  0,
		}
		err = wallet.CreateData(global.DB)
		if err != nil {
			l.Logger.Errorf("创建钱包失败: %v", err)
			return nil, errors.New("钱包初始化失败")
		}
	}

	// 5. 同一账号互踢：旧token加入黑名单
	oldKey := fmt.Sprintf("driver:token:%d", driverModel.ID)
	oldToken, _ := global.RDB.Get(global.Ctx, oldKey).Result()
	if oldToken != "" {
		// 使用统一的jwt黑名单机制（与退出登录保持一致）
		_ = pkg.BlacklistToken(oldToken, pkg.APP_KEY)
	}

	// 6. 签发新token
	token, err := pkg.TokenHandler(strconv.FormatInt(int64(driverModel.ID), 10))
	if err != nil {
		return nil, errors.New("token签发失败")
	}
	// 将最新token存入Redis，7天有效
	global.RDB.Set(global.Ctx, oldKey, token, time.Hour*24*7)

	// 7. 查询认证状态和车辆认证状态
	var authModel user3.DriverAuth
	var cartModel user3.DriverCar
	global.DB.Where("driver_id = ?", driverModel.ID).First(&authModel)
	global.DB.Where("driver_id = ?", driverModel.ID).First(&cartModel)

	// 8. 登录成功：清除限流计数器 + 状态改为在线
	global.RDB.Del(global.Ctx, limitKey)
	global.DB.Model(&driverModel).Update("status", 1)

	l.Logger.Infof("司机 %d 登录成功", driverModel.ID)
	return &driver.DriverLoginResp{
		DriverId:    int64(driverModel.ID),
		Token:       token,
		AuditStatus: int64(authModel.AuditStatus),
		CartStatus:  int64(cartModel.Status),
	}, nil
}
