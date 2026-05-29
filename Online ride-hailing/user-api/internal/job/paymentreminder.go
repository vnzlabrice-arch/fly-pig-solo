package job

import (
	"context"
	"sync"
	"time"

	"user-api/internal/svc"
	"user-api/internal/ws"
	"user-srv/pb/usermodel"

	"github.com/zeromicro/go-zero/core/logx"
)

// PaymentReminderJob 支付提醒和超时取消定时任务
type PaymentReminderJob struct {
	svcCtx *svc.ServiceContext
	hub    *ws.Hub

	// 记录已提醒订单，避免重复推送
	remindedOrders map[string]int
	mu             sync.RWMutex
}

type timeoutCancelMessage struct {
	OrderID      string `json:"order_id"`
	CancelReason string `json:"cancel_reason"`
	CancelTime   int64  `json:"cancel_time"`
}

// NewPaymentReminderJob 创建定时任务
func NewPaymentReminderJob(svcCtx *svc.ServiceContext, hub *ws.Hub) *PaymentReminderJob {
	return &PaymentReminderJob{
		svcCtx:         svcCtx,
		hub:            hub,
		remindedOrders: make(map[string]int),
	}
}

// Run 执行支付提醒任务
func (j *PaymentReminderJob) Run() {
	logx.Info("开始执行支付提醒任务")
	defer logx.Info("支付提醒任务执行完成")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := j.svcCtx.UserRpc.GetUnpaidOrders(ctx, &usermodel.GetUnpaidOrdersReq{
		HoursBefore: 24,
	})
	if err != nil {
		logx.Errorf("查询未支付订单失败: %v", err)
		return
	}

	if resp.Code != 200 {
		logx.Errorf("查询未支付订单返回异常 code=%d, message=%s", resp.Code, resp.Message)
		return
	}

	logx.Infof("查询到 %d 个未支付订单", len(resp.Orders))
	for _, order := range resp.Orders {
		j.processOrder(order)
	}
}

// processOrder 处理单个支付提醒订单
func (j *PaymentReminderJob) processOrder(order *usermodel.PassengerOrderInfo) {
	j.mu.RLock()
	remindCount := j.remindedOrders[order.OrderId]
	j.mu.RUnlock()

	if remindCount >= 3 {
		logx.Infof("订单 %s 已提醒 %d 次，不再重复提醒", order.OrderId, remindCount)
		return
	}

	messageData := map[string]interface{}{
		"order_id":      order.OrderId,
		"start_address": order.StartAddress,
		"end_address":   order.EndAddress,
		"actual_price":  order.ActualPrice,
		"finish_time":   order.FinishTime,
		"remind_count":  remindCount + 1,
	}

	j.hub.SendToUser(order.PassengerId, "payment_reminder", messageData)
	logx.Infof("已向用户 %d 发送订单 %s 的支付提醒，第 %d 次", order.PassengerId, order.OrderId, remindCount+1)

	j.mu.Lock()
	j.remindedOrders[order.OrderId] = remindCount + 1
	j.mu.Unlock()
}

// ClearRemindedOrder 清除订单提醒记录
func (j *PaymentReminderJob) ClearRemindedOrder(orderID string) {
	j.mu.Lock()
	delete(j.remindedOrders, orderID)
	j.mu.Unlock()
	logx.Infof("已清除订单 %s 的提醒记录", orderID)
}

// TimeoutCancelOrders 超时自动取消订单并推送通知（待实现）
func (j *PaymentReminderJob) TimeoutCancelOrders() {
	// 暂未实现超时取消订单功能，因为 user-srv 中未定义此 RPC 方法
	logx.Info("TimeoutCancelOrders 功能待实现")
}
