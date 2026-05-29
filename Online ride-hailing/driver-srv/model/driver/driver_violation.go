package driver

import (
	"gorm.io/gorm"
)

// DriverViolation 司机违规记录
type DriverViolation struct {
	gorm.Model
	DriverID      int64  `Gorm:"not null;column:driver_id"`      // 司机ID
	OrderID       string `Gorm:"size:64;column:order_id"`        // 订单ID
	PassengerID   int64  `gorm:"default:0;column:passenger_id"`  // 关联乘客ID
	ViolationType int8   `Gorm:"not null;column:violation_type"` // 违规类型 1-拒单 2-取消订单 3-绕路 4-被投诉
	Reason        string `Gorm:"size:255;column:reason"`         // 违规原因
	Score         int    `gorm:"default:0;column:score"`         // 扣分数
	//ScoreDeduct   int       `Gorm:"default:0;column:score_deduct"`      // 迣分
}

func (DriverViolation) TableName() string {
	return "driver_violations"
}
