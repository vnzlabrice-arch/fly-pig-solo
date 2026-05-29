package system

import "gorm.io/gorm"

// SystemDict 系统字典
type SystemDict struct {
	gorm.Model
	ID        int64  `Gorm:"primaryKey;autoIncrement;column:id"`
	DictType  string `Gorm:"size:32;not null;column:dict_type"`
	DictKey   string `Gorm:"size:64;not null;column:dict_key"`
	DictValue string `Gorm:"size:64;not null;column:dict_value"`
	Sort      int    `Gorm:"default:0;column:sort"`
}

func (SystemDict) TableName() string {
	return "system_dicts"
}
