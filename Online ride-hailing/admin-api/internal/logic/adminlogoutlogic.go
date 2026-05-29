// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"admin-api/internal/svc"
	"admin-api/internal/types"
	admin "admin-srv/adminclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLogoutLogic {
	return &AdminLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLogoutLogic) AdminLogout(req *types.AdminLogoutRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.AdminLogout(l.ctx, &admin.AdminLogoutRequest{
		AdminId: req.AdminId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "登出失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "登出成功",
		Data: map[string]interface{}{
			"success": rpcResp.Success,
		},
	}, nil
}
