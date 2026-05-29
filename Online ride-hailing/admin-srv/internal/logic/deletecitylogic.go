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

type DeleteCityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCityLogic {
	return &DeleteCityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCityLogic) DeleteCity(in *admin.DeleteCityRequest) (*admin.DeleteCityResponse, error) {
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

	// 删除城市
	result = global.DB.Delete(&city)
	if result.Error != nil {
		return nil, errors.New("删除城市失败")
	}

	return &admin.DeleteCityResponse{
		Success: true,
	}, nil
}
