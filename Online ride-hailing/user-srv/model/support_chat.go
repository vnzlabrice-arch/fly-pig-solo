package model

import "time"

// PassengerSupportChat 乘客客服会话
type PassengerSupportChat struct {
	ID          int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	PassengerID int64     `Gorm:"not null;column:passenger_id"`
	Status      int8      `Gorm:"default:1;column:status"`
	CreatedAt   time.Time `Gorm:"autoCreateTime;column:created_at"`
	UpdatedAt   time.Time `Gorm:"autoUpdateTime;column:updated_at"`
}

func (PassengerSupportChat) TableName() string {
	return "passenger_support_chats"
}
