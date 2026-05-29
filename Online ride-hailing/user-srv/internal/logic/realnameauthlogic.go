package logic

import (
	"context"
	"errors"
	"user-srv/global"
	"user-srv/internal/svc"
	user2 "user-srv/model"
	"user-srv/user"

	"user-srv/pkg"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RealNameAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRealNameAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RealNameAuthLogic {
	return &RealNameAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RealNameAuthLogic) RealNameAuth(in *user.RealNameAuthReq) (*user.RealNameAuthResp, error) {
	if in.UserId == 0 {
		return &user.RealNameAuthResp{
			Code:    400,
			Message: "用户ID不能为空",
			IsAuth:  false,
		}, nil
	}

	if in.RealName == "" {
		return &user.RealNameAuthResp{
			Code:    400,
			Message: "真实姓名不能为空",
			IsAuth:  false,
		}, nil
	}

	if in.IdCard == "" {
		return &user.RealNameAuthResp{
			Code:    400,
			Message: "身份证号不能为空",
			IsAuth:  false,
		}, nil
	}

	var passengerUser user2.PassengerUser
	err := global.DB.Where("id = ?", in.UserId).First(&passengerUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.RealNameAuthResp{
				Code:    404,
				Message: "用户不存在",
				IsAuth:  false,
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &user.RealNameAuthResp{
			Code:    500,
			Message: "系统错误",
			IsAuth:  false,
		}, nil
	}

	if passengerUser.RealName != "" && passengerUser.IDCardHash != "" {
		return &user.RealNameAuthResp{
			Code:    400,
			Message: "用户已完成实名认证",
			IsAuth:  true,
		}, nil
	}

	isValid, err := pkg.ShenFen(in.RealName, in.IdCard)
	if err != nil {
		l.Errorf("身份验证失败: %v", err)
		return &user.RealNameAuthResp{
			Code:    500,
			Message: "身份验证失败，请稍后重试",
			IsAuth:  false,
		}, nil
	}

	if !isValid {
		return &user.RealNameAuthResp{
			Code:    400,
			Message: "身份证信息与姓名不匹配",
			IsAuth:  false,
		}, nil
	}

	err = global.DB.Model(&passengerUser).Updates(map[string]interface{}{
		"real_name":    in.RealName,
		"id_card_hash": in.IdCard, // 实际项目中应该进行哈希处理
	}).Error
	if err != nil {
		l.Errorf("更新实名认证信息失败: %v", err)
		return &user.RealNameAuthResp{
			Code:    500,
			Message: "系统错误",
			IsAuth:  false,
		}, nil
	}

	return &user.RealNameAuthResp{
		Code:    200,
		Message: "实名认证成功",
		IsAuth:  true,
	}, nil
}
