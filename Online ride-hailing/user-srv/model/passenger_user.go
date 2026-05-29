package model

import (
	"time"
)

// PassengerUser 用户表 - 存储乘客用户的基本信息
type PassengerUser struct {
	ID            int64      `Gorm:"primaryKey;autoIncrement;column:id"`          // 用户ID，主键，自增
	Phone         string     `Gorm:"size:20;unique;not null;column:phone"`        // 手机号，唯一标识，不为空
	Nickname      string     `Gorm:"size:50;column:nickname"`                     // 用户昵称
	AvatarURL     string     `Gorm:"size:255;column:avatar_url"`                  // 头像URL地址
	RealName      string     `Gorm:"size:50;column:real_name"`                    // 真实姓名（实名认证后存储）
	IDCardHash    string     `Gorm:"size:255;column:id_card_hash"`                // 身份证号哈希值（加密存储）
	Gender        int8       `Gorm:"default:0;column:gender"`                     // 性别：0-未知，1-男，2-女
	MemberLevel   int8       `Gorm:"default:0;column:member_level"`               // 会员等级：0-普通用户，1-白银，2-黄金，3-铂金，4-钻石
	Balance       float64    `Gorm:"type:decimal(10,2);default:0;column:balance"` // 账户余额（元）
	FlowerCoin    int        `Gorm:"default:0;column:flower_coin"`                // 花小猪金币数量
	Integral      int        `Gorm:"default:0;column:integral"`                   // 积分数量
	Status        int8       `Gorm:"not null;default:1;column:status"`            // 用户状态：1-正常，2-冻结，3-注销
	LastLoginTime *time.Time `Gorm:"column:last_login_time"`                      // 最后登录时间
	CreatedAt     time.Time  `Gorm:"autoCreateTime;column:created_at"`            // 创建时间（注册时间）
	UpdatedAt     time.Time  `Gorm:"autoUpdateTime;column:updated_at"`            // 更新时间
}
