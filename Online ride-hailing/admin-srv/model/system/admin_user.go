package system

import (
	"time"

	"gorm.io/gorm"
)

// AdminUser 管理员用户
type AdminUser struct {
	gorm.Model
	ID            int64      `Gorm:"primaryKey;autoIncrement;column:id"`
	Username      string     `Gorm:"size:32;unique;not null;column:username"`
	Password      string     `Gorm:"size:128;not null;column:password"`
	RoleID        int64      `Gorm:"not null;column:role_id"`
	Status        int8       `Gorm:"default:1;column:status"`
	LastLoginTime *time.Time `Gorm:"column:last_login_time"`
	CreatedAt     time.Time  `Gorm:"autoCreateTime;column:created_at"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}
