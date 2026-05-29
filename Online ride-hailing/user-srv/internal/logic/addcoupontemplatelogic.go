package logic

import (
	"context"
	"time"
	"user-srv/model"
	"user-srv/user"

	"user-srv/global"
	"user-srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCouponTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCouponTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCouponTemplateLogic {
	return &AddCouponTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCouponTemplateLogic) AddCouponTemplate(in *user.AddCouponTemplateReq) (*user.AddCouponTemplateResp, error) {
	if in.Name == "" {
		return &user.AddCouponTemplateResp{
			Code:    400,
			Message: "优惠券名称不能为空",
		}, nil
	}

	if in.Type <= 0 || in.Type > 4 {
		return &user.AddCouponTemplateResp{
			Code:    400,
			Message: "优惠券类型无效",
		}, nil
	}

	if in.ValidType <= 0 || in.ValidType > 2 {
		return &user.AddCouponTemplateResp{
			Code:    400,
			Message: "有效期类型无效",
		}, nil
	}

	if in.ValidType == 1 {
		if in.ValidStart <= 0 || in.ValidEnd <= 0 {
			return &user.AddCouponTemplateResp{
				Code:    400,
				Message: "固定有效期需要设置开始和结束时间",
			}, nil
		}
		if in.ValidEnd <= in.ValidStart {
			return &user.AddCouponTemplateResp{
				Code:    400,
				Message: "结束时间必须大于开始时间",
			}, nil
		}
	} else {
		if in.ValidDays <= 0 {
			return &user.AddCouponTemplateResp{
				Code:    400,
				Message: "领取后有效天数必须大于0",
			}, nil
		}
	}

	template := model.CouponTemplate{
		Name:         in.Name,
		Type:         int8(in.Type),
		Discount:     in.Discount,
		MinAmount:    in.MinAmount,
		ReduceAmount: in.ReduceAmount,
		MaxReduce:    in.MaxReduce,
		ValidType:    int8(in.ValidType),
		Total:        int(in.Total),
		PerLimit:     int(in.PerLimit),
		CityCode:     in.CityCode,
		StartRegion:  in.StartRegion,
		EndRegion:    in.EndRegion,
		CarType:      in.CarType,
		UseTime:      in.UseTime,
		Status:       1,
	}

	if in.ValidType == 1 {
		template.ValidStart = time.Unix(in.ValidStart, 0)
		template.ValidEnd = time.Unix(in.ValidEnd, 0)
	} else {
		template.ValidDays = int(in.ValidDays)
	}

	err := global.DB.Create(&template).Error
	if err != nil {
		l.Errorf("创建优惠券模板失败: %v", err)
		return &user.AddCouponTemplateResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	l.Infof("成功创建优惠券模板: %s, ID: %d", in.Name, template.ID)

	return &user.AddCouponTemplateResp{
		Code:       200,
		Message:    "创建优惠券模板成功",
		TemplateId: template.ID,
	}, nil
}
