package driver

import "time"

// DriverOnlineStatus 司机在线状态
type DriverOnlineStatus struct {
	ID           int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	DriverID     int64     `Gorm:"unique;not null;column:driver_id"`
	OnlineStatus int8      `Gorm:"default:0;column:online_status"`
	AcceptOrder  int8      `Gorm:"default:0;column:accept_order"`
	Lng          float64   `Gorm:"type:decimal(10,6);default:0;column:lng"`
	Lat          float64   `Gorm:"type:decimal(10,6);default:0;column:lat"`
	CarStatus    int8      `Gorm:"default:1;column:car_status"`
	OnlineType   int8      `Gorm:"default:1;column:online_type"`
	RestReminder string    `Gorm:"size:255;column:rest_reminder"`
	UpdatedAt    time.Time `Gorm:"autoUpdateTime;column:updated_at"`
}

func (DriverOnlineStatus) TableName() string {
	return "driver_online_status"
}
