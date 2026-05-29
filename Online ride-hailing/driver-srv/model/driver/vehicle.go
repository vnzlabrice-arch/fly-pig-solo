package driver

// DriverVehicle 车辆信息
type DriverVehicle struct {
	ID       int64  `Gorm:"primaryKey;autoIncrement;column:id"` // 主键
	DriverID int64  `Gorm:"not null;column:driver_id"`          // 司机ID
	PlateNo  string `Gorm:"size:16;not null;column:plate_no"`   // 车牌号
	Brand    string `Gorm:"size:32;column:brand"`               // 车辆品牌
	Model    string `Gorm:"size:32;column:model"`               // 车辆型号
	Color    string `Gorm:"size:16;column:color"`               // 车辆颜色
	Status   int8   `Gorm:"default:1;column:status"`            // 车辆状态
}

func (DriverVehicle) TableName() string {
	return "driver_vehicles"
}
