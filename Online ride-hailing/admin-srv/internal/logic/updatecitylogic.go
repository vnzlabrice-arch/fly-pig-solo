package logic

import (
	"context"
	"errors"

	"admin-srv/global"
	"admin-srv/internal/svc"
	"admin-srv/model/system"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCityLogic {
	return &UpdateCityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCityLogic) UpdateCity(in *admin.UpdateCityRequest) (*admin.UpdateCityResponse, error) {
	// 验证必填参数
	if in.Id == 0 {
		return nil, errors.New("城市ID不能为空")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 检查城市是否存在
	var city system.CityConfig
	result := global.DB.First(&city, in.Id)
	if result.Error != nil {
		return nil, errors.New("城市不存在")
	}

	// 检查城市编码是否与其他城市冲突
	if in.CityCode != "" && in.CityCode != city.CityCode {
		var count int64
		global.DB.Model(&system.CityConfig{}).Where("city_code = ? AND id != ?", in.CityCode, in.Id).Count(&count)
		if count > 0 {
			return nil, errors.New("城市编码已存在")
		}
	}

	// 更新城市
	updates := map[string]interface{}{}
	if in.CityCode != "" {
		updates["city_code"] = in.CityCode
	}
	if in.CityName != "" {
		updates["city_name"] = in.CityName
	}
	if in.Status != 0 {
		updates["status"] = in.Status
	}

	if len(updates) > 0 {
		err := global.DB.Model(&city).Updates(updates)
		if err != nil {
			return nil, errors.New("更新城市失败")
		}
	}

	return &admin.UpdateCityResponse{
		Success: true,
	}, nil
}
