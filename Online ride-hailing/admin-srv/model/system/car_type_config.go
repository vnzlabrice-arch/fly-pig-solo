package system

import "gorm.io/gorm"

// CarTypeConfig 车辆类型配置表
// 存储不同车型的基础价格、里程单价、时长单价等计费规则
type CarTypeConfig struct {
	gorm.Model
	ID          int64   `gorm:"primaryKey;autoIncrement;column:id" comment:"主键ID"`
	TypeName    string  `gorm:"size:32;not null;column:type_name" comment:"车辆类型名称，如：经济型、舒适型、商务型"`
	BasePrice   float64 `gorm:"type:decimal(10,2);default:0;column:base_price" comment:"基础价格（起步价）"`
	KmPrice     float64 `gorm:"type:decimal(10,2);default:0;column:km_price" comment:"每公里单价"`
	MinutePrice float64 `gorm:"type:decimal(10,2);default:0;column:minute_price" comment:"每分钟单价（时长费）"`
	Status      int8    `gorm:"default:1;column:status" comment:"状态：1-启用 2-禁用"`
}

func (CarTypeConfig) TableName() string {
	return "car_type_configs"
}
