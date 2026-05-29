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

type DeleteCarTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCarTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCarTypeLogic {
	return &DeleteCarTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCarTypeLogic) DeleteCarType(req *types.DeleteCarTypeRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.DeleteCarType(l.ctx, &admin.DeleteCarTypeRequest{
		Id: req.Id,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "删除车型失败: " + err.Error(),
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
