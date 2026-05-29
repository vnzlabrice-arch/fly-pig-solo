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

type GetMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuListLogic {
	return &GetMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuListLogic) GetMenuList(req *types.GetMenuListRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetMenuList(l.ctx, &admin.GetMenuListRequest{})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "查询菜单列表失败: " + err.Error(),
		}, nil
	}

	var list []interface{}
	for _, item := range rpcResp.List {
		list = append(list, map[string]interface{}{
			"id":        item.Id,
			"parent_id": item.ParentId,
			"name":      item.Name,
			"path":      item.Path,
			"icon":      item.Icon,
			"sort":      item.Sort,
		})
	}

	return &types.Request{
		Code:    200,
		Message: "查询成功",
		Data: map[string]interface{}{
			"list": list,
		},
	}, nil
}
