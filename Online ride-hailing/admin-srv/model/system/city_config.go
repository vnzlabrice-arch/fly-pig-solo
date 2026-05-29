package system

import "gorm.io/gorm"

// CityConfig 城市配置
type CityConfig struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey;autoIncrement;column:id;comment:主键ID"`
	CityCode string `gorm:"size:16;unique;not null;column:city_code;comment:城市编码（唯一）"`
	CityName string `gorm:"size:32;not null;column:city_name;comment:城市名称"`
	Status   int8   `gorm:"default:1;column:status;comment:状态 1=启用 0=禁用"`
}

func (CityConfig) TableName() string {
	return "city_configs"
}
