package model

import "time"

// Service 服务表
type Service struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id"`
	Name      string    `gorm:"type:varchar(50);column:name"`
	Price     float64   `gorm:"type:decimal(10,2);column:price"`
	CarType   int8      `gorm:"type:tinyint;column:car_type"`
	Status    int8      `gorm:"type:tinyint;default:1;column:status"`
	Sort      int       `gorm:"type:int;default:0;column:sort"`
	Sales     int       `gorm:"type:int;default:0;column:sales"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Service) TableName() string {
	return "service"
}
