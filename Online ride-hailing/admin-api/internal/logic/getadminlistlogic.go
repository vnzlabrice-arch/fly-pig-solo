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

type GetAdminListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminListLogic {
	return &GetAdminListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminListLogic) GetAdminList(req *types.GetAdminListRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetAdminList(l.ctx, &admin.GetAdminListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "查询失败: " + err.Error(),
		}, nil
	}

	var list []interface{}
	for _, item := range rpcResp.List {
		list = append(list, map[string]interface{}{
			"id":              item.Id,
			"username":        item.Username,
			"role_id":         item.RoleId,
			"role_name":       item.RoleName,
			"status":          item.Status,
			"last_login_time": item.LastLoginTime,
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
