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

type GetUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserDetailLogic) GetUserDetail(req *types.UserDetailReq) (resp *types.Request, err error) {

	data, err := l.svcCtx.UserRpc.GetUserDetail(l.ctx, &usermodel.UserDetailReq{
		UserId: req.UserId,
	})
	if err != nil {
		logx.Error("RPC调用失败: ", err)
		return &types.Request{
			Code:    500,
			Message: "用户查询详情失败",
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
			"user_id":         data.UserId,
			"phone":           data.Phone,
			"nickname":        data.Nickname,
			"avatar":          data.Avatar,
			"user_type":       data.UserType,
			"create_time":     data.CreateTime,
			"last_login_time": data.LastLoginTime,
		},
	}, nil
}
