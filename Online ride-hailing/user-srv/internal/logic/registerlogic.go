package logic

import (
	"context"
	"errors"
	"regexp"
	"user-srv/global"
	"user-srv/internal/svc"
	user2 "user-srv/model"
	"user-srv/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 验证手机号
	if in.Phone == "" {
		return &user.RegisterResp{
			Code:    400,
			Message: "手机号不能为空",
		}, nil
	}

	phoneReg := regexp.MustCompile(`^1[3-9]\d{9}&`)

	matchString := phoneReg.MatchString(in.Phone)

	if !matchString {
		return nil, errors.New("手机号格式错误")
	}
	// 验证密码

	// 验证验证码（这里简化处理，实际应该从Redis中获取并验证）
	if in.Code == "" {
		return &user.RegisterResp{
			Code:    400,
			Message: "验证码不能为空",
		}, nil
	}

	// 检查手机号是否已注册
	var existingUser user2.PassengerUser
	err := global.DB.Where("phone = ?", in.Phone).First(&existingUser).Error
	if err == nil {
		return &user.RegisterResp{
			Code:    400,
			Message: "手机号已注册",
		}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Errorf("查询用户失败: %v", err)
		return &user.RegisterResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 创建用户
	newUser := user2.PassengerUser{
		Phone:     in.Phone,
		Nickname:  in.Nickname,
		AvatarURL: in.Avatar,
		Status:    1,
	}

	err = global.DB.Create(&newUser).Error
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return &user.RegisterResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	return &user.RegisterResp{
		Code:    200,
		Message: "注册成功",
		UserId:  newUser.ID,
	}, nil
}
