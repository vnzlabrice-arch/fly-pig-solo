package model

import "time"

// PassengerTripSafetyLog 行程安全记录表 - 存储行程中的安全事件记录
type PassengerTripSafetyLog struct {
	ID          int64     `Gorm:"primaryKey;autoIncrement;column:id"` // 记录ID，主键，自增
	OrderID     string    `Gorm:"size:32;not null;column:order_id"`   // 关联订单ID
	PassengerID int64     `Gorm:"not null;column:passenger_id"`       // 用户ID，关联用户表
	LogType     int8      `Gorm:"not null;column:log_type"`           // 记录类型：1-行程分享，2-紧急求助，3-行程录音，4-行程录像，5-安全中心操作
	ContentURL  string    `Gorm:"size:255;column:content_url"`        // 内容URL（如录音、录像文件地址）
	ContactInfo string    `Gorm:"size:255;column:contact_info"`       // 紧急联系人信息
	Description string    `Gorm:"size:255;column:description"`        // 事件描述
	CreatedAt   time.Time `Gorm:"autoCreateTime;column:created_at"`   // 创建时间
}
