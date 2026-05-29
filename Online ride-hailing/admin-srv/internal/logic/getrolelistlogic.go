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

type GetRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleListLogic) GetRoleList(in *admin.GetRoleListRequest) (*admin.GetRoleListResponse, error) {
	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 默认分页参数
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 构建查询
	query := global.DB.Model(&system.AdminRole{})

	// 关键词搜索
	if in.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+in.Keyword+"%")
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 查询列表
	var roles []system.AdminRole
	result := query.Offset(int(offset)).Limit(int(pageSize)).Order("id ASC").Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}

	// 获取角色菜单映射
	roleMenuMap := make(map[int64][]int64)
	var roleMenus []system.AdminRoleMenu
	global.DB.Find(&roleMenus)
	for _, rm := range roleMenus {
		roleMenuMap[rm.RoleID] = append(roleMenuMap[rm.RoleID], rm.MenuID)
	}

	// 转换结果
	list := make([]*admin.RoleListItem, 0, len(roles))
	for _, role := range roles {
		menuIDs := roleMenuMap[role.ID]
		var menuID int32
		if len(menuIDs) > 0 {
			menuID = int32(menuIDs[0]) // 取第一个菜单ID
		}

		item := &admin.RoleListItem{
			Id:      int32(role.ID),
			Name:    role.Name,
			Remark:  role.Remark,
			MenuIds: menuID,
		}
		list = append(list, item)
	}

	return &admin.GetRoleListResponse{
		Total: int32(total),
		List:  list,
	}, nil
}
