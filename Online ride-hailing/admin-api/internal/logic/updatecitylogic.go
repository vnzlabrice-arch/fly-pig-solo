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

type UpdateCityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCityLogic {
	return &UpdateCityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCityLogic) UpdateCity(req *types.UpdateCityRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.UpdateCity(l.ctx, &admin.UpdateCityRequest{
		Id:       req.Id,
		CityCode: req.CityCode,
		CityName: req.CityName,
		Status:   req.Status,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "更新城市失败: " + err.Error(),
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
