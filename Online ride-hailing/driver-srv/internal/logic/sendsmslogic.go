package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	"driver-srv/pb/driver"
	"errors"
	"math/rand"
	"strconv"
	"time"
	"user-srv/pkg"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsLogic {
	return &SendSmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendSmsLogic) SendSms(in *driver.SendSmsReq) (*driver.SendSmsResp, error) {
	if in.Phone == "" {
		return nil, errors.New("手机号不能为空")
	}
	//if global.RDB == nil {
	//	return nil, errors.New("redis初始化失败")
	//}
	i, _ := global.RDB.Get(global.Ctx, "sendNum"+in.Phone).Int()
	if i >= 3 {
		return nil, errors.New("一分钟只允许发送3次")
	}

	code := rand.Intn(900000) + 100000

	send := pkg.SmsSend(in.Phone, strconv.Itoa(code))
	if send.Code != 2 {
		return nil, errors.New("验证码发送失败")
	}

	global.RDB.Incr(global.Ctx, "sendNum"+in.Phone)
	global.RDB.Expire(global.Ctx, "sendNum"+in.Phone, time.Minute*1)
	global.RDB.Set(global.Ctx, "send"+in.Phone, code, time.Minute*5)

	return &driver.SendSmsResp{
		Success: true,
	}, nil
}
