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

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.UpdateRole(l.ctx, &admin.UpdateRoleRequest{
		RoleId:  req.RoleId,
		Name:    req.Name,
		Remark:  req.Remark,
		MenuIds: req.MenuIds,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "更新角色失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "更新成功",
		Data: map[string]interface{}{
			"success": rpcResp.Success,
		},
	}, nil
}
