package model

import "time"

// PassengerMemberBenefit 会员权益表 - 存储用户的会员权益信息
type PassengerMemberBenefit struct {
	ID          int64      `Gorm:"primaryKey;autoIncrement;column:id"`   // 权益记录ID，主键，自增
	PassengerID int64      `Gorm:"not null;column:passenger_id"`         // 用户ID，关联用户表
	BenefitType string     `Gorm:"size:50;not null;column:benefit_type"` // 权益类型：如"免费升级"、"优先派单"、"专属客服"
	Status      int8       `Gorm:"default:1;column:status"`              // 权益状态：1-有效，2-已使用，3-已过期
	ExpireTime  *time.Time `Gorm:"column:expire_time"`                   // 权益过期时间
	CreatedAt   time.Time  `Gorm:"autoCreateTime;column:created_at"`     // 创建时间（获得权益时间）
}
