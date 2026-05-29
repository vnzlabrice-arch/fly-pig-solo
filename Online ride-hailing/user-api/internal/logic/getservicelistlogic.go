package logic

import (
	"context"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServiceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetServiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServiceListLogic {
	return &GetServiceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetServiceListLogic) GetServiceList(req *types.ServiceListReq) (resp *types.Request, err error) {
	// 暂未实现服务列表查询功能，因为 user-srv 中未定义此 RPC 方法
	return &types.Request{
		Code:    200,
		Message: "服务列表查询功能待实现",
		Data:    nil,
	}, nil
}
