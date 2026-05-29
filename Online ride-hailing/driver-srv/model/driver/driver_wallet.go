package driver

import (
	"gorm.io/gorm"
)

// DriverWallet 司机钱包
type DriverWallet struct {
	gorm.Model
	DriverID     int64   `Gorm:"unique;not null;column:driver_id"`                 // 司机ID
	Balance      float64 `Gorm:"type:decimal(12,2);default:0;column:balance"`      // 余额
	Withdrawable float64 `Gorm:"type:decimal(12,2);default:0;column:withdrawable"` // 可提现金额
	Frozen       float64 `Gorm:"type:decimal(12,2);default:0;column:frozen"`       // 冻结金额
	TotalIncome  float64 `gorm:"default:0.00;column:total_income"`                 // 总收入
}

func (DriverWallet) TableName() string {
	return "driver_wallets"
}

func (w *DriverWallet) CreateData(db *gorm.DB) error {
	return db.Debug().Create(w).Error
}
