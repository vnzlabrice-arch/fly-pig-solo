package model

import "time"

// PassengerMessage 乘客系统消息
type PassengerMessage struct {
	ID          int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	PassengerID int64     `Gorm:"not null;column:passenger_id"`
	Title       string    `Gorm:"size:64;column:title"`
	Content     string    `Gorm:"size:512;column:content"`
	MsgType     int8      `Gorm:"default:0;column:msg_type"`
	IsRead      int8      `Gorm:"default:0;column:is_read"`
	CreatedAt   time.Time `Gorm:"autoCreateTime;column:created_at"`
}

func (PassengerMessage) TableName() string {
	return "passenger_messages"
}
