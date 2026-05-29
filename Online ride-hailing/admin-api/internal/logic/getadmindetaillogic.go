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

type GetAdminDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminDetailLogic {
	return &GetAdminDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminDetailLogic) GetAdminDetail(req *types.GetAdminDetailRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetAdminDetail(l.ctx, &admin.GetAdminDetailRequest{
		AdminId: req.AdminId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "查询详情失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "查询成功",
		Data: map[string]interface{}{
			"id":              rpcResp.Id,
			"username":        rpcResp.Username,
			"role_id":         rpcResp.RoleId,
			"role_name":       rpcResp.RoleName,
			"status":          rpcResp.Status,
			"last_login_time": rpcResp.LastLoginTime,
		},
	}, nil
}
