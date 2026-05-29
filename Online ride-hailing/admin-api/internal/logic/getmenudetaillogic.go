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

type GetMenuDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuDetailLogic {
	return &GetMenuDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuDetailLogic) GetMenuDetail(req *types.GetMenuDetailRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetMenuDetail(l.ctx, &admin.GetMenuDetailRequest{
		MenuId: req.MenuId,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "查询菜单详情失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "查询成功",
		Data: map[string]interface{}{
			"id":        rpcResp.Id,
			"parent_id": rpcResp.ParentId,
			"name":      rpcResp.Name,
			"path":      rpcResp.Path,
			"icon":      rpcResp.Icon,
			"sort":      rpcResp.Sort,
		},
	}, nil
}
