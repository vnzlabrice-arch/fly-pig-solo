package system

import "gorm.io/gorm"

// AdminMenu 菜单
// 用于管理后台左侧菜单、按钮、路由的层级结构
type AdminMenu struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey;autoIncrement;column:id" comment:"菜单主键ID"`
	ParentID int64  `gorm:"default:0;column:parent_id" comment:"父级菜单ID，0=顶级菜单"`
	Name     string `gorm:"size:32;not null;column:name" comment:"菜单名称（显示用）"`
	Path     string `gorm:"size:128;column:path" comment:"前端路由路径"`
	Icon     string `gorm:"size:32;column:icon" comment:"菜单图标"`
	Sort     int    `gorm:"default:0;column:sort" comment:"排序值，数字越小越靠前"`
}

func (AdminMenu) TableName() string {
	return "admin_menus"
}
