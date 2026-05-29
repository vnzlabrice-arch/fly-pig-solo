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

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.CreateRole(l.ctx, &admin.CreateRoleRequest{
		Name:    req.Name,
		Remark:  req.Remark,
		MenuIds: req.MenuIds,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "创建角色失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "创建成功",
		Data: map[string]interface{}{
			"id":      rpcResp.Id,
			"success": rpcResp.Success,
		},
	}, nil
}
