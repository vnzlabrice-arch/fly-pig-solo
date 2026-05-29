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

type AddAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAddressLogic) AddAddress(req *types.AddAddressReq) (resp *types.Request, err error) {
	rpcResp, err := l.svcCtx.UserRpc.AddAddress(l.ctx, &usermodel.AddAddressReq{
		UserId:    req.UserId,
		Tag:       req.Tag,
		Address:   req.Address,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		return &types.Request{
			Code:    500,
			Message: "添加地址失败",
			Data:    nil,
		}, nil
	}

	return &types.Request{
		Code:    200,
		Message: "添加地址成功",
		Data:    rpcResp,
	}, nil
}
