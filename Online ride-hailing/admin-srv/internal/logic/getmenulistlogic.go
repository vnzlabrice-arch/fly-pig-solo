package logic

import (
	"admin-srv/global"
	"context"

	"admin-srv/internal/svc"
	"admin-srv/model/system"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuListLogic {
	return &GetMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuListLogic) GetMenuList(in *admin.GetMenuListRequest) (*admin.GetMenuListResponse, error) {
	// 查询所有菜单
	var menus []system.AdminMenu
	result := global.DB.Order("sort ASC, id ASC").Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	// 构建树形结构
	tree := buildMenuTree(menus, 0)

	return &admin.GetMenuListResponse{
		List: tree,
	}, nil
}

// buildMenuTree 构建菜单树形结构
func buildMenuTree(menus []system.AdminMenu, parentId int64) []*admin.MenuListItem {
	tree := make([]*admin.MenuListItem, 0)
	for _, menu := range menus {
		if menu.ParentID == parentId {
			item := &admin.MenuListItem{
				Id:       int32(menu.ID),
				ParentId: int32(menu.ParentID),
				Name:     menu.Name,
				Path:     menu.Path,
				Icon:     menu.Icon,
				Sort:     int32(menu.Sort),
				Children: buildMenuTree(menus, menu.ID),
			}
			// 如果没有子菜单，设置为 nil 而不是空切片
			if len(item.Children) == 0 {
				item.Children = nil
			}
			tree = append(tree, item)
		}
	}
	return tree
}
