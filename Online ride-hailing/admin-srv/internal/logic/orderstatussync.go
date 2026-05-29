package logic

import (
	"errors"
	"time"

	"admin-srv/model/order"
)

// buildOrderStatusUpdates
// 作用：根据订单当前状态，构建订单更新所需的字段映射（map）
// 用于订单状态流转时，生成需要更新到数据库/消息队列的字段集合
// 参数：
//
//	status      - 订单状态枚举值
//	driverID    - 司机ID（有值时才会加入更新字段）
//	finalPrice  - 订单最终价格（仅完成状态需要）
//	cancelReason- 取消原因（仅取消状态需要）
//	now         - 当前时间，用于记录到达/开始/完成时间
//
// 返回：
//
//	订单更新字段map、错误（仅状态非法时返回）
func buildOrderStatusUpdates(status int32, driverID int64, finalPrice float64, cancelReason string, now time.Time) (map[string]interface{}, error) {
	// 校验订单状态是否在合法范围内
	if status < order.OrderStatusPendingAccept || status > order.OrderStatusCancelled {
		return nil, errors.New("invalid order status")
	}

	// 初始化更新map，默认只更新订单状态（转成int8节省空间）
	updates := map[string]interface{}{
		"status": int8(status),
	}

	// 如果司机ID有效（大于0），则加入driver_id更新字段
	if driverID > 0 {
		updates["driver_id"] = driverID
	}

	// 根据不同订单状态，追加对应的时间/价格/取消原因等字段
	switch status {
	case order.OrderStatusArrived:
		// 司机已到达：更新取货时间
		updates["pickup_time"] = &now
	case order.OrderStatusInProgress:
		// 订单进行中：更新开始服务时间
		updates["start_time"] = &now
	case order.OrderStatusCompleted:
		// 订单已完成：更新结束时间
		updates["end_time"] = &now
		// 最终价格>0时才更新，避免无效价格覆盖
		if finalPrice > 0 {
			updates["final_price"] = finalPrice
		}
	case order.OrderStatusCancelled:
		// 订单已取消：有取消原因时才更新该字段
		if cancelReason != "" {
			updates["cancel_reason"] = cancelReason
		}
	}

	// 返回构建好的订单更新字段
	return updates, nil
}
