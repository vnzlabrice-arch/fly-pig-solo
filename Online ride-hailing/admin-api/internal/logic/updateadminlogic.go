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

type UpdateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminLogic {
	return &UpdateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdminLogic) UpdateAdmin(req *types.UpdateAdminRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.UpdateAdmin(l.ctx, &admin.UpdateAdminRequest{
		AdminId:  req.AdminId,
		Username: req.Username,
		RoleId:   req.RoleId,
		Status:   req.Status,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "更新失败: " + err.Error(),
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
