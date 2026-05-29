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

type GetCurrentAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCurrentAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentAdminLogic {
	return &GetCurrentAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCurrentAdminLogic) GetCurrentAdmin(in *admin.GetCurrentAdminRequest) (*admin.GetCurrentAdminResponse, error) {
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
	result := global.DB.First(&adminUser, in.AdminId)
	if result.Error != nil {
		return nil, errors.New("管理员不存在")
	}

	// 查询角色
	var role system.AdminRole
	roleName := ""
	if err := global.DB.First(&role, adminUser.RoleID).Error; err == nil {
		roleName = role.Name
	}

	return &admin.GetCurrentAdminResponse{
		Id:       int32(adminUser.ID),
		Username: adminUser.Username,
		RoleId:   int32(adminUser.RoleID),
		Status:   int32(adminUser.Status),
		RoleName: roleName,
	}, nil
}
