package driver

import (
	"time"

	"gorm.io/gorm"
)

// DriverProfile 司机资料详情
type DriverProfile struct {
	ID            int64      `Gorm:"primaryKey;autoIncrement;column:id"` // 主键
	DriverID      int64      `Gorm:"unique;not null;column:driver_id"`   // 司机ID
	LicenseNo     string     `Gorm:"size:32;column:license_no"`          // 行驶证号
	LicenseExpire *time.Time `Gorm:"type:date;column:license_expire"`    // 行驶证过期时间
	DriveYears    int        `Gorm:"default:0;column:drive_years"`       // 驾驶年限
}

func (DriverProfile) TableName() string {
	return "driver_profiles"
}

func (p *DriverProfile) UpdateData(db *gorm.DB) error {
	return db.Debug().Updates(p).Error
}
