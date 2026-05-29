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

type GetMenuDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuDetailLogic {
	return &GetMenuDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuDetailLogic) GetMenuDetail(in *admin.GetMenuDetailRequest) (*admin.GetMenuDetailResponse, error) {
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
	result := global.DB.First(&menu, in.MenuId)
	if result.Error != nil {
		return nil, errors.New("菜单不存在")
	}

	return &admin.GetMenuDetailResponse{
		Id:       int32(menu.ID),
		ParentId: int32(menu.ParentID),
		Name:     menu.Name,
		Path:     menu.Path,
		Icon:     menu.Icon,
		Sort:     int32(menu.Sort),
	}, nil
}
