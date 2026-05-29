package model

import "time"

// PassengerAddressBook 用户地址簿表 - 存储用户的常用地址信息
type PassengerAddressBook struct {
	ID          int64     `Gorm:"primaryKey;autoIncrement;column:id"`     // 地址ID，主键，自增
	PassengerID int64     `Gorm:"not null;column:passenger_id"`           // 用户ID，关联用户表
	Tag         string    `Gorm:"size:20;column:tag"`                     // 地址标签：如"家"、"公司"、"学校"
	Address     string    `Gorm:"size:255;not null;column:address"`       // 详细地址描述
	Lng         float64   `Gorm:"type:decimal(10,6);not null;column:lng"` // 经度（经纬度坐标）
	Lat         float64   `Gorm:"type:decimal(10,6);not null;column:lat"` // 纬度（经纬度坐标）
	IsDefault   int8      `Gorm:"default:0;column:is_default"`            // 是否默认地址：0-否，1-是
	CreatedAt   time.Time `Gorm:"autoCreateTime;column:created_at"`       // 创建时间
}
