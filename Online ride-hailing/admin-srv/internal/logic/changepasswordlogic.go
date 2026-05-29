package logic

import (
	"admin-srv/global"
	"admin-srv/pkg"
	"context"
	"errors"

	"admin-srv/internal/svc"
	"admin-srv/model/system"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *admin.ChangePasswordRequest) (*admin.ChangePasswordResponse, error) {
	// 参数校验
	if in.AdminId <= 0 {
		return nil, errors.New("管理员ID无效")
	}
	if in.OldPassword == "" {
		return nil, errors.New("旧密码不能为空")
	}
	if in.NewPassword == "" {
		return nil, errors.New("新密码不能为空")
	}
	if len(in.NewPassword) < 6 {
		return nil, errors.New("新密码长度不能少于6位")
	}

	// 查询管理员
	var adminUser system.AdminUser
	if err := global.DB.First(&adminUser, in.AdminId).Error; err != nil {
		return nil, errors.New("管理员不存在")
	}

	// 验证旧密码
	if pkg.MD5(in.OldPassword) != adminUser.Password {
		return nil, errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword := pkg.MD5(in.NewPassword)

	// 更新密码
	if err := global.DB.Model(&adminUser).Update("password", hashedPassword).Error; err != nil {
		return nil, errors.New("密码更新失败")
	}

	return &admin.ChangePasswordResponse{
		Success: true,
		Message: "密码修改成功",
	}, nil
}
