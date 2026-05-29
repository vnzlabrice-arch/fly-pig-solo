package model

import "time"

// PassengerVerificationCode 乘客验证码表
type PassengerVerificationCode struct {
	ID         int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	Phone      string    `Gorm:"size:11;not null;column:phone"`
	Code       string    `Gorm:"size:8;not null;column:code"`
	ExpireTime time.Time `Gorm:"not null;column:expire_time"`
	IsUsed     int8      `Gorm:"default:0;column:is_used"`
	CreatedAt  time.Time `Gorm:"autoCreateTime;column:created_at"`
}

func (PassengerVerificationCode) TableName() string {
	return "passenger_verification_codes"
}
