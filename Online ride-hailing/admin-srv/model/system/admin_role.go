package system

import "gorm.io/gorm"

// AdminRole 角色
type AdminRole struct {
	gorm.Model
	ID     int64  `Gorm:"primaryKey;autoIncrement;column:id"`
	Name   string `Gorm:"size:32;not null;column:name"`
	Remark string `Gorm:"size:128;column:remark"`
}

func (AdminRole) TableName() string {
	return "admin_roles"
}
