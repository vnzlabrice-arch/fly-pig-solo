package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"admin-srv/global"
	"admin-srv/internal/svc"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLogoutLogic {
	return &AdminLogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminLogoutLogic) AdminLogout(in *admin.AdminLogoutRequest) (*admin.AdminLogoutResponse, error) {
	// 参数校验
	if in.AdminId <= 0 {
		return nil, errors.New("管理员ID无效")
	}

	// 登出逻辑：清除Redis中的Token等
	// 简单实现，实际可添加Token黑名单等逻辑
	if global.RDB != nil {
		global.RDB.Del(l.ctx, "admin:token:"+strconv.FormatInt(int64(in.AdminId), 10))
	}

	fmt.Sprintf("管理员 %d 登出成功", in.AdminId)

	return &admin.AdminLogoutResponse{
		Success: true,
	}, nil
}
