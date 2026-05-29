package logic

import (
	"context"
	"user-srv/model"
	"user-srv/user"

	"user-srv/global"
	"user-srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOrderListLogic {
	return &GetUserOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserOrderList 获取用户订单列表
func (l *GetUserOrderListLogic) GetUserOrderList(in *user.UserOrderListReq) (*user.UserOrderListResp, error) {
	// 参数验证
	if in.UserId == 0 {
		return &user.UserOrderListResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

	// 设置默认分页参数
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 查询乘客信息
	var passenger model.PassengerUser
	err := global.DB.Where("id = ?", in.UserId).First(&passenger).Error
	if err != nil {
		l.Errorf("查询乘客信息失败: %v", err)
	}

	// 查询订单总数
	var total int64
	query := global.DB.Model(&model.PassengerOrder{}).Where("passenger_id = ?", in.UserId)

	// 状态筛选（0表示全部）
	if in.Status > 0 {
		query = query.Where("status = ?", in.Status)
	}

	err = query.Count(&total).Error
	if err != nil {
		l.Errorf("查询订单总数失败: %v", err)
		return &user.UserOrderListResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 查询订单列表
	var orders []model.PassengerOrder
	err = query.Order("book_time DESC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&orders).Error
	if err != nil {
		l.Errorf("查询订单列表失败: %v", err)
		return &user.UserOrderListResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 构造订单列表响应
	var orderList []*user.PassengerOrderInfo
	for _, order := range orders {
		// 处理时间字段
		var bookTime int64 = order.BookTime.Unix()
		var pickupTime int64
		if order.PickupTime != nil {
			pickupTime = order.PickupTime.Unix()
		}
		var startTime int64
		if order.StartTime != nil {
			startTime = order.StartTime.Unix()
		}
		var endTime int64
		if order.EndTime != nil {
			endTime = order.EndTime.Unix()
		}
		var cancelTime int64
		if order.Status == 6 {
			cancelTime = order.UpdatedAt.Unix()
		}

		orderList = append(orderList, &user.PassengerOrderInfo{
			OrderId:        order.OrderID,
			PassengerId:    order.PassengerID,
			PassengerName:  order.PassengerName,
			PassengerPhone: order.PassengerPhone,
			PassengerAvatar: passenger.AvatarURL,
			OrderType:      int32(order.OrderType),
			CarType:        order.CarType,
			Status:         int32(order.Status),
			PayStatus:      int32(order.PayStatus),
			StartAddress:   order.StartAddress,
			StartLng:       order.StartLng,
			StartLat:       order.StartLat,
			EndAddress:     order.EndAddress,
			EndLng:         order.EndLng,
			EndLat:         order.EndLat,
			CouponName:     order.CouponName,
			EstimatedPrice: order.EstimatedPrice,
			ActualPrice:    order.FinalPrice,
			Remark:         order.PassRemark,
			BookTime:       bookTime,
			AcceptTime:     pickupTime,
			StartTime:      startTime,
			FinishTime:     endTime,
			CancelTime:     cancelTime,
			CancelReason:   order.CancelReason,
			CreateTime:     order.CreatedAt.Unix(),
			UpdateTime:     order.UpdatedAt.Unix(),
		})
	}

	return &user.UserOrderListResp{
		Code:     200,
		Message:  "获取成功",
		Orders:   orderList,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
	}, nil
}
