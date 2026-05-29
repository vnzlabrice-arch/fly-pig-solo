package driver

import (
	"gorm.io/gorm"
)

// 司机反馈投诉表

type DriverFeedback struct {
	gorm.Model
	DriverID    int64  `gorm:"not null;column:driver_id"`         // 司机ID
	OrderID     int64  `gorm:"default:0;column:order_id"`         // 关联订单ID
	PassengerID int64  `gorm:"default:0;column:passenger_id"`     // 关联乘客ID
	Content     string `gorm:"type:text;not null;column:content"` // 反馈内容
	Img         string `gorm:"size:500;column:img"`               // 图片
	Status      int8   `gorm:"default:0;column:status"`           // 0待处理 1已处理
	Reply       string `gorm:"size:255;column:reply"`             // 管理员回复
}
