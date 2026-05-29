package driver

import "time"

// DriverOrderRecommend 司机订单推荐表
type DriverOrderRecommend struct {
	ID             int64      `Gorm:"primaryKey;autoIncrement;column:id"`
	DriverID       int64      `Gorm:"not null;column:driver_id"`
	OrderID        string     `Gorm:"size:64;not null;column:order_id"`
	Distance       int        `Gorm:"default:0;column:distance"`
	EstimatePrice  float64    `Gorm:"type:decimal(10,2);default:0;column:estimate_price"`
	Status         int8       `Gorm:"default:1;column:status"`
	OrderLockTime  *time.Time `Gorm:"column:order_lock_time"`
	AntiFraudToken string     `Gorm:"size:64;column:anti_fraud_token"`
	CreatedAt      time.Time  `Gorm:"autoCreateTime;column:created_at"`
}

func (DriverOrderRecommend) TableName() string {
	return "driver_order_recommends"
}
