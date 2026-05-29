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

type DeleteAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdminLogic {
	return &DeleteAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAdminLogic) DeleteAdmin(in *admin.DeleteAdminRequest) (*admin.DeleteAdminResponse, error) {
	// 参数校验
	if in.AdminId <= 0 {
		return nil, errors.New("管理员ID无效")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 查询管理员
	var adminUser system.AdminUser
	if err := global.DB.First(&adminUser, in.AdminId).Error; err != nil {
		return nil, errors.New("管理员不存在")
	}

	// 删除管理员
	if err := global.DB.Delete(&adminUser).Error; err != nil {
		return nil, errors.New("删除管理员失败")
	}

	return &admin.DeleteAdminResponse{
		Success: true,
	}, nil
}
