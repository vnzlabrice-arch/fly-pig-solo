package logic

import (
	"context"
	"errors"
	"time"
	"user-srv/model"
	"user-srv/user"

	"user-srv/global"
	"user-srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetAvailableCouponsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAvailableCouponsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailableCouponsLogic {
	return &GetAvailableCouponsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAvailableCouponsLogic) GetAvailableCoupons(in *user.GetAvailableCouponsReq) (*user.GetAvailableCouponsResp, error) {
	if in.UserId == 0 {
		return &user.GetAvailableCouponsResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

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
	now := time.Now()

	var templates []model.CouponTemplate
	var total int64

	query := global.DB.Model(&model.CouponTemplate{}).Where("status = ?", 1)

	err := query.Count(&total).Error
	if err != nil {
		l.Errorf("查询可领取优惠券总数失败: %v", err)
		return &user.GetAvailableCouponsResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	err = query.Offset(int(offset)).Limit(int(in.PageSize)).Order("id DESC").Find(&templates).Error
	if err != nil {
		l.Errorf("查询可领取优惠券列表失败: %v", err)
		return &user.GetAvailableCouponsResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	var userCouponCountMap map[int64]int
	userCouponCounts, err := l.getUserCouponCounts(in.UserId)
	if err == nil {
		userCouponCountMap = userCouponCounts
	} else {
		userCouponCountMap = make(map[int64]int)
	}

	availableTemplates := make([]*user.AvailableCouponTemplateInfo, 0, len(templates))

	for _, t := range templates {
		if t.ValidType == 1 {
			if now.After(t.ValidEnd) {
				continue
			}
		}

		if t.Total > 0 && t.Received >= t.Total {
			continue
		}

		userReceived := userCouponCountMap[t.ID]
		if t.PerLimit > 0 && userReceived >= t.PerLimit {
			continue
		}

		remain := 0
		if t.Total > 0 {
			remain = t.Total - t.Received
		} else {
			remain = -1
		}

		validStart := int64(0)
		validEnd := int64(0)
		if !t.ValidStart.IsZero() {
			validStart = t.ValidStart.Unix()
		}
		if !t.ValidEnd.IsZero() {
			validEnd = t.ValidEnd.Unix()
		}

		templateInfo := &user.AvailableCouponTemplateInfo{
			TemplateId:   t.ID,
			Name:         t.Name,
			Type:         int32(t.Type),
			Discount:     t.Discount,
			MinAmount:    t.MinAmount,
			ReduceAmount: t.ReduceAmount,
			MaxReduce:    t.MaxReduce,
			ValidType:    int32(t.ValidType),
			ValidStart:   validStart,
			ValidEnd:     validEnd,
			ValidDays:    int32(t.ValidDays),
			Total:        int32(t.Total),
			Received:     int32(t.Received),
			Remain:       int32(remain),
			PerLimit:     int32(t.PerLimit),
			UserReceived: int32(userReceived),
			CityCode:     t.CityCode,
			StartRegion:  t.StartRegion,
			EndRegion:    t.EndRegion,
			CarType:      t.CarType,
			UseTime:      t.UseTime,
		}

		availableTemplates = append(availableTemplates, templateInfo)
	}

	return &user.GetAvailableCouponsResp{
		Code:      200,
		Message:   "获取可领取优惠券列表成功",
		Templates: availableTemplates,
		Total:     int32(len(availableTemplates)),
		Page:      in.Page,
		PageSize:  in.PageSize,
	}, nil
}

func (l *GetAvailableCouponsLogic) getUserCouponCounts(userId int64) (map[int64]int, error) {
	var results []struct {
		TemplateID int64
		Count      int
	}

	err := global.DB.Model(&model.UserCoupon{}).
		Select("template_id, COUNT(*) as count").
		Where("user_id = ?", userId).
		Group("template_id").
		Find(&results).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	countMap := make(map[int64]int)
	for _, r := range results {
		countMap[r.TemplateID] = r.Count
	}

	return countMap, nil
}
