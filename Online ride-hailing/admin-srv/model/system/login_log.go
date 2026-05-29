package system

import (
	"time"
)

// AdminLoginLog 登录日志（增强版）
type AdminLoginLog struct {
	ID          int64      `gorm:"primaryKey;autoIncrement;column:id"`
	AdminID     int64      `gorm:"not null;column:admin_id"`           // 管理员ID（失败时为0）
	AdminName   string     `gorm:"size:32;column:admin_name"`           // 登录用户名
	IP          string     `gorm:"size:45;column:ip"`                  // IP地址(IPv6最长45字符)
	Location    string     `gorm:"size:100;column:location"`            // IP地理位置
	DeviceType  string     `gorm:"size:20;column:device_type"`          // 设备类型: PC/Mobile/Tablet
	Browser     string     `gorm:"size:50;column:browser"`              // 浏览器
	OS          string     `gorm:"size:50;column:os"`                   // 操作系统
	UserAgent   string     `gorm:"size:500;column:user_agent"`          // 完整User-Agent
	Status      int8       `gorm:"not null;default:0;column:status"`    // 0=失败 1=成功 2=账号锁定
	FailReason  string     `gorm:"size:200;column:fail_reason"`         // 失败原因
	LoginTime   time.Time  `gorm:"autoCreateTime;column:login_time"`    // 登录时间
	CreatedAt   time.Time  `gorm:"autoCreateTime;column:created_at"`
}

// TableName 表名
func (AdminLoginLog) TableName() string {
	return "admin_login_logs"
}

const (
	LoginStatusFailed = 0 // 登录失败
	LoginStatusSuccess = 1 // 登录成功
	LoginStatusLocked = 2 // 账号被锁定
)

// LoginLogItem 日志记录项（用于简化日志创建）
type LoginLogItem struct {
	AdminID    int64
	Username   string
	IP         string
	Location   string
	DeviceType string
	Browser    string
	OS         string
	UserAgent  string
	Status     int8
	FailReason string
}
