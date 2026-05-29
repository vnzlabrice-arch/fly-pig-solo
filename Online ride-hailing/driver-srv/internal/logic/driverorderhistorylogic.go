package logic

import (
	"context"
	"driver-srv/global"
	"driver-srv/internal/svc"
	user3 "driver-srv/model/driver"
	"driver-srv/pb/driver"
	"driver-srv/pkg"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DriverOrderHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDriverOrderHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriverOrderHistoryLogic {
	return &DriverOrderHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取历史订单列表
func (l *DriverOrderHistoryLogic) DriverOrderHistory(in *driver.DriverOrderHistoryReq) (*driver.DriverOrderHistoryResp, error) {
	query := global.DB.Model(&user3.PassengerOrder{}).Where("driver_id = ?", in.DriverId)

	if in.Status != 0 {
		query = query.Where("status = ?", in.Status)
	} else {
		query = query.Where("status = 5 OR status = 6")
	}

	if in.PayStatus != 0 {
		query = query.Where("pay_status = ?", in.PayStatus)
	}

	if strings.TrimSpace(in.SearchKeyword) != "" {
		key := "%" + strings.TrimSpace(in.SearchKeyword) + "%"
		query = query.Where("order_id LIKE ? OR passenger_name LIKE ? OR passenger_phone LIKE ? OR start_address LIKE ? OR end_address LIKE ?",
			key, key, key, key, key)
	}

	if strings.TrimSpace(in.StartDate) != "" {
		startTime, err := time.Parse("2006-01-02", in.StartDate)
		if err == nil {
			query = query.Where("book_time >= ?", startTime)
		}
	}

	if strings.TrimSpace(in.EndDate) != "" {
		endTime, err := time.Parse("2006-01-02", in.EndDate)
		if err == nil {
			query = query.Where("book_time < ?", endTime.Add(24*time.Hour))
		}
	}

	if in.MinIncome > 0 || in.MaxIncome > 0 {
		query = query.Joins("left join driver_incomes on driver_incomes.order_id = passenger_orders.order_id and driver_incomes.driver_id = passenger_orders.driver_id")
		if in.MinIncome > 0 {
			query = query.Where("driver_incomes.actual_income >= ?", in.MinIncome)
		}
		if in.MaxIncome > 0 {
			query = query.Where("driver_incomes.actual_income <= ?", in.MaxIncome)
		}
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, errors.New("查询失败")
	}

	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 10
	}

	var passengerOrders []user3.PassengerOrder
	err = query.Order("book_time desc").Scopes(pkg.Paginate(int(in.Page), int(in.Size))).Find(&passengerOrders).Error
	if err != nil {
		return nil, errors.New("查询失败")
	}

	serviceScore := float32(0)
	// ServiceScore 字段当前未在 DriverUser 模型中启用

	var orderList []*driver.DriverOrderHistoryItem
	for _, order := range passengerOrders {
		bookTime := ""
		if !order.BookTime.IsZero() {
			bookTime = order.BookTime.Format("2006-01-02T15:04:05Z")
		}

		endTime := ""
		if order.EndTime != nil && !order.EndTime.IsZero() {
			endTime = order.EndTime.Format("2006-01-02T15:04:05Z")
		}

		var income user3.DriverIncome
		actualIncome := float32(0)
		err = global.DB.Where("driver_id = ? AND order_id = ?", in.DriverId, order.OrderID).First(&income).Error
		if err == nil {
			actualIncome = float32(income.ActualIncome)
		}

		orderList = append(orderList, &driver.DriverOrderHistoryItem{
			OrderId:      order.OrderID,
			OrderType:    int64(order.OrderType),
			CarType:      order.CarType,
			Status:       int64(order.Status),
			StartAddress: order.StartAddress,
			EndAddress:   order.EndAddress,
			FinalPrice:   float32(order.FinalPrice),
			ActualIncome: actualIncome,
			BookTime:     bookTime,
			EndTime:      endTime,
			PayStatus:    int64(order.PayStatus),
			ServiceScore: serviceScore,
		})
	}

	return &driver.DriverOrderHistoryResp{
		DriverOrderHistoryItem: orderList,
		Total:                  total,
	}, nil
}
