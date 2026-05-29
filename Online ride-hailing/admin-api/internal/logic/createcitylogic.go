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

type CreateCityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCityLogic {
	return &CreateCityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCityLogic) CreateCity(req *types.CreateCityRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.CreateCity(l.ctx, &admin.CreateCityRequest{
		CityCode: req.CityCode,
		CityName: req.CityName,
		Status:   req.Status,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "创建城市失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "创建成功",
		Data: map[string]interface{}{
			"id": rpcResp.Id,
		},
	}, nil
}
