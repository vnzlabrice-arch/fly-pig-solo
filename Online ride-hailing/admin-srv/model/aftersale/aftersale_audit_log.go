package aftersale

import (
	"time"

	"gorm.io/gorm"
)

// AftersaleAuditLog 售后审核日志
type AftersaleAuditLog struct {
	gorm.Model
	ID          int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	AftersaleID int64     `Gorm:"not null;column:aftersale_id"`
	AdminID     int64     `Gorm:"not null;column:admin_id"`
	Status      int8      `Gorm:"not null;column:status"`
	Note        string    `Gorm:"size:255;column:note"`
	CreatedAt   time.Time `Gorm:"autoCreateTime;column:created_at"`
}

func (AftersaleAuditLog) TableName() string {
	return "aftersale_audit_logs"
}
