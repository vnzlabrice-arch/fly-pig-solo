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

type DeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.DeleteRoleRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.DeleteRole(l.ctx, &admin.DeleteRoleRequest{
		RoleId: req.RoleId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "删除角色失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "删除成功",
		Data: map[string]interface{}{
			"success": rpcResp.Success,
		},
	}, nil
}
