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

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMenuLogic) CreateMenu(in *admin.CreateMenuRequest) (*admin.CreateMenuResponse, error) {
	// 参数校验
	if in.Name == "" {
		return nil, errors.New("菜单名称不能为空")
	}

	if in.Path == "" {
		return nil, errors.New("菜单路径不能为空")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 检查父级菜单是否存在
	if in.ParentId > 0 {
		var parentMenu system.AdminMenu
		if err := global.DB.First(&parentMenu, in.ParentId).Error; err != nil {
			return nil, errors.New("父级菜单不存在")
		}
	}

	// 创建菜单
	menu := system.AdminMenu{
		ParentID: int64(in.ParentId),
		Name:     in.Name,
		Path:     in.Path,
		Icon:     in.Icon,
		Sort:     int(in.Sort),
	}

	if err := global.DB.Create(&menu).Error; err != nil {
		return nil, errors.New("创建菜单失败")
	}

	return &admin.CreateMenuResponse{
		Id:      int32(menu.ID),
		Success: true,
	}, nil
}
