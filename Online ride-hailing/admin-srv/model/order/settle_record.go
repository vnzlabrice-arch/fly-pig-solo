package order

import (
	"time"

	"gorm.io/gorm"
)

// SettleRecord 订单结算
type SettleRecord struct {
	gorm.Model
	ID         int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	OrderID    string    `Gorm:"size:64;unique;not null;column:order_id"`
	RealAmount float64   `Gorm:"type:decimal(10,2);not null;column:real_amount"`
	Status     int8      `Gorm:"not null;column:status"`
	CreatedAt  time.Time `Gorm:"autoCreateTime;column:created_at"`
}

func (SettleRecord) TableName() string {
	return "settle_records"
}
