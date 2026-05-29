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

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *admin.UpdateRoleRequest) (*admin.UpdateRoleResponse, error) {
	// 参数校验
	if in.RoleId <= 0 {
		return nil, errors.New("角色ID无效")
	}
	if in.Name == "" {
		return nil, errors.New("角色名称不能为空")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 查询角色
	var role system.AdminRole
	if err := global.DB.First(&role, in.RoleId).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	// 检查角色名称是否被其他角色使用
	var existing system.AdminRole
	if err := global.DB.Where("name = ? AND id != ?", in.Name, in.RoleId).First(&existing).Error; err == nil {
		return nil, errors.New("角色名称已存在")
	}

	// 更新角色
	updates := map[string]interface{}{
		"name":   in.Name,
		"remark": in.Remark,
	}
	if err := global.DB.Model(&role).Updates(updates).Error; err != nil {
		return nil, errors.New("更新角色失败")
	}

	// 更新角色菜单 - 先删除旧的关联
	global.DB.Where("role_id = ?", in.RoleId).Delete(&system.AdminRoleMenu{})

	// 添加新的菜单关联
	for _, menuID := range in.MenuIds {
		roleMenu := system.AdminRoleMenu{
			RoleID: role.ID,
			MenuID: int64(menuID),
		}
		global.DB.Create(&roleMenu)
	}

	return &admin.UpdateRoleResponse{
		Success: true,
	}, nil
}
