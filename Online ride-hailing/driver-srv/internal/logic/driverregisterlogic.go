package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type DriverRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDriverRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverRegisterLogic {
	return &DriverRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DriverRegisterLogic) DriverRegister(in *driver.DriverRegisterReq) (*driver.DriverRegisterResp, error) {
	if in.Code == "" {
		return nil, errors.New("验证码不能为空")
	}

	if in.Phone == "" {
		return nil, errors.New("手机号不能为空")
	}

	result, err := global.RDB.Get(global.Ctx, "send"+in.Phone).Result()
	if err != nil || result != in.Code {
		return nil, errors.New("验证码错误或已过期")
	}

	var count int64
	err = global.DB.Model(&user3.DriverUser{}).Where("phone = ?", in.Phone).Count(&count).Error
	if err != nil {
		return nil, errors.New("数据查询失败")
	}
	if count > 0 {
		return nil, errors.New("手机号已注册")
	}

	data := user3.DriverUser{
		Phone:    in.Phone,
		Nickname: in.Nickname,
		Status:   1,
	}
	err = global.DB.Create(&data).Error
	if err != nil {
		return nil, errors.New("司机注册失败")
	}

	return &driver.DriverRegisterResp{
		DriverId: int64(data.ID),
	}, nil
}
