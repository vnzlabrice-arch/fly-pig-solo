package logic

import (
	"testing"
	"time"

	"admin-srv/model/order"
)

// TestBuildOrderStatusUpdatesByStatus
// 测试核心功能：根据不同订单状态，正确生成订单更新字段
func TestBuildOrderStatusUpdatesByStatus(t *testing.T) {
	// 固定一个测试时间，避免每次运行时间不同导致测试失败
	now := time.Date(2026, 5, 25, 9, 30, 0, 0, time.Local)

	// 定义测试用例切片：包含各种订单状态的场景
	tests := []struct {
		name          string  // 测试用例名称（用于区分场景）
		status        int32   // 订单状态
		driverID      int64   // 司机ID
		finalPrice    float64 // 最终价格
		cancelReason  string  // 取消原因
		wantTimeKey   string  // 期望生成的时间字段名（pickup_time/start_time/end_time）
		wantFinal     bool    // 是否期望生成 final_price 字段
		wantCancelled bool    // 是否期望生成 cancel_reason 字段
	}{
		// 用例1：订单已接单 → 只更新状态+司机ID
		{
			name:     "accepted syncs driver and status",
			status:   order.OrderStatusAccepted,
			driverID: 1001,
		},
		// 用例2：司机已到达 → 更新状态+取货时间
		{
			name:        "arrived syncs pickup time",
			status:      order.OrderStatusArrived,
			wantTimeKey: "pickup_time",
		},
		// 用例3：订单进行中 → 更新状态+开始时间
		{
			name:        "in progress syncs start time",
			status:      order.OrderStatusInProgress,
			wantTimeKey: "start_time",
		},
		// 用例4：订单已完成 → 更新状态+结束时间+最终价格
		{
			name:        "completed syncs end time and final price",
			status:      order.OrderStatusCompleted,
			finalPrice:  56.78,
			wantTimeKey: "end_time",
			wantFinal:   true,
		},
		// 用例5：订单已取消 → 更新状态+取消原因
		{
			name:          "cancelled syncs cancel reason",
			status:        order.OrderStatusCancelled,
			cancelReason:  "passenger cancelled",
			wantCancelled: true,
		},
	}

	// 遍历所有测试用例，逐个执行
	for _, tt := range tests {
		// t.Run 子测试：每个用例独立运行，互不影响
		t.Run(tt.name, func(t *testing.T) {
			// 调用被测试函数
			updates, err := buildOrderStatusUpdates(tt.status, tt.driverID, tt.finalPrice, tt.cancelReason, now)

			// 1. 校验：正常场景不应该返回错误
			if err != nil {
				t.Fatalf("buildOrderStatusUpdates() error = %v", err)
			}

			// 2. 校验：状态字段必须正确
			if got := updates["status"]; got != int8(tt.status) {
				t.Fatalf("status update = %v, want %d", got, tt.status)
			}

			// 3. 校验：有司机ID时，必须包含 driver_id 字段
			if tt.driverID > 0 && updates["driver_id"] != tt.driverID {
				t.Fatalf("driver_id update = %v, want %d", updates["driver_id"], tt.driverID)
			}

			// 4. 校验：对应状态必须生成正确的时间字段
			if tt.wantTimeKey != "" {
				got, ok := updates[tt.wantTimeKey]
				if !ok {
					t.Fatalf("%s update missing", tt.wantTimeKey)
				}
				// 校验时间是指针且时间值完全一致
				gotTime, ok := got.(*time.Time)
				if !ok || gotTime == nil || !gotTime.Equal(now) {
					t.Fatalf("%s update = %#v, want pointer to %v", tt.wantTimeKey, got, now)
				}
			}

			// 5. 校验：完成状态必须包含正确的最终价格
			if tt.wantFinal && updates["final_price"] != tt.finalPrice {
				t.Fatalf("final_price update = %v, want %.2f", updates["final_price"], tt.finalPrice)
			}

			// 6. 校验：取消状态必须包含正确的取消原因
			if tt.wantCancelled && updates["cancel_reason"] != tt.cancelReason {
				t.Fatalf("cancel_reason update = %v, want %q", updates["cancel_reason"], tt.cancelReason)
			}
		})
	}
}

// TestBuildOrderStatusUpdatesRejectsInvalidStatus
// 测试异常场景：非法订单状态必须返回错误
func TestBuildOrderStatusUpdatesRejectsInvalidStatus(t *testing.T) {
	// 传入一个不存在的状态 99
	_, err := buildOrderStatusUpdates(99, 0, 0, "", time.Now())
	// 必须返回错误，否则测试失败
	if err == nil {
		t.Fatal("expected invalid status error")
	}
}
