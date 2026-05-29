package model

import "time"

// UserCoupon 用户优惠券表（用户实际持有的券）
type UserCoupon struct {
	ID         int64      `Gorm:"primaryKey;autoIncrement;column:id"`
	UserID     int64      `Gorm:"not null;column:user_id"`           // 用户ID
	TemplateID int64      `Gorm:"not null;column:template_id"`       // 模板ID
	CouponNo   string     `Gorm:"size:64;not null;column:coupon_no"` // 券唯一编码
	Status     int8       `Gorm:"not null;column:status"`            // 状态 1-未使用 2-已使用 3-已过期 4-已作废
	UsedTime   *time.Time `Gorm:"column:used_time"`                  // 使用时间
	OrderNo    string     `Gorm:"size:64;column:order_no"`           // 核销订单号
	StartTime  time.Time  `Gorm:"not null;column:start_time"`        // 生效时间
	EndTime    time.Time  `Gorm:"not null;column:end_time"`          // 过期时间
	CreatedAt  time.Time  `Gorm:"autoCreateTime;column:created_at"`
	UpdatedAt  time.Time  `Gorm:"autoUpdateTime;column:updated_at"`
}
