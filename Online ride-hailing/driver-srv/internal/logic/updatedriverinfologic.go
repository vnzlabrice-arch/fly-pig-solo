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

type UpdateDriverInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDriverInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDriverInfoLogic {
	return &UpdateDriverInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDriverInfoLogic) UpdateDriverInfo(in *driver.UpdateDriverInfoReq) (*driver.UpdateDriverInfoResp, error) {
	//从token中获取id
	//driverID, ok := l.ctx.Value("driver_id").(int64)
	//if !ok {
	//	return nil, errors.New("invalid token")
	//}

	user := user3.DriverUser{
		Model:     gorm.Model{ID: uint(in.DriverID)},
		Nickname:  in.Nickname,
		AvatarURL: in.AvatarURL,
	}

	err := global.DB.Updates(&user).Error
	if err != nil {
		return nil, errors.New("修改失败")
	}

	return &driver.UpdateDriverInfoResp{
		Success: true,
	}, nil
}
