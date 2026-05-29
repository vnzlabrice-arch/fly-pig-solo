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

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.UpdateMenu(l.ctx, &admin.UpdateMenuRequest{
		MenuId:   req.MenuId,
		ParentId: req.ParentId,
		Name:     req.Name,
		Path:     req.Path,
		Icon:     req.Icon,
		Sort:     req.Sort,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "更新菜单失败: " + err.Error(),
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
