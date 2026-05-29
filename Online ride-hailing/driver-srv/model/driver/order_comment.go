package driver

import (
	"gorm.io/gorm"
)

// OrderComment 订单评价表
type OrderComment struct {
	gorm.Model
	DriverID         int64  `gorm:"not null;column:driver_id"`         // 司机ID
	OrderID          int64  `gorm:"not null;column:order_id"`          // 订单ID
	PassengerID      int64  `gorm:"not null;column:passenger_id"`      // 乘客ID
	DriverScore      int8   `gorm:"not null;column:driver_score"`      // 乘客对司机评分 1-5
	PassengerScore   int8   `gorm:"default:0;column:passenger_score"`  // 司机对乘客评分 1-5
	DriverContent    string `gorm:"size:255;column:driver_content"`    // 司机评价乘客内容
	PassengerContent string `gorm:"size:255;column:passenger_content"` // 乘客评价司机内容
}
