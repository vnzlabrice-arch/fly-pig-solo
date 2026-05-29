package system

import "gorm.io/gorm"

// AdminRoleMenu 角色菜单关联
type AdminRoleMenu struct {
	gorm.Model
	ID     int64 `Gorm:"primaryKey;autoIncrement;column:id"`
	RoleID int64 `Gorm:"not null;column:role_id"`
	MenuID int64 `Gorm:"not null;column:menu_id"`
}

func (AdminRoleMenu) TableName() string {
	return "admin_role_menus"
}
