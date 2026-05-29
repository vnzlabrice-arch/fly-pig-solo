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

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleLogic) CreateRole(in *admin.CreateRoleRequest) (*admin.CreateRoleResponse, error) {
	// 参数校验
	if in.Name == "" {
		return nil, errors.New("角色名称不能为空")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 检查角色名称是否已存在
	var existing system.AdminRole
	if err := global.DB.Where("name = ?", in.Name).First(&existing).Error; err == nil {
		return nil, errors.New("角色名称已存在")
	}

	// 创建角色
	role := system.AdminRole{
		Name:   in.Name,
		Remark: in.Remark,
	}

	if err := global.DB.Create(&role).Error; err != nil {
		return nil, errors.New("创建角色失败")
	}

	// 关联菜单
	for _, menuID := range in.MenuIds {
		roleMenu := system.AdminRoleMenu{
			RoleID: role.ID,
			MenuID: int64(menuID),
		}
		global.DB.Create(&roleMenu)
	}

	return &admin.CreateRoleResponse{
		Id:      int32(role.ID),
		Success: true,
	}, nil
}
