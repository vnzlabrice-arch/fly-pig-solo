package driver

import (
	"gorm.io/gorm"
)

// DriverWithdraw 司机提现记录
type DriverWithdraw struct {
	gorm.Model
	DriverID int64   `Gorm:"not null;column:driver_id"`                 // 司机ID
	Amount   float64 `Gorm:"type:decimal(10,2);not null;column:amount"` // 提现金额
	BankName string  `gorm:"size:50;not null;column:bank_name"`         // 银行名称
	BankCard string  `Gorm:"size:32;column:bank_card"`                  // 银行卡号
	BankUser string  `gorm:"size:20;not null;column:bank_user"`         // 开户人
	Status   int8    `Gorm:"default:1;column:status"`                   // 状态 ：1-申请中 2-已打款 3-失败
}

func (DriverWithdraw) TableName() string {
	return "driver_withdraws"
}
