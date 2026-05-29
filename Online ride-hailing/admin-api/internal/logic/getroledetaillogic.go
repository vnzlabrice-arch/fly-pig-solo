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

type GetRoleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleDetailLogic {
	return &GetRoleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleDetailLogic) GetRoleDetail(req *types.GetRoleDetailRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetRoleDetail(l.ctx, &admin.GetRoleDetailRequest{
		RoleId: req.RoleId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "查询角色详情失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "查询成功",
		Data: map[string]interface{}{
			"id":        rpcResp.Id,
			"name":      rpcResp.Name,
			"remark":    rpcResp.Remark,
			"menu_ids":  rpcResp.MenuIds,
		},
	}, nil
}
