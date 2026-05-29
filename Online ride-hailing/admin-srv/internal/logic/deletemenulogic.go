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

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(in *admin.DeleteMenuRequest) (*admin.DeleteMenuResponse, error) {
	// 参数校验
	if in.MenuId <= 0 {
		return nil, errors.New("菜单ID无效")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 查询菜单
	var menu system.AdminMenu
	if err := global.DB.First(&menu, in.MenuId).Error; err != nil {
		return nil, errors.New("菜单不存在")
	}

	// 检查是否有子菜单
	var childCount int64
	global.DB.Model(&system.AdminMenu{}).Where("parent_id = ?", in.MenuId).Count(&childCount)
	if childCount > 0 {
		return nil, errors.New("该菜单下有子菜单，无法删除")
	}

	// 检查是否有角色使用该菜单
	global.DB.Where("menu_id = ?", in.MenuId).Delete(&system.AdminRoleMenu{})

	// 删除菜单
	if err := global.DB.Delete(&menu).Error; err != nil {
		return nil, errors.New("删除菜单失败")
	}

	return &admin.DeleteMenuResponse{
		Success: true,
	}, nil
}
