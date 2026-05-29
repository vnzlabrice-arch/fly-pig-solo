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

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.CreateMenu(l.ctx, &admin.CreateMenuRequest{
		ParentId: req.ParentId,
		Name:     req.Name,
		Path:     req.Path,
		Icon:     req.Icon,
		Sort:     req.Sort,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "创建菜单失败: " + err.Error(),
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
