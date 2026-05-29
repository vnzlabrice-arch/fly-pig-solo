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

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type OutLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOutLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OutLoginLogic {
	return &OutLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OutLoginLogic) OutLogin(in *driver.OutLoginReq) (*driver.OutLoginResp, error) {
	// 1. 从RPC上下文中获取token
	token, err := l.svcCtx.GetRpcToken(l.ctx)
	if err != nil {
		return nil, errors.New("未登录或token已失效")
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

	// 3. 将token加入黑名单（使用JWT剩余有效期作为TTL，防止已过期token重复使用）
	err = pkg.BlacklistToken(token, pkg.APP_KEY)
	if err != nil {
		return nil, err
	}

	// 4. 删除Redis中的driver:token缓存，清空登录状态
	key := fmt.Sprintf("driver:token:%d", driverId)
	global.RDB.Del(global.Ctx, key)

	// 5. 将司机状态更新为离线（2-离线）
	err = global.DB.Model(&user3.DriverUser{}).Where("id = ?", driverId).Update("status", 2).Error
	if err != nil {
		l.Logger.Errorf("更新司机离线状态失败: %v", err)
		return nil, errors.New("退出登录失败")
	}

	l.Logger.Infof("司机 %d 退出登录成功", driverId)
	return &driver.OutLoginResp{
		Success: true,
	}, nil
}
