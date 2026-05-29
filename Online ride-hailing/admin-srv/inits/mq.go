package inits

import (
	"admin-srv/internal/logic"
	"admin-srv/pkg/mq"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"admin-srv/global"
)

// InitMQ 初始化 RocketMQ（生产者+消费者+定时检查器）
func InitMQ() {
	// 检查配置
	if global.Config.RocketMQ.NameServers == nil || len(global.Config.RocketMQ.NameServers) == 0 {
		fmt.Println("⚠️  RocketMQ 未配置，跳过消息队列初始化")
		return
	}

	// 启动生产者
	if err := mq.InitProducer(); err != nil {
		fmt.Printf("❌ RocketMQ Producer 初始化失败: %v\n", err)
		fmt.Println("⚠️  程序继续运行，但派单功能将不可用")
		return
	}

	// 启动消费者（异步处理派单）
	if err := mq.StartDispatchConsumer(); err != nil {
		fmt.Printf("❌ RocketMQ Consumer 初始化失败: %v\n", err)
	}

	// 启动待派单定时检查器（每秒查询 status=0 的订单）
	logic.StartPendingOrderChecker()

	// 注册优雅关闭信号处理
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		fmt.Println("\n🛑 正在关闭服务...")
		logic.Stop() // 停止定时检查器
		if err := mq.CloseProducer(); err != nil {
			fmt.Printf("关闭 Producer 失败: %v\n", err)
		}
		fmt.Println("✅ 服务已安全关闭")
	}()

	fmt.Println("✅ RocketMQ 消息队列 + 定时检查器初始化完成")
}
