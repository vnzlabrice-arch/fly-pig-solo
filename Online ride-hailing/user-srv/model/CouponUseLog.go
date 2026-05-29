package model

import "time"

type CouponUseLog struct {
	ID           int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	UserID       int64     `Gorm:"not null;column:user_id"`
	TemplateID   int64     `Gorm:"not null;column:template_id"`
	UserCouponID int64     `Gorm:"not null;column:user_coupon_id"`
	OrderNo      string    `Gorm:"size:64;not null;column:order_no"`
	OrderAmount  float64   `Gorm:"type:decimal(10,2);not null;column:order_amount"`  // 订单原价
	ReduceAmount float64   `Gorm:"type:decimal(10,2);not null;column:reduce_amount"` // 减免金额
	UseTime      time.Time `Gorm:"autoCreateTime;column:use_time"`
}
