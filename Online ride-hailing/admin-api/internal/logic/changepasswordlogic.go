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

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.ChangePassword(l.ctx, &admin.ChangePasswordRequest{
		AdminId:     req.AdminId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "修改密码失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: rpcResp.Message,
		Data: map[string]interface{}{
			"success": rpcResp.Success,
		},
	}, nil
}
