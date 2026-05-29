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

type GetUserCouponListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCouponListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCouponListLogic {
	return &GetUserCouponListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserCouponListLogic) GetUserCouponList(in *user.UserCouponListReq) (*user.UserCouponListResp, error) {
	if in.UserId == 0 {
		return &user.UserCouponListResp{
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

	var userCoupons []model.UserCoupon
	var total int64

	query := global.DB.Model(&model.UserCoupon{}).Where("user_id = ?", in.UserId)

	if in.Status > 0 {
		query = query.Where("status = ?", in.Status)
	}

	err := query.Count(&total).Error
	if err != nil {
		l.Errorf("查询用户优惠券总数失败: %v", err)
		return &user.UserCouponListResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	err = query.Offset(int(offset)).Limit(int(in.PageSize)).Order("created_at DESC").Find(&userCoupons).Error
	if err != nil {
		l.Errorf("查询用户优惠券列表失败: %v", err)
		return &user.UserCouponListResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	now := time.Now()

	coupons := make([]*user.UserCouponInfo, 0, len(userCoupons))
	for _, uc := range userCoupons {
		status := uc.Status
		if status == 1 && uc.EndTime.Before(now) {
			status = 3
		}

		coupons = append(coupons, &user.UserCouponInfo{
			Id:         uc.ID,
			UserId:     uc.UserID,
			TemplateId: uc.TemplateID,
			CouponNo:   uc.CouponNo,
			Status:     int32(status),
			StartTime:  uc.StartTime.Unix(),
			EndTime:    uc.EndTime.Unix(),
			CreateTime: uc.CreatedAt.Unix(),
		})
	}

	return &user.UserCouponListResp{
		Code:     200,
		Message:  "获取用户优惠券列表成功",
		Coupons:  coupons,
		Total:    int32(total),
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
