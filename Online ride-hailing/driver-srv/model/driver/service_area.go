package driver

// DriverServiceArea 司机服务区域
type DriverServiceArea struct {
	ID       int64  `Gorm:"primaryKey;autoIncrement;column:id"` // 主键
	DriverID int64  `Gorm:"not null;column:driver_id"`          // 司机ID
	CityCode string `Gorm:"size:16;not null;column:city_code"`  // 城市编码
	Status   int8   `Gorm:"default:1;column:status"`            // 服务状态
}

func (DriverServiceArea) TableName() string {
	return "driver_service_areas"
}
