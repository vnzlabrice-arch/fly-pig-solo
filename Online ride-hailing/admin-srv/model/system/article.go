package system

import (
	"time"
)

// Article 文章表
type Article struct {
	// 注意：gorm.Model 自带 ID/CreatedAt/UpdatedAt/DeletedAt
	// 你这里重复定义了 ID 和 CreatedAt，我保留你自定义写法，去掉 Model 避免冲突
	ID         int64     `gorm:"primaryKey;autoIncrement;column:id;comment:文章ID"`
	Title      string    `gorm:"size:128;not null;column:title;comment:文章标题"`
	CategoryID int64     `gorm:"not null;column:category_id;comment:分类ID"`
	Content    string    `gorm:"type:text;not null;column:content;comment:文章内容"`
	Status     int8      `gorm:"default:1;column:status;comment:状态 1正常 0禁用"`
	CreatedAt  time.Time `gorm:"autoCreateTime;column:created_at;comment:创建时间"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;column:updated_at;comment:更新时间"`
	// 软删除字段（可选，需要就加）
	// DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at;comment:删除时间"`
}

func (Article) TableName() string {
	return "articles"
}
