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

type CreateAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminLogic) CreateAdmin(req *types.CreateAdminRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.CreateAdmin(l.ctx, &admin.CreateAdminRequest{
		Username: req.Username,
		Password: req.Password,
		RoleId:   req.RoleId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "创建失败: " + err.Error(),
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
