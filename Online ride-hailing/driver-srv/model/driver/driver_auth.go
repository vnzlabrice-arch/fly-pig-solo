package driver

import (
	"gorm.io/gorm"
)

// DriverAuth 司机认证审核
type DriverAuth struct {
	gorm.Model
	DriverID    int64  `Gorm:"not null;column:driver_id"`              // 司机ID
	RealName    string `Gorm:"size:32;column:real_name"`               // 真实姓名
	IDCard      string `gorm:"size:18;not null;column:id_card"`        // 身份证号
	IDCardFront string `gorm:"size:255;not null;column:id_card_front"` // 身份证正面
	IDCardBack  string `gorm:"size:255;not null;column:id_card_back"`  // 身份证反面
	LicenseImg  string `gorm:"size:255;not null;column:license_img"`   // 驾驶证照片
	FaceImg     string `gorm:"size:255;not null;column:face_img"`      // 人脸核验照
	AuditStatus int8   `Gorm:"default:0;column:audit_status"`          // 审核状态 1-待审核 2-通过 3-驳回
	Reason      string `Gorm:"varchar(255);column:reason"`             // 拒绝原因
	//AuthType    int8      `Gorm:"not null;column:auth_type"`        // 认证类型 1-实名认证 2-驾照认证
	//PicURL      string    `Gorm:"size:255;column:pic_url"`          // 认证图片URL
	//AuditTime   time.Time `Gorm:"column:audit_time"`                // 审核时间

}

func (DriverAuth) TableName() string {
	return "driver_auths"
}

func (a *DriverAuth) CreateData(db *gorm.DB) error {
	return db.Debug().Create(a).Error
}

func (a *DriverAuth) UpdateData(db *gorm.DB) error {
	return db.Debug().Updates(a).Error
}

func (a *DriverAuth) CreateAuth(db *gorm.DB) error {
	return db.Debug().Create(a).Error
}
