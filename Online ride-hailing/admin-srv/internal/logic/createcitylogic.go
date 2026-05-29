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

type CreateCityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCityLogic {
	return &CreateCityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCityLogic) CreateCity(in *admin.CreateCityRequest) (*admin.CreateCityResponse, error) {
	// 验证必填参数
	if in.CityCode == "" {
		return nil, errors.New("城市编码不能为空")
	}
	if in.CityName == "" {
		return nil, errors.New("城市名称不能为空")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 检查城市编码是否已存在
	var count int64
	global.DB.Model(&system.CityConfig{}).Where("city_code = ?", in.CityCode).Count(&count)
	if count > 0 {
		return nil, errors.New("城市编码已存在")
	}

	// 创建城市
	city := &system.CityConfig{
		CityCode: in.CityCode,
		CityName: in.CityName,
		Status:   int8(in.Status),
	}

	result := global.DB.Create(city)
	if result.Error != nil {
		return nil, result.Error
	}

	return &admin.CreateCityResponse{
		Id: int32(city.ID),
	}, nil
}
