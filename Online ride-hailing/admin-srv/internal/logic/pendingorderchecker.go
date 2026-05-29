package logic

import (
	"fmt"
	"sync"
	"time"

	"admin-srv/global"
	"admin-srv/model/order"
	"github.com/zeromicro/go-zero/core/logx"
)

// PendingOrderChecker 待派单定时检查器
// 老师要求：每秒查询待派单订单(status=0)，用于后台派单监控
type PendingOrderChecker struct {
	ticker    *time.Ticker
	stopCh    chan struct{}
	isRunning bool
	mu        sync.Mutex
}

var checker *PendingOrderChecker

// StartPendingOrderChecker 启动待派单定时检查器
// 每1秒查询一次 status=0 (待派单) 的订单
func StartPendingOrderChecker() {
	checker = &PendingOrderChecker{
		ticker: time.NewTicker(1 * time.Second), // 每秒执行
		stopCh: make(chan struct{}),
	}
	checker.isRunning = true
	go checker.run()
	logx.Infof("🔄 待派单定时检查器已启动（每秒查询 status=0 的订单）")
}

// Stop 停止检查器
func Stop() {
	if checker != nil {
		checker.stop()
	}
}

// run 主循环
func (c *PendingOrderChecker) run() {
	defer c.ticker.Stop()
	for {
		select {
		case <-c.ticker.C:
			c.checkPendingOrders()
		case <-c.stopCh:
			c.mu.Lock()
			c.isRunning = false
			c.mu.Unlock()
			logx.Info("⏹️  待派单定时检查器已停止")
			return
		}
	}
}

// stop 停止运行
func (c *PendingOrderChecker) stop() {
	close(c.stopCh)
}

// checkPendingOrders 查询待派单订单（status=0）
// 业务场景：用户下单后订单状态为0，后台需要知道有哪些订单待派单
func (c *PendingOrderChecker) checkPendingOrders() {
	var pendingOrders []struct {
		OrderID      string    `gorm:"column:order_id"`
		PassengerID   int64     `gorm:"column:passenger_id"`
		PassengerName string    `gorm:"column:passenger_name"`
		PassengerPhone string   `gorm:"column:passenger_phone"`
		StartAddress  string    `gorm:"column:start_address"`
		EndAddress    string    `gorm:"column:end_address"`
		CarType       string    `gorm:"column:car_type"`
		Status        int8      `gorm:"column:status"`
		CreatedAt     time.Time `gorm:"column:created_at"`
	}

	result := global.DB.Table("passenger_orders").
		Select("order_id, passenger_id, passenger_name, passenger_phone, start_address, end_address, car_type, status, created_at").
		Where("status = ?", order.OrderStatusPendingDispatch). // 只查 status=0 (待派单)
		Order("created_at ASC"). // 按下单时间升序，优先处理早下的单
		Limit(20). // 限制数量避免性能问题
		Find(&pendingOrders)

	if result.Error != nil {
		logx.Errorf("❌ 查询待派单失败: %v", result.Error)
		return
	}

	if len(pendingOrders) > 0 {
		logx.Infof("📋 当前待派单数量: %d 单 (status=0)", len(pendingOrders))
		
		// 打印前3个待派单信息供后台查看
		showCount := len(pendingOrders)
		if showCount > 3 {
			showCount = 3
		}
		for i := 0; i < showCount; i++ {
			o := pendingOrders[i]
			fmt.Printf("   [%d] 订单号:%s | 乘客:%s(%s) | %s → %s | 车型:%s | 下单时间:%s\n",
				i+1, o.OrderID, o.PassengerName, o.PassengerPhone,
				o.StartAddress, o.EndAddress, o.CarType,
				o.CreatedAt.Format("01-02 15:04:05"))
		}
		if len(pendingOrders) > 3 {
			fmt.Printf("   ... 还有 %d 个待派单\n", len(pendingOrders)-3)
		}
	}
}

// GetPendingOrderList 获取待派单列表（供API调用）
// 返回当前所有待派单的订单号和基本信息
func GetPendingOrderList() ([]map[string]interface{}, error) {
	var orders []struct {
		OrderID      string `gorm:"column:order_id"`
		PassengerName string `gorm:"column:passenger_name"`
		PassengerPhone string `gorm:"column:passenger_phone"`
		StartAddress  string `gorm:"column:start_address"`
		EndAddress    string `gorm:"column:end_address"`
		CreatedAt     time.Time `gorm:"column:created_at"`
	}

	result := global.DB.Table("passenger_orders").
		Select("order_id, passenger_name, passenger_phone, start_address, end_address, created_at").
		Where("status = ?", order.OrderStatusPendingDispatch).
		Order("created_at ASC").
		Find(&orders)

	if result.Error != nil {
		return nil, fmt.Errorf("查询待派单列表失败: %w", result.Error)
	}

	list := make([]map[string]interface{}, len(orders))
	for i, o := range orders {
		list[i] = map[string]interface{}{
			"order_id":       o.OrderID,
			"passenger_name": o.PassengerName,
			"passenger_phone":o.PassengerPhone,
			"start_address":  o.StartAddress,
			"end_address":    o.EndAddress,
			"created_at":     o.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return list, nil
}
