package order

import (
	"time"

	"gorm.io/gorm"
)

// PaymentOrder 付款单（退款/赔付）
type PaymentOrder struct {
	gorm.Model
	ID        int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	PaymentNo string    `Gorm:"size:32;unique;not null;column:payment_no"`
	Amount    float64   `Gorm:"type:decimal(10,2);not null;column:amount"`
	Status    int8      `Gorm:"not null;column:status"`
	CreatedAt time.Time `Gorm:"autoCreateTime;column:created_at"`
}

func (PaymentOrder) TableName() string {
	return "payment_orders"
}
