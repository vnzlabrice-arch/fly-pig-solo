package driver

import (
	"gorm.io/gorm"
)

// DriverMessage 司机消息
type DriverMessage struct {
	gorm.Model
	DriverID int64  `Gorm:"not null;column:driver_id"` // 司机ID
	OrderID  int64  `gorm:"default:0;column:order_id"` // 关联订单ID
	Title    string `Gorm:"size:64;column:title"`      // 消息标题
	Content  string `Gorm:"size:512;column:content"`   // 消息内容
	Type     int8   `gorm:"default:1;column:type"`     // 1订单通知 2系统通知 3提现通知
	IsRead   int8   `Gorm:"default:0;column:is_read"`  // 是否已读:1-未读 2-已读
}

func (DriverMessage) TableName() string {
	return "driver_messages"
}
