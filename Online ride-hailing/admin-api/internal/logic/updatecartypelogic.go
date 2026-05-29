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

type UpdateCarTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCarTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCarTypeLogic {
	return &UpdateCarTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCarTypeLogic) UpdateCarType(req *types.UpdateCarTypeRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.UpdateCarType(l.ctx, &admin.UpdateCarTypeRequest{
		Id:          req.Id,
		TypeName:    req.TypeName,
		BasePrice:   req.BasePrice,
		KmPrice:     req.KmPrice,
		MinutePrice: req.MinutePrice,
		Status:      req.Status,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "更新车型失败: " + err.Error(),
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
