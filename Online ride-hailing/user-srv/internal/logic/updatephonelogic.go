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

type UpdatePhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhoneLogic {
	return &UpdatePhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePhoneLogic) UpdatePhone(in *user.UpdatePhoneReq) (*user.UpdatePhoneResp, error) {
	if in.UserId == 0 {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

	if in.OldPhone == "" {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "旧手机号不能为空",
		}, nil
	}

	if in.NewPhone == "" {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "新手机号不能为空",
		}, nil
	}

	if in.Code == "" {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "验证码不能为空",
		}, nil
	}

	if in.OldPhone == in.NewPhone {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "新手机号与旧手机号相同",
		}, nil
	}

	var passengerUser user2.PassengerUser
	err := global.DB.Where("id = ?", in.UserId).First(&passengerUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.UpdatePhoneResp{
				Code:    404,
				Message: "用户不存在",
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &user.UpdatePhoneResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	if passengerUser.Phone != in.OldPhone {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "旧手机号与用户当前手机号不匹配",
		}, nil
	}

	key := "sms:code:" + in.OldPhone
	storedCode, err := global.RDB.Get(global.Ctx, key).Result()
	if err != nil {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "验证码已过期或不存在",
		}, nil
	}

	if storedCode != in.Code {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "验证码错误",
		}, nil
	}

	var existingUser user2.PassengerUser
	err = global.DB.Where("phone = ? AND id != ?", in.NewPhone, in.UserId).First(&existingUser).Error
	if err == nil {
		return &user.UpdatePhoneResp{
			Code:    400,
			Message: "新手机号已被其他用户使用",
		}, nil
	}

	err = global.DB.Model(&passengerUser).Update("phone", in.NewPhone).Error
	if err != nil {
		l.Errorf("更新手机号失败: %v", err)
		return &user.UpdatePhoneResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	global.RDB.Del(global.Ctx, key)

	return &user.UpdatePhoneResp{
		Code:    200,
		Message: "手机号修改成功",
	}, nil
}
