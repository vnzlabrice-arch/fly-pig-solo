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

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(in *admin.DeleteRoleRequest) (*admin.DeleteRoleResponse, error) {
	// 参数校验
	if in.RoleId <= 0 {
		return nil, errors.New("角色ID无效")
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

	// 检查是否有管理员使用该角色
	var adminCount int64
	global.DB.Model(&system.AdminUser{}).Where("role_id = ?", in.RoleId).Count(&adminCount)
	if adminCount > 0 {
		return nil, errors.New("该角色下有管理员，无法删除")
	}

	// 删除角色菜单关联
	global.DB.Where("role_id = ?", in.RoleId).Delete(&system.AdminRoleMenu{})

	// 删除角色
	if err := global.DB.Delete(&role).Error; err != nil {
		return nil, errors.New("删除角色失败")
	}

	return &admin.DeleteRoleResponse{
		Success: true,
	}, nil
}
