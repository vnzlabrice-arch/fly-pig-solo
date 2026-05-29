package model

import "time"

type UserCouponLimit struct {
	ID         int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	UserID     int64     `Gorm:"not null;column:user_id"`
	TemplateID int64     `Gorm:"not null;column:template_id"`
	TodayCount int       `Gorm:"default:0;column:today_count"` // 今日领取数
	TotalCount int       `Gorm:"default:0;column:total_count"` // 累计领取数
	UpdatedAt  time.Time `Gorm:"autoUpdateTime;column:updated_at"`
}
