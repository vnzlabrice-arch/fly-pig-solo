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

type AdminLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLoginLogic) AdminLogin(req *types.AdminLoginRequest) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.AdminSrv.AdminLogin(l.ctx, &admin.AdminLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "登录失败: " + err.Error(),
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "登录成功",
		Data: map[string]interface{}{
			"token":     rpcResp.Token,
			"admin_id":  rpcResp.AdminId,
			"username":  rpcResp.Username,
		},
	}, nil
}
