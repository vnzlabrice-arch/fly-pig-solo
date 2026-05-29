package system

import (
	"time"

	"gorm.io/gorm"
)

// AdminOperationLog 操作日志
type AdminOperationLog struct {
	gorm.Model
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id;comment:主键ID"`
	AdminID   int64     `gorm:"not null;column:admin_id;comment:操作人ID"`
	Module    string    `gorm:"size:32;column:module;comment:操作模块"`
	Operation string    `gorm:"size:32;column:operation;comment:操作类型"`
	Detail    string    `gorm:"size:255;column:detail;comment:操作详情"`
	IP        string    `gorm:"size:32;column:ip;comment:操作IP地址"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;comment:创建时间"`
}

func (AdminOperationLog) TableName() string {
	return "admin_operation_logs"
}
