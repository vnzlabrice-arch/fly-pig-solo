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

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.UserListReq) (resp *types.Request, err error) {
	data, err := l.svcCtx.UserRpc.GetUserList(l.ctx, &usermodel.UserListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return &types.Request{
			Code:    400,
			Message: "用户查询详情失败",
			Data:    nil,
		}, nil
	}
	return &types.Request{
		Code:    0,
		Message: "用户列表查询成功",
		Data:    data,
	}, nil
}
