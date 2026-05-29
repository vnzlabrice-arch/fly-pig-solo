package system

import "gorm.io/gorm"

// SystemConfig 系统配置
type SystemConfig struct {
	gorm.Model
	ID          int64  `Gorm:"primaryKey;autoIncrement;column:id"`
	ConfigKey   string `Gorm:"size:64;unique;not null;column:config_key"`
	ConfigValue string `Gorm:"size:255;column:config_value"`
	Remark      string `Gorm:"size:128;column:remark"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}
