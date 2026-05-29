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

type GetRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleListLogic) GetRoleList(req *types.GetRoleListRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetRoleList(l.ctx, &admin.GetRoleListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "查询角色列表失败: " + err.Error(),
		}, nil
	}

	var list []interface{}
	for _, item := range rpcResp.List {
		list = append(list, map[string]interface{}{
			"id":       item.Id,
			"name":     item.Name,
			"remark":   item.Remark,
			"menu_ids": item.MenuIds,
		})
	}

	return &types.Request{
		Code:    200,
		Message: "查询成功",
		Data: map[string]interface{}{
			"total": rpcResp.Total,
			"list":  list,
		},
	}, nil
}
