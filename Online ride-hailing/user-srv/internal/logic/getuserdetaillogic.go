package logic

import (
	"context"
	"errors"
	"user-srv/global"
	"user-srv/internal/svc"
	user2 "user-srv/model"
	"user-srv/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserDetailLogic) GetUserDetail(in *user.UserDetailReq) (*user.UserDetailResp, error) {
	if in.UserId == 0 {
		return &user.UserDetailResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

	var passengerUser user2.PassengerUser
	err := global.DB.Where("id = ?", in.UserId).First(&passengerUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.UserDetailResp{
				Code:    404,
				Message: "用户不存在",
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &user.UserDetailResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	var lastLoginTime int64
	if passengerUser.LastLoginTime != nil {
		lastLoginTime = passengerUser.LastLoginTime.Unix()
	}

	return &user.UserDetailResp{
		Code:          200,
		Message:       "获取用户详情成功",
		UserId:        passengerUser.ID,
		Phone:         passengerUser.Phone,
		Nickname:      passengerUser.Nickname,
		Avatar:        passengerUser.AvatarURL,
		UserType:      1,
		CreateTime:    passengerUser.CreatedAt.Unix(),
		LastLoginTime: lastLoginTime,
	}, nil
}
