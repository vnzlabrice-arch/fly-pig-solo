package model

import "time"

// PassengerOrder 订单主表 - 存储乘客的出行订单信息
type PassengerOrder struct {
	OrderID        string     `Gorm:"primaryKey;size:32;column:order_id"`                 // 订单ID，主键，唯一标识
	PassengerID    int64      `Gorm:"not null;column:passenger_id"`                       // 乘客ID，关联用户表
	PassengerName  string     `Gorm:"size:50;column:passenger_name"`                      // 乘客姓名
	PassengerPhone string     `Gorm:"size:20;column:passenger_phone"`                     // 乘客电话
	CouponID       int64      `Gorm:"column:coupon_id"`                                   // 优惠券ID（使用优惠劵时填充）
	CouponName     string     `Gorm:"size:64;column:coupon_name"`                         // 优惠券名称
	DriverID       int64      `Gorm:"column:driver_id"`                                   // 司机ID（接单后填充）
	OrderType      int8       `Gorm:"not null;column:order_type"`                         // 订单类型：1-即时用车，2-预约用车，3-拼车
	CarType        string     `Gorm:"size:20;not null;column:car_type"`                   // 车型：如"经济型"、"舒适型"、"商务型"
	Status         int8       `Gorm:"not null;column:status"`                             // 订单状态：1-待接单，2-已接单，3-司机已到，4-行程中，5-已完成，6-已取消
	StartLng       float64    `Gorm:"type:decimal(10,6);not null;column:start_lng"`       // 起点经度
	StartLat       float64    `Gorm:"type:decimal(10,6);not null;column:start_lat"`       // 起点纬度
	StartAddress   string     `Gorm:"size:255;not null;column:start_address"`             // 起点地址描述
	EndLng         float64    `Gorm:"type:decimal(10,6);not null;column:end_lng"`         // 终点经度
	EndLat         float64    `Gorm:"type:decimal(10,6);not null;column:end_lat"`         // 终点纬度
	EndAddress     string     `Gorm:"size:255;not null;column:end_address"`               // 终点地址描述
	PassRemark     string     `Gorm:"size:255;column:pass_remark"`                        // 乘客备注信息
	EstimatedPrice float64    `Gorm:"type:decimal(10,2);not null;column:estimated_price"` // 预估价格（元）
	FinalPrice     float64    `Gorm:"type:decimal(10,2);column:final_price"`              // 最终价格（元，行程结束后计算）
	PayStatus      int8       `Gorm:"default:0;column:pay_status"`                        // 支付状态：0-未支付，1-已支付，2-部分支付，3-退款中，4-已退款
	BookTime       time.Time  `Gorm:"not null;column:book_time"`                          // 下单时间
	AppointTime    *time.Time `Gorm:"column:appoint_time"`                                // 预约时间（仅预约单有值）
	PickupTime     *time.Time `Gorm:"column:pickup_time"`                                 // 司机到达上车点时间
	StartTime      *time.Time `Gorm:"column:start_time"`                                  // 行程开始时间
	EndTime        *time.Time `Gorm:"column:end_time"`                                    // 行程结束时间
	CancelReason   string     `Gorm:"size:255;column:cancel_reason"`                      // 取消原因（取消订单时填写）
	CreatedAt      time.Time  `Gorm:"autoCreateTime;column:created_at"`                   // 创建时间
	UpdatedAt      time.Time  `Gorm:"autoUpdateTime;column:updated_at"`                   // 更新时间
}
