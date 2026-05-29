package logic

import (
	"context"
	"errors"
	"time"
	"user-srv/global"
	"user-srv/internal/svc"
	user2 "user-srv/model"
	"user-srv/user"

	"user-srv/pkg"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 验证手机号
	if in.Phone == "" {
		return &user.LoginResp{
			Code:    400,
			Message: "手机号不能为空",
		}, nil
	}

	// 验证密码或验证码
	if in.Code == "" {
		return &user.LoginResp{
			Code:    400,
			Message: "验证码不能为空",
		}, nil
	}

	// 查找用户
	var existingUser user2.PassengerUser
	err := global.DB.Where("phone = ?", in.Phone).First(&existingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.LoginResp{
				Code:    400,
				Message: "用户不存在",
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &user.LoginResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 验证验证码
	if in.Code != "" {
		// 从Redis中获取验证码
		key := "sms:code:" + in.Phone
		storedCode, err := global.RDB.Get(global.Ctx, key).Result()
		if err != nil {
			l.Errorf("获取验证码失败: %v", err)
			return &user.LoginResp{
				Code:    400,
				Message: "验证码已过期或不存在",
			}, nil
		}

		// 验证验证码
		if storedCode != in.Code {
			return &user.LoginResp{
				Code:    400,
				Message: "验证码错误",
			}, nil
		}

		// 验证码验证成功后，删除Redis中的验证码
		err = global.RDB.Del(global.Ctx, key).Err()
		if err != nil {
			l.Errorf("删除验证码失败: %v", err)
		}
	}

	// 更新最后登录时间
	now := time.Now()
	existingUser.LastLoginTime = &now
	err = global.DB.Save(&existingUser).Error
	if err != nil {
		l.Errorf("更新登录时间失败: %v", err)
	}

	// 生成JWT令牌
	token, err := pkg.GenerateToken(existingUser.ID, in.Phone)
	if err != nil {
		l.Errorf("生成令牌失败: %v", err)
		return &user.LoginResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	return &user.LoginResp{
		Code:    200,
		Message: "登录成功",
		UserId:  existingUser.ID,
		Token:   token,
	}, nil
}
