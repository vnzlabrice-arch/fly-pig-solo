package driver

import (
	"gorm.io/gorm"
)

//系统规则表

type SysRule struct {
	gorm.Model
	RuleKey   string `gorm:"size:50;not null;unique;column:rule_key"` // 规则唯一标识
	RuleValue string `gorm:"size:500;not null;column:rule_value"`     // 规则内容
	Title     string `gorm:"size:50;not null;column:title"`           // 规则名称
}
