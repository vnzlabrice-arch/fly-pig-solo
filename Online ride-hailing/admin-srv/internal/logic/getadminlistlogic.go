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

type GetAdminListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminListLogic {
	return &GetAdminListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminListLogic) GetAdminList(in *admin.GetAdminListRequest) (*admin.GetAdminListResponse, error) {
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
	query := global.DB.Model(&system.AdminUser{})

	// 关键词搜索
	if in.Keyword != "" {
		query = query.Where("username LIKE ?", "%"+in.Keyword+"%")
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 查询列表
	var adminUsers []system.AdminUser
	result := query.Offset(int(offset)).Limit(int(pageSize)).Order("id ASC").Find(&adminUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	// 获取角色名称映射
	roleMap := make(map[int64]string)
	var roles []system.AdminRole
	global.DB.Find(&roles)
	for _, role := range roles {
		roleMap[role.ID] = role.Name
	}

	// 转换结果
	list := make([]*admin.AdminListItem, 0, len(adminUsers))
	for _, user := range adminUsers {
		item := &admin.AdminListItem{
			Id:       int32(user.ID),
			Username: user.Username,
			RoleId:   int32(user.RoleID),
			RoleName: roleMap[user.RoleID],
			Status:   int32(user.Status),
		}

		// 格式化时间
		if user.LastLoginTime != nil {
			item.LastLoginTime = user.LastLoginTime.Format("2006-01-02 15:04:05")
		}

		list = append(list, item)
	}

	return &admin.GetAdminListResponse{
		Total: int32(total),
		List:  list,
	}, nil
}
