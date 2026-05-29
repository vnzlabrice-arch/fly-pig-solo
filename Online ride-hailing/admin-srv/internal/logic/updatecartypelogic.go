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

type UpdateCarTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCarTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCarTypeLogic {
	return &UpdateCarTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCarTypeLogic) UpdateCarType(in *admin.UpdateCarTypeRequest) (*admin.UpdateCarTypeResponse, error) {
	// 验证必填参数
	if in.Id == 0 {
		return nil, errors.New("车型ID不能为空")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 检查车型是否存在
	var carType system.CarTypeConfig
	result := global.DB.First(&carType, in.Id)
	if result.Error != nil {
		return nil, errors.New("车型不存在")
	}

	// 更新车型
	updates := map[string]interface{}{}
	if in.TypeName != "" {
		updates["type_name"] = in.TypeName
	}
	if in.BasePrice > 0 {
		updates["base_price"] = in.BasePrice
	}
	if in.KmPrice > 0 {
		updates["km_price"] = in.KmPrice
	}
	if in.MinutePrice > 0 {
		updates["minute_price"] = in.MinutePrice
	}
	if in.Status != 0 {
		updates["status"] = in.Status
	}

	if len(updates) > 0 {
		global.DB.Model(&carType).Updates(updates)
	}

	return &admin.UpdateCarTypeResponse{
		Success: true,
	}, nil
}
