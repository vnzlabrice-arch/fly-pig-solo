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

type GetRoleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleDetailLogic {
	return &GetRoleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleDetailLogic) GetRoleDetail(in *admin.GetRoleDetailRequest) (*admin.GetRoleDetailResponse, error) {
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
	result := global.DB.First(&role, in.RoleId)
	if result.Error != nil {
		return nil, errors.New("角色不存在")
	}

	// 查询角色菜单
	var roleMenus []system.AdminRoleMenu
	global.DB.Where("role_id = ?", in.RoleId).Find(&roleMenus)

	menuIDs := make([]int32, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, int32(rm.MenuID))
	}

	return &admin.GetRoleDetailResponse{
		Id:      int32(role.ID),
		Name:    role.Name,
		Remark:  role.Remark,
		MenuIds: menuIDs,
	}, nil
}
