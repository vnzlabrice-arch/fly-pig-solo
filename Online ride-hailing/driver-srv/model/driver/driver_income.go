package driver

import (
	"gorm.io/gorm"
)

// DriverIncome 司机收入明细
type DriverIncome struct {
	gorm.Model
	DriverID     int64   `Gorm:"not null;column:driver_id"`                         // 司机ID
	OrderID      string  `Gorm:"size:64;column:order_id"`                           // 订单ID
	PassengerID  int64   `gorm:"not null;column:passenger_id"`                      // 关联乘客ID
	OrderNo      string  `gorm:"size:32;not null;column:order_no"`                  // 订单号
	TotalFee     float64 `Gorm:"type:decimal(10,2);default:0;column:total_fee"`     // 订单总金额
	PlatformFee  float64 `Gorm:"type:decimal(10,2);default:0;column:platform_fee"`  // 平台费
	ActualIncome float64 `Gorm:"type:decimal(10,2);default:0;column:actual_income"` // 司机实际收入
	Type         int8    `gorm:"default:1;column:type"`                             // 1-订单收入 2-奖励 3-补贴
}

func (DriverIncome) TableName() string {
	return "driver_incomes"
}
