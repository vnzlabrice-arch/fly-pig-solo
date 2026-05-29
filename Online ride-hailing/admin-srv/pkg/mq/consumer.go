package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"admin-srv/global"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/zeromicro/go-zero/core/logx"
)

// StartDispatchConsumer 启动派单消息消费者（异步处理派单逻辑）
// 流程图步骤：MQ消费者异步接收消息 → 校验 → 更新订单 → 推送给司机
func StartDispatchConsumer() error {
	if global.Config.RocketMQ.NameServers == nil || len(global.Config.RocketMQ.NameServers) == 0 {
		return fmt.Errorf("RocketMQ NameServers 未配置")
	}

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(global.Config.RocketMQ.NameServers),
		consumer.WithGroupName(global.Config.RocketMQ.GroupName+"_consumer"),
		consumer.WithConsumerModel(consumer.Clustering), // 集群消费模式
	)
	if err != nil {
		return fmt.Errorf("创建消费者失败: %w", err)
	}

	err = c.Subscribe(global.Config.RocketMQ.Topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			if err := handleSingleMessage(ctx, msg); err != nil {
				return consumer.ConsumeRetryLater, err
			}
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return fmt.Errorf("订阅主题失败: %w", err)
	}

	if err = c.Start(); err != nil {
		return fmt.Errorf("启动消费者失败: %w", err)
	}

	fmt.Println("✅ RocketMQ Dispatch Consumer 启动成功，等待接收派单消息...")
	return nil
}

// handleSingleMessage 处理单条派单消息（核心业务逻辑）
func handleSingleMessage(ctx context.Context, msg *primitive.MessageExt) error {
	var dispatchMsg DispatchMessage
	if err := json.Unmarshal(msg.Body, &dispatchMsg); err != nil {
		logx.Errorf("解析派单消息失败: %v, body=%s", err, string(msg.Body))
		return fmt.Errorf("消息格式错误: %w", err)
	}

	logx.Infof("📥 接收到派单消息: order_id=%s, driver_id=%d",
		dispatchMsg.OrderID, dispatchMsg.DriverID)

	// ========== 1. 查询订单是否存在 ==========
	var existingOrder struct {
		OrderID string `gorm:"column:order_id"`
		DriverID int64  `gorm:"column:driver_id"`
		Status   int8   `gorm:"column:status"`
	}

	result := global.DB.Table("passenger_orders").
		Select("order_id, driver_id, status").
		Where("order_id = ?", dispatchMsg.OrderID).
		First(&existingOrder)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			logx.Slowf("❌ 订单不存在: order_id=%s", dispatchMsg.OrderID)
			return nil // 丢弃消息
		}
		logx.Errorf("查询订单失败: %v, order_id=%s", result.Error, dispatchMsg.OrderID)
		return result.Error // 重试
	}

	// ========== 2. 校验订单状态是否是待派单 (status=0) ==========
	if existingOrder.Status != 0 { // OrderStatusPendingDispatch = 0
		logx.Slowf("⚠️  订单状态不允许派单(当前status=%d, 仅允许status=0待派单): order_id=%s",
			existingOrder.Status, dispatchMsg.OrderID)
		return nil // 丢弃消息，派单失败
	}

	// ========== 3. 检查是否已派过单 ==========
	if existingOrder.DriverID > 0 {
		logx.Slowf("⚠️  该订单已派单给其他司机(driver_id=%d): order_id=%s",
			existingOrder.DriverID, dispatchMsg.OrderID)
		return nil // 丢弃消息，避免重复派单
	}

	// ========== 4. 更新订单：设置 driver_id 和 dispatch_time ==========
	now := time.Now()
	updateResult := global.DB.Table("passenger_orders").
		Where("order_id = ? AND status = ? AND driver_id IS NULL",
			dispatchMsg.OrderID, 0). // 确保是待派单且未被分配
		Updates(map[string]interface{}{
			"driver_id":      dispatchMsg.DriverID,
			"dispatch_time":  now,
			"status":         1, // 更新为：待接单(OrderStatusPendingAccept=1)
		})

	if updateResult.Error != nil {
		logx.Errorf("更新订单失败: %v, order_id=%s", updateResult.Error, dispatchMsg.OrderID)
		return updateResult.Error // 重试
	}

	if updateResult.RowsAffected == 0 {
		logx.Slowf("⚠️  更新失败（可能已被其他请求处理）: order_id=%s", dispatchMsg.OrderID)
		return nil
	}

	// ========== 5. 推送派单消息给司机端 ==========
	// TODO: 调用司机端WebSocket或推送服务通知司机有新订单
	logx.Infof("📤 推送派单消息给司机端: order_id=%s, driver_id=%d",
		dispatchMsg.OrderID, dispatchMsg.DriverID)

	// ========== 6. 派单完成 ==========
	logx.Infof("✅ 后台派单处理完成: order_id=%s → driver_id=%d, dispatch_time=%s",
		dispatchMsg.OrderID, dispatchMsg.DriverID, now.Format("2006-01-02 15:04:05"))

	return nil
}
