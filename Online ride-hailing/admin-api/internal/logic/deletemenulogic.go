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

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteMenuRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.DeleteMenu(l.ctx, &admin.DeleteMenuRequest{
		MenuId: req.MenuId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "删除菜单失败: " + err.Error(),
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
