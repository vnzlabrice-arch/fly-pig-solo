package logic

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user-srv/model"
	"user-srv/user"

	"user-srv/global"
	"user-srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCouponLogic {
	return &AddCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCouponLogic) AddCoupon(in *user.AddCouponReq) (*user.AddCouponResp, error) {
	if in.UserId == 0 {
		return &user.AddCouponResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

	if in.TemplateId == 0 {
		return &user.AddCouponResp{
			Code:    400,
			Message: "优惠券模板ID不能为空",
		}, nil
	}

	count := in.Count
	if count <= 0 {
		count = 1
	}

	var template model.CouponTemplate
	err := global.DB.Where("id = ? AND status = 1", in.TemplateId).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.AddCouponResp{
				Code:    404,
				Message: "优惠券模板不存在或已下架",
			}, nil
		}
		l.Errorf("查询优惠券模板失败: %v", err)
		return &user.AddCouponResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	if template.Total > 0 && template.Received >= template.Total {
		return &user.AddCouponResp{
			Code:    400,
			Message: "该优惠券已领完",
		}, nil
	}

	var userCouponCount int64
	err = global.DB.Model(&model.UserCoupon{}).Where("user_id = ? AND template_id = ?", in.UserId, in.TemplateId).Count(&userCouponCount).Error
	if err != nil {
		l.Errorf("查询用户已领取优惠券数量失败: %v", err)
		return &user.AddCouponResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	if template.PerLimit > 0 && userCouponCount >= int64(template.PerLimit) {
		return &user.AddCouponResp{
			Code:    400,
			Message: fmt.Sprintf("每人限领%d张，您已领取%d张", template.PerLimit, userCouponCount),
		}, nil
	}

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var successCount int32
	for i := int32(0); i < count; i++ {
		if template.Total > 0 && template.Received >= template.Total {
			break
		}

		couponNo := generateCouponNo()

		var startTime, endTime time.Time
		if template.ValidType == 1 {
			startTime = template.ValidStart
			endTime = template.ValidEnd
		} else {
			startTime = time.Now()
			endTime = startTime.Add(time.Duration(template.ValidDays) * 24 * time.Hour)
		}

		userCoupon := model.UserCoupon{
			UserID:     in.UserId,
			TemplateID: in.TemplateId,
			CouponNo:   couponNo,
			Status:     1,
			StartTime:  startTime,
			EndTime:    endTime,
		}

		err := tx.Create(&userCoupon).Error
		if err != nil {
			tx.Rollback()
			l.Errorf("创建用户优惠券失败: %v", err)
			return &user.AddCouponResp{
				Code:    500,
				Message: "系统错误",
			}, nil
		}

		err = tx.Model(&model.CouponTemplate{}).Where("id = ?", in.TemplateId).Update("received", gorm.Expr("received + 1")).Error
		if err != nil {
			tx.Rollback()
			l.Errorf("更新优惠券领取数量失败: %v", err)
			return &user.AddCouponResp{
				Code:    500,
				Message: "系统错误",
			}, nil
		}

		template.Received++
		successCount++
	}

	err = tx.Commit().Error
	if err != nil {
		l.Errorf("提交事务失败: %v", err)
		return &user.AddCouponResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	l.Infof("成功为用户%d发放%d张优惠券，模板ID: %d", in.UserId, successCount, in.TemplateId)

	return &user.AddCouponResp{
		Code:         200,
		Message:      "发放优惠券成功",
		SuccessCount: successCount,
	}, nil
}

func generateCouponNo() string {
	return fmt.Sprintf("CP%d%08d", time.Now().Unix(), randInt(10000000, 99999999))
}

func randInt(min, max int) int {
	return min + int(time.Now().UnixNano())%(max-min+1)
}
