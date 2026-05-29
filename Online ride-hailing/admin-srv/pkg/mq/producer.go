package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"admin-srv/global"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var dispatchProducer rocketmq.Producer

// DispatchMessage 派单消息结构
type DispatchMessage struct {
	OrderID  string `json:"order_id"`  // 订单号
	DriverID int64  `json:"driver_id"` // 司机ID
}

// InitProducer 初始化 RocketMQ 生产者（在程序启动时调用）
func InitProducer() error {
	if global.Config.RocketMQ.NameServers == nil || len(global.Config.RocketMQ.NameServers) == 0 {
		return fmt.Errorf("RocketMQ NameServers 未配置")
	}

	p, err := rocketmq.NewProducer(
		producer.WithNameServer(global.Config.RocketMQ.NameServers),
		producer.WithGroupName(global.Config.RocketMQ.GroupName),
		producer.WithRetry(2), // 失败重试2次
	)
	if err != nil {
		return fmt.Errorf("创建生产者失败: %w", err)
	}

	if err = p.Start(); err != nil {
		return fmt.Errorf("启动生产者失败: %w", err)
	}

	dispatchProducer = p
	fmt.Println("✅ RocketMQ Producer 启动成功")
	return nil
}

// SendDispatchMessage 发送派单消息到 MQ
// 流程图步骤：派单API参数校验 → 发送派单消息到 MQ 消息队列
func SendDispatchMessage(orderID string, driverID int64) error {
	if dispatchProducer == nil {
		return fmt.Errorf("RocketMQ Producer 未初始化")
	}

	msg := DispatchMessage{
		OrderID:  orderID,
		DriverID: driverID,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	m := &primitive.Message{
		Topic: global.Config.RocketMQ.Topic,
		Body:  body,
	}
	m.WithTag("dispatch_order") // 设置标签

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := dispatchProducer.SendSync(ctx, m)
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}

	fmt.Printf("📤 派单消息发送成功: order_id=%s, driver_id=%d, msg_id=%s\n",
		orderID, driverID, result.MsgID)
	return nil
}

// CloseProducer 关闭生产者（程序退出时调用）
func CloseProducer() error {
	if dispatchProducer != nil {
		return dispatchProducer.Shutdown()
	}
	return nil
}
