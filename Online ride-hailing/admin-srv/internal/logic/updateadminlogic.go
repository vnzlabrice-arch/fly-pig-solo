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

type UpdateAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminLogic {
	return &UpdateAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAdminLogic) UpdateAdmin(in *admin.UpdateAdminRequest) (*admin.UpdateAdminResponse, error) {
	// 参数校验
	if in.AdminId <= 0 {
		return nil, errors.New("管理员ID无效")
	}
	if in.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if in.RoleId <= 0 {
		return nil, errors.New("请选择角色")
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

	// 检查用户名是否被其他管理员使用
	var existing system.AdminUser
	if err := global.DB.Where("username = ? AND id != ?", in.Username, in.AdminId).First(&existing).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查角色是否存在
	var role system.AdminRole
	if err := global.DB.First(&role, in.RoleId).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	// 更新管理员
	updates := map[string]interface{}{
		"username": in.Username,
		"role_id":  int64(in.RoleId),
		"status":   in.Status,
	}

	if err := global.DB.Model(&adminUser).Updates(updates).Error; err != nil {
		return nil, errors.New("更新管理员失败")
	}

	return &admin.UpdateAdminResponse{
		Success: true,
	}, nil
}
