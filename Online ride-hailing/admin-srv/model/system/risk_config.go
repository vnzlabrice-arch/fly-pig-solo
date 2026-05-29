package system

import "gorm.io/gorm"

// RiskConfig 风控配置
type RiskConfig struct {
	gorm.Model
	ID            int64 `Gorm:"primaryKey;autoIncrement;column:id"`
	NeedRealname  int8  `Gorm:"default:1;column:need_realname"`
	BlacklistDays int   `Gorm:"default:7;column:blacklist_days"`
}

func (RiskConfig) TableName() string {
	return "risk_configs"
}
