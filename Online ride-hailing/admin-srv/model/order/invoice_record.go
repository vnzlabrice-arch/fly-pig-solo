package order

import (
	"time"

	"gorm.io/gorm"
)

// InvoiceRecord 发票记录
type InvoiceRecord struct {
	gorm.Model
	ID        int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	InvoiceNo string    `Gorm:"size:32;unique;not null;column:invoice_no"`
	OrderID   string    `Gorm:"size:64;not null;column:order_id"`
	Status    int8      `Gorm:"not null;column:status"`
	CreatedAt time.Time `Gorm:"autoCreateTime;column:created_at"`
}

func (InvoiceRecord) TableName() string {
	return "invoice_records"
}
