package model

import "time"

type CouponGrantTask struct {
	ID         int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	TemplateID int64     `Gorm:"not null;column:template_id"`
	GrantType  int8      `Gorm:"not null;column:grant_type"` // 发放类型 1-新人 2-充值 3-活动 4-手动
	GrantNum   int       `Gorm:"default:1;column:grant_num"` // 每次发放张数
	Status     int8      `Gorm:"default:1;column:status"`    // 1-启用 0-禁用
	CreatedAt  time.Time `Gorm:"autoCreateTime;column:created_at"`
	UpdatedAt  time.Time `Gorm:"autoUpdateTime;column:updated_at"`
}
