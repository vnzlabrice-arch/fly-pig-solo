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

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(in *admin.UpdateMenuRequest) (*admin.UpdateMenuResponse, error) {
	// 参数校验
	if in.MenuId <= 0 {
		return nil, errors.New("菜单ID无效")
	}
	if in.Name == "" {
		return nil, errors.New("菜单名称不能为空")
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

	// 检查父级菜单是否存在（不能将自己设为父级）
	if in.ParentId == in.MenuId {
		return nil, errors.New("不能将自身设为父级菜单")
	}
	if in.ParentId > 0 {
		var parentMenu system.AdminMenu
		if err := global.DB.First(&parentMenu, in.ParentId).Error; err != nil {
			return nil, errors.New("父级菜单不存在")
		}
	}

	// 更新菜单
	updates := map[string]interface{}{
		"parent_id": int64(in.ParentId),
		"name":      in.Name,
		"path":      in.Path,
		"icon":      in.Icon,
		"sort":      int(in.Sort),
	}

	if err := global.DB.Model(&menu).Updates(updates).Error; err != nil {
		return nil, errors.New("更新菜单失败")
	}

	return &admin.UpdateMenuResponse{
		Success: true,
	}, nil
}
