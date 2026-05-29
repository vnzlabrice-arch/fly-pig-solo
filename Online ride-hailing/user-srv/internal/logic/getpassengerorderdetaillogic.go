package logic

import (
	"context"
	"errors"
	"user-srv/model"
	"user-srv/user"

	"user-srv/global"
	"user-srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetPassengerOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPassengerOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPassengerOrderDetailLogic {
	return &GetPassengerOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetPassengerOrderDetail 获取订单详情
func (l *GetPassengerOrderDetailLogic) GetPassengerOrderDetail(in *user.PassengerOrderDetailReq) (*user.PassengerOrderDetailResp, error) {
	// 参数验证
	if in.UserId == 0 {
		return &user.PassengerOrderDetailResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

	if in.OrderId == "" {
		return &user.PassengerOrderDetailResp{
			Code:    400,
			Message: "订单号不能为空",
		}, nil
	}

	// 查询订单
	var order model.PassengerOrder
	err := global.DB.Where("order_id = ? AND passenger_id = ?", in.OrderId, in.UserId).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.PassengerOrderDetailResp{
				Code:    404,
				Message: "订单不存在或不属于当前用户",
			}, nil
		}

		l.Errorf("查询订单详情失败: %v", err)
		return &user.PassengerOrderDetailResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	// 查询乘客信息
	var passenger model.PassengerUser
	err = global.DB.Where("id = ?", in.UserId).First(&passenger).Error
	if err != nil {
		l.Errorf("查询乘客信息失败: %v", err)
	}

	// 查询费用明细
	var feeDetails []model.OrderFeeDetail
	err = global.DB.Where("order_id = ?", in.OrderId).Find(&feeDetails).Error
	if err != nil {
		l.Errorf("查询费用明细失败: %v", err)
	}

	// 构造费用明细响应
	var feeList []*user.OrderFeeDetail
	feeTypeMap := map[string]int32{
		"起步价": 1,
		"里程费": 2,
		"时长费": 3,
		"远途费": 4,
		"夜间费": 5,
		"其他":  6,
	}
	for _, fee := range feeDetails {
		feeType, ok := feeTypeMap[fee.FeeType]
		if !ok {
			feeType = 6
		}
		feeList = append(feeList, &user.OrderFeeDetail{
			FeeType:     feeType,
			FeeName:     fee.FeeType,
			Amount:      fee.Amount,
			Description: fee.Description,
		})
	}

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

	// 构造响应
	return &user.PassengerOrderDetailResp{
		Code:    200,
		Message: "获取成功",
		Order: &user.PassengerOrderInfo{
			OrderId:         order.OrderID,
			PassengerId:     order.PassengerID,
			PassengerName:   order.PassengerName,
			PassengerPhone:  order.PassengerPhone,
			PassengerAvatar: passenger.AvatarURL,
			OrderType:       int32(order.OrderType),
			CarType:         order.CarType,
			Status:          int32(order.Status),
			PayStatus:       int32(order.PayStatus),
			StartAddress:    order.StartAddress,
			StartLng:        order.StartLng,
			StartLat:        order.StartLat,
			EndAddress:      order.EndAddress,
			EndLng:          order.EndLng,
			EndLat:          order.EndLat,
			CouponName:      order.CouponName,
			EstimatedPrice:  order.EstimatedPrice,
			ActualPrice:     order.FinalPrice,
			DriverId:        order.DriverID,
			DriverName:      "司佳帅",
			DriverPhone:     "54838945594",
			LicensePlate:    "黑A5201314",
			Remark:          order.PassRemark,
			BookTime:        bookTime,
			AcceptTime:      pickupTime,
			StartTime:       startTime,
			FinishTime:      endTime,
			CancelTime:      cancelTime,
			CancelReason:    order.CancelReason,
			CreateTime:      order.CreatedAt.Unix(),
			UpdateTime:      order.UpdatedAt.Unix(),
		},
		FeeDetails: feeList,
	}, nil
}
