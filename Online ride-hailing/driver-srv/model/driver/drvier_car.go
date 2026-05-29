package driver

import (
	"gorm.io/gorm"
)

// 车辆认证表
type DriverCar struct {
	gorm.Model
	DriverID       int64  `gorm:"not null;column:driver_id"`                // 司机ID
	CarPlate       string `gorm:"size:20;not null;column:car_plate"`        // 车牌号
	CarModel       string `gorm:"size:50;not null;column:car_model"`        // 车型
	CarColor       string `gorm:"size:20;not null;column:car_color"`        // 颜色
	DrivingLicense string `gorm:"size:255;not null;column:driving_license"` // 行驶证
	CarImg         string `gorm:"size:255;not null;column:car_img"`         // 车辆照片
	Status         int8   `gorm:"default:0;column:status"`                  // 1-未提交 2-审核中 3-通过 4-驳回
	RejectReason   string `gorm:"size:255;column:reject_reason"`            // 驳回原因
}

func (c *DriverCar) CreateData(db *gorm.DB) error {
	return db.Debug().Create(c).Error
}
