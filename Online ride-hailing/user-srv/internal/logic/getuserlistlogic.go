package logic

import (
	"context"
	"user-srv/global"
	"user-srv/internal/svc"
	user2 "user-srv/model"
	"user-srv/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(in *user.UserListReq) (*user.UserListResp, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	if in.PageSize > 100 {
		in.PageSize = 100
	}

	offset := (in.Page - 1) * in.PageSize

	var passengerUsers []user2.PassengerUser
	var total int64

	query := global.DB.Model(&user2.PassengerUser{})

	err := query.Count(&total).Error
	if err != nil {
		l.Errorf("查询用户总数失败: %v", err)
		return &user.UserListResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	err = query.Offset(int(offset)).Limit(int(in.PageSize)).Find(&passengerUsers).Error
	if err != nil {
		l.Errorf("查询用户列表失败: %v", err)
		return &user.UserListResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	users := make([]*user.UserInfo, 0, len(passengerUsers))
	for _, pu := range passengerUsers {
		users = append(users, &user.UserInfo{
			UserId:     pu.ID,
			Phone:      pu.Phone,
			Nickname:   pu.Nickname,
			Avatar:     pu.AvatarURL,
			UserType:   1,
			CreateTime: pu.CreatedAt.Unix(),
		})
	}

	return &user.UserListResp{
		Code:     200,
		Message:  "获取用户列表成功",
		Users:    users,
		Total:    int32(total),
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
