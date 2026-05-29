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

type DeleteCarTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCarTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCarTypeLogic {
	return &DeleteCarTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCarTypeLogic) DeleteCarType(in *admin.DeleteCarTypeRequest) (*admin.DeleteCarTypeResponse, error) {
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
	err := global.DB.First(&carType, in.Id)
	if err != nil {
		return nil, errors.New("车型不存在")
	}

	// 删除车型
	err = global.DB.Delete(&carType)
	if err != nil {
		return nil, errors.New("删除车型失败")
	}

	return &admin.DeleteCarTypeResponse{
		Success: true,
	}, nil
}
