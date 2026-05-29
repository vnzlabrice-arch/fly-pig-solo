package model

import "time"

// PassengerSupportMessage 乘客客服消息
type PassengerSupportMessage struct {
	ID         int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	ChatID     int64     `Gorm:"not null;column:chat_id"`
	SenderType int8      `Gorm:"not null;column:sender_type"`
	Content    string    `Gorm:"type:text;column:content"`
	CreatedAt  time.Time `Gorm:"autoCreateTime;column:created_at"`
}

func (PassengerSupportMessage) TableName() string {
	return "passenger_support_messages"
}
