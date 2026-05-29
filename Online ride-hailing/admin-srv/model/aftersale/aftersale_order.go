package aftersale

import (
	"time"

	"gorm.io/gorm"
)

// AftersaleOrder 售后工单
type AftersaleOrder struct {
	gorm.Model
	ID          int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	AftersaleNo string    `Gorm:"size:32;unique;not null;column:aftersale_no"`
	OrderID     string    `Gorm:"size:64;not null;column:order_id"`
	Type        int8      `Gorm:"not null;column:type"`
	Status      int8      `Gorm:"not null;column:status"`
	CreatedAt   time.Time `Gorm:"autoCreateTime;column:created_at"`
}

func (AftersaleOrder) TableName() string {
	return "aftersale_orders"
}
