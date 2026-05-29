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

type GetAdminDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminDetailLogic {
	return &GetAdminDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminDetailLogic) GetAdminDetail(in *admin.GetAdminDetailRequest) (*admin.GetAdminDetailResponse, error) {
	// 参数校验
	if in.AdminId <= 0 {
		return nil, errors.New("管理员ID无效")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 查询管理员
	var adminUser system.AdminUser
	result := global.DB.First(&adminUser, in.AdminId)
	if result.Error != nil {
		return nil, errors.New("管理员不存在")
	}

	// 查询角色
	var role system.AdminRole
	roleName := ""
	if err := global.DB.First(&role, adminUser.RoleID).Error; err == nil {
		roleName = role.Name
	}

	// 格式化响应
	lastLoginTime := ""
	if adminUser.LastLoginTime != nil {
		lastLoginTime = adminUser.LastLoginTime.Format("2006-01-02 15:04:05")
	}

	return &admin.GetAdminDetailResponse{
		Id:            int32(adminUser.ID),
		Username:      adminUser.Username,
		RoleId:        int32(adminUser.RoleID),
		RoleName:      roleName,
		Status:        int32(adminUser.Status),
		LastLoginTime: lastLoginTime,
	}, nil
}
