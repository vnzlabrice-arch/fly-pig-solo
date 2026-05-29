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

type GetCarTypeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCarTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCarTypeListLogic {
	return &GetCarTypeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCarTypeListLogic) GetCarTypeList(req *types.GetCarTypeListRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.GetCarTypeList(l.ctx, &admin.GetCarTypeListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "查询车型列表失败: " + err.Error(),
		}, nil
	}

	var list []interface{}
	for _, item := range rpcResp.List {
		list = append(list, map[string]interface{}{
			"id":           item.Id,
			"type_name":    item.TypeName,
			"base_price":   item.BasePrice,
			"km_price":     item.KmPrice,
			"minute_price": item.MinutePrice,
			"status":       item.Status,
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
