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

type CreateCarTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCarTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCarTypeLogic {
	return &CreateCarTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCarTypeLogic) CreateCarType(req *types.CreateCarTypeRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.CreateCarType(l.ctx, &admin.CreateCarTypeRequest{
		TypeName:    req.TypeName,
		BasePrice:   req.BasePrice,
		KmPrice:     req.KmPrice,
		MinutePrice: req.MinutePrice,
		Status:      req.Status,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "创建车型失败: " + err.Error(),
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
