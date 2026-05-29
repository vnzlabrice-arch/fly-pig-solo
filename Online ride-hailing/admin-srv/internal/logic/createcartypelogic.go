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

type CreateCarTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCarTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCarTypeLogic {
	return &CreateCarTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCarTypeLogic) CreateCarType(in *admin.CreateCarTypeRequest) (*admin.CreateCarTypeResponse, error) {
	// 验证必填参数
	if in.TypeName == "" {
		return nil, errors.New("车型名称不能为空")
	}
	if in.BasePrice <= 0 {
		return nil, errors.New("起步价必须大于0")
	}
	if in.KmPrice <= 0 {
		return nil, errors.New("公里单价必须大于0")
	}
	if in.MinutePrice <= 0 {
		return nil, errors.New("时长单价必须大于0")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 创建车型
	carType := &system.CarTypeConfig{
		TypeName:    in.TypeName,
		BasePrice:   in.BasePrice,
		KmPrice:     in.KmPrice,
		MinutePrice: in.MinutePrice,
		Status:      int8(in.Status),
	}

	err := global.DB.Create(carType)
	if err != nil {
		return nil, errors.New("创建车型失败")
	}

	return &admin.CreateCarTypeResponse{
		Id: int32(carType.ID),
	}, nil
}
