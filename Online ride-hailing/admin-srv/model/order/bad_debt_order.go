package order

import "gorm.io/gorm"

// BadDebtOrder 坏账订单
type BadDebtOrder struct {
	gorm.Model
	ID          int64   `Gorm:"primaryKey;autoIncrement;column:id"`
	OrderID     string  `Gorm:"size:64;unique;not null;column:order_id"`
	DebtAmount  float64 `Gorm:"type:decimal(10,2);not null;column:debt_amount"`
	OverdueDays int     `Gorm:"default:0;column:overdue_days"`
}

func (BadDebtOrder) TableName() string {
	return "bad_debt_orders"
}
