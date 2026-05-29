package system

import "gorm.io/gorm"

// AppUpdate APP升级配置
type AppUpdate struct {
	gorm.Model
	ID          int64  `Gorm:"primaryKey;autoIncrement;column:id"`
	Platform    int8   `Gorm:"not null;column:platform"`
	Version     string `Gorm:"size:32;not null;column:version"`
	DownloadURL string `Gorm:"size:255;column:download_url"`
	ForceUpdate int8   `Gorm:"default:0;column:force_update"`
	Content     string `Gorm:"type:text;column:content"`
}

func (AppUpdate) TableName() string {
	return "app_updates"
}
