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

type GetCurrentAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentAdminLogic {
	return &GetCurrentAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentAdminLogic) GetCurrentAdmin(req *types.GetCurrentAdminRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetCurrentAdmin(l.ctx, &admin.GetCurrentAdminRequest{
		AdminId: req.AdminId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "获取信息失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"id":        rpcResp.Id,
			"username":  rpcResp.Username,
			"role_id":   rpcResp.RoleId,
			"status":    rpcResp.Status,
			"role_name": rpcResp.RoleName,
		},
	}, nil
}
