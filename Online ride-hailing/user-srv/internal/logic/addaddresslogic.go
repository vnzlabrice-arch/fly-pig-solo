package logic

import (
	"context"
	"errors"
	"user-srv/global"
	"user-srv/internal/svc"
	"user-srv/model"
	"user-srv/user"

	"user-srv/pkg"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAddressLogic) AddAddress(in *user.AddAddressReq) (*user.AddAddressResp, error) {
	if in.UserId == 0 {
		return &user.AddAddressResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

	if in.Address == "" {
		return &user.AddAddressResp{
			Code:    400,
			Message: "地址不能为空",
		}, nil
	}

	// 自动地址转经纬度
	lng, lat, geoErr := pkg.AddressToLngLat(in.Address)
	if geoErr != nil {
		l.Errorf("地址转经纬度失败: %v", geoErr)
		return &user.AddAddressResp{
			Code:    400,
			Message: "地址解析失败，请检查地址是否正确",
		}, nil
	}

	var passengerUser model.PassengerUser
	err := global.DB.Where("id = ?", in.UserId).First(&passengerUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.AddAddressResp{
				Code:    404,
				Message: "用户不存在",
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &user.AddAddressResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	if in.IsDefault {
		err = global.DB.Model(&model.PassengerAddressBook{}).
			Where("passenger_id = ?", in.UserId).
			Update("is_default", 0).Error
		if err != nil {
			l.Errorf("重置默认地址失败: %v", err)
			return &user.AddAddressResp{
				Code:    500,
				Message: "系统错误",
			}, nil
		}
	}

	isDefaultValue := int8(0)
	if in.IsDefault {
		isDefaultValue = 1
	}

	addressBook := model.PassengerAddressBook{
		PassengerID: in.UserId,
		Tag:         in.Tag,
		Address:     in.Address,
		Lng:         lng,
		Lat:         lat,
		IsDefault:   isDefaultValue,
	}

	err = global.DB.Create(&addressBook).Error
	if err != nil {
		l.Errorf("创建地址失败: %v", err)
		return &user.AddAddressResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}
	return &user.AddAddressResp{
		Code:      200,
		Message:   "地址添加成功",
		AddressId: addressBook.ID,
	}, nil
}
