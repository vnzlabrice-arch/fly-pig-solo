package logic

import (
	"admin-srv/pkg"
	"context"
	"errors"

	"admin-srv/global"
	"admin-srv/internal/svc"
	"admin-srv/model/system"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAdminLogic) CreateAdmin(in *admin.CreateAdminRequest) (*admin.CreateAdminResponse, error) {
	// 参数校验
	if in.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if in.Password == "" {
		return nil, errors.New("密码不能为空")
	}
	if in.RoleId <= 0 {
		return nil, errors.New("请选择角色")
	}

	// 检查数据库连接
	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 检查用户名是否已存在
	var existing system.AdminUser
	if err := global.DB.Where("username = ?", in.Username).First(&existing).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查角色是否存在
	var role system.AdminRole
	if err := global.DB.First(&role, in.RoleId).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	// 加密密码（MD5）
	hashedPassword := pkg.MD5(in.Password)

	// 创建管理员
	adminUser := system.AdminUser{
		Username: in.Username,
		Password: hashedPassword,
		RoleID:   int64(in.RoleId),
		Status:   1,
	}

	if err := global.DB.Create(&adminUser).Error; err != nil {
		return nil, errors.New("创建管理员失败")
	}

	return &admin.CreateAdminResponse{
		Id:      int32(adminUser.ID),
		Success: true,
	}, nil
}
