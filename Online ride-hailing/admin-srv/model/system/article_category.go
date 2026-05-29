package system

import "gorm.io/gorm"

// ArticleCategory 文章分类
type ArticleCategory struct {
	gorm.Model
	ID       int64  `Gorm:"primaryKey;autoIncrement;column:id"`
	Name     string `Gorm:"size:64;not null;column:name"`
	ParentID int64  `Gorm:"default:0;column:parent_id"`
	Status   int8   `Gorm:"default:1;column:status"`
}

func (ArticleCategory) TableName() string {
	return "article_categories"
}
