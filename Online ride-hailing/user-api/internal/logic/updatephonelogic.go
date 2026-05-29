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

type UpdatePhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhoneLogic {
	return &UpdatePhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePhoneLogic) UpdatePhone(req *types.UpdatePhoneReq) (resp *types.Request, err error) {

	data, err := l.svcCtx.UserRpc.UpdatePhone(l.ctx, &usermodel.UpdatePhoneReq{

		UserId:   req.UserId,
		OldPhone: req.OldPhone,
		NewPhone: req.NewPhone,
		Code:     req.Code,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "手机号修改失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "手机号修改成功",
		Data:    data,
	}, nil
}
