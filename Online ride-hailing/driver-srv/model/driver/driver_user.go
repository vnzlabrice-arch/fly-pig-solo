package driver

import (
	"gorm.io/gorm"
)

// DriverUser 司机主表
type DriverUser struct {
	gorm.Model
	Phone     string `Gorm:"size:11;unique;not null;column:phone"` // 手机号
	Nickname  string `Gorm:"size:32;column:nickname"`              // 昵称
	AvatarURL string `Gorm:"size:255;column:avatar_url"`           // 头像URL
	Password  string `gorm:"size:64;column:password"`              //密码
	Status    int8   `Gorm:"default:1;column:status"`              // 状态: 1-在线 2-离线

	//Gender        int8       `Gorm:"default:0;column:gender"`                             // 性别
	//ServiceScore  float64    `Gorm:"type:decimal(3,2);default:5.00;column:service_score"` // 服务评分
	//OrderCount    int64      `Gorm:"default:0;column:order_count"`                        // 订单数量
	//LastLoginTime *time.Time `Gorm:"column:last_login_time"`                              // 最后登录时间
}

func (DriverUser) TableName() string {
	return "driver_users"
}

//
//func (u *DriverUser) FindDriverById(db *gorm.DB, id int64) error {
//	return db.Debug().Where("id = ?", id).Limit(1).Find(u).Error
//}
//
//func (u *DriverUser) UpdateData(db *gorm.DB) error {
//	return db.Debug().Updates(u).Error
//}
//
//func (u *DriverUser) UpdateDataById(db *gorm.DB) error {
//	return db.Debug().Updates(u).Error
//}

func (u *DriverUser) CreateData(db *gorm.DB) error {
	return db.Debug().Create(u).Error
}
