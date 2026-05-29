package logic

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"
	"user-srv/global"
	"user-srv/internal/svc"
	"user-srv/user"

	"user-srv/pkg"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCodeLogic {
	return &SendCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendCodeLogic) SendCode(in *user.SendCodeReq) (*user.SendCodeResp, error) {
	// 验证手机号
	if in.Phone == "" {
		return &user.SendCodeResp{
			Code:    400,
			Message: "手机号不能为空",
		}, nil
	}

	// 生成随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := strconv.Itoa(rnd.Intn(900000) + 100000) // 生成6位随机验证码

	// 模拟发送短信（实际项目中需要调用真实的短信服务）
	smsResult := pkg.SmsSend(in.Phone, code)
	if smsResult.Code != 2 {
		l.Errorf("发送短信失败: %s", smsResult.Msg)
		return &user.SendCodeResp{
			Code:    500,
			Message: "发送验证码失败",
		}, nil
	}
	l.Infof("模拟发送短信成功，手机号: %s, 验证码: %s", in.Phone, code)

	// 模拟存储验证码到Redis（实际项目中需要使用真实的Redis）
	key := "sms:code:" + in.Phone
	err := global.RDB.Set(global.Ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		l.Errorf("存储验证码失败: %v", err)
		return &user.SendCodeResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	i, err := global.RDB.Incr(global.Ctx, "sms:num:"+in.Phone).Result()

	if i > 0 {
		return nil, errors.New("不能重复发送短信")
	}

	l.Infof("发送验证码成功，手机号: %s, 验证码: %s", in.Phone, code)

	return &user.SendCodeResp{
		Code:    200,
		Message: "发送验证码成功",
	}, nil

}
