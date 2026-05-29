package order

import "time"

// PassengerOrder 乘客订单（admin-srv 用于后台派单）
type PassengerOrder struct {
	OrderID        string     `gorm:"primaryKey;size:32;column:order_id"`
	PassengerID    int64      `gorm:"not null;column:passenger_id"`
	PassengerName  string     `gorm:"size:50;column:passenger_name"`
	PassengerPhone string     `gorm:"size:20;column:passenger_phone"`
	DriverID       int64      `gorm:"column:driver_id"`
	OrderType      int8       `gorm:"not null;column:order_type"`
	CarType        string     `gorm:"size:20;not null;column:car_type"`
	Status         int8       `gorm:"not null;column:status"`
	StartLng       float64    `gorm:"type:decimal(10,6);not null;column:start_lng"`
	StartLat       float64    `gorm:"type:decimal(10,6);not null;column:start_lat"`
	StartAddress   string     `gorm:"size:255;not null;column:start_address"`
	EndLng         float64    `gorm:"type:decimal(10,6);not null;column:end_lng"`
	EndLat         float64    `gorm:"type:decimal(10,6);not null;column:end_lat"`
	EndAddress     string     `gorm:"size:255;not null;column:end_address"`
	PassRemark     string     `gorm:"size:255;column:pass_remark"`
	EstimatedPrice float64    `gorm:"type:decimal(10,2);not null;column:estimated_price"`
	FinalPrice     float64    `gorm:"type:decimal(10,2);column:final_price"`
	PayStatus      int8       `gorm:"default:0;column:pay_status"`
	BookTime       time.Time  `gorm:"not null;column:book_time"`
	AppointTime    *time.Time `gorm:"column:appoint_time"`
	DispatchTime   *time.Time `gorm:"column:dispatch_time"` // 派单时间（后台派单时填充）
	PickupTime     *time.Time `gorm:"column:pickup_time"`
	StartTime      *time.Time `gorm:"column:start_time"`
	EndTime        *time.Time `gorm:"column:end_time"`
	CancelReason   string     `gorm:"size:255;column:cancel_reason"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime;column:updated_at"`
}

func (PassengerOrder) TableName() string {
	return "passenger_orders"
}

// 订单状态常量
const (
	OrderStatusPendingDispatch = 0 // 待派单（用户下单后等待后台分配司机）
	OrderStatusPendingAccept   = 1 // 待接单
	OrderStatusAccepted        = 2 // 已接单（司机已接单）
	OrderStatusArrived         = 3 // 司机已到达
	OrderStatusInProgress      = 4 // 行程中
	OrderStatusCompleted       = 5 // 已完成
	OrderStatusCancelled       = 6 // 已取消
)

// 车型映射
var CarTypeMap = map[int32]string{
	1: "经济型",
	2: "舒适型",
	3: "商务型",
	4: "豪华型",
}
