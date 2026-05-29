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

type DeleteAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdminLogic {
	return &DeleteAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAdminLogic) DeleteAdmin(req *types.DeleteAdminRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.DeleteAdmin(l.ctx, &admin.DeleteAdminRequest{
		AdminId: req.AdminId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "删除失败: " + err.Error(),
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
