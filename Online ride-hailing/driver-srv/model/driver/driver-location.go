package driver

import (
	"gorm.io/gorm"
)

// 司机实时位置表
type DriverLocation struct {
	gorm.Model
	DriverID  int64   `gorm:"not null;column:driver_id"`  // 司机ID
	Lng       float64 `gorm:"column:lng"`                 // 经度
	Lat       float64 `gorm:"column:lat"`                 // 纬度
	Speed     float64 `gorm:"default:0;column:speed"`     // 速度km/h
	Direction int     `gorm:"default:0;column:direction"` // 方向角度
}
