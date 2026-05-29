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

type DeleteCityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCityLogic {
	return &DeleteCityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCityLogic) DeleteCity(req *types.DeleteCityRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.DeleteCity(l.ctx, &admin.DeleteCityRequest{
		Id: req.Id,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "删除城市失败: " + err.Error(),
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
