// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"user-api/internal/svc"
	"user-api/internal/types"
	"user-srv/pb/usermodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.Request, err error) {
	if req.Phone == "" || req.Code == "" {
		return &types.Request{
			Code:    400,
			Message: "手机号和验证码不能为空",
			Data:    nil,
		}, nil
	}

	data, err := l.svcCtx.UserRpc.Login(l.ctx, &usermodel.LoginReq{
		Phone: req.Phone,
		Code:  req.Code,
	})
	if err != nil {
		logx.Error("RPC调用失败: ", err)
		return &types.Request{
			Code:    500,
			Message: "登录失败",
			Data:    nil,
		}, nil
	}

	if data.Code != 200 {
		return &types.Request{
			Code:    data.Code,
			Message: data.Message,
			Data:    nil,
		}, nil
	}

	return &types.Request{
		Code:    0,
		Message: "success",
		Data: map[string]interface{}{
			"user_id": data.UserId,
			"token":   data.Token,
		},
	}, nil
}
