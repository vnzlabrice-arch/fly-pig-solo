package init

import (
	"fmt"
	"sync"
	"time"
	"user-srv/global"
	"user-srv/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlInit() {
	data := global.ConfigData.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		data.User,
		data.Password,
		data.Host,
		data.Port,
		data.Database,
	)
	var err error
	once := sync.Once{}
	once.Do(func() {
		global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("数据库连接成功")
	err = global.DB.AutoMigrate(
		&model.PassengerUser{},
		&model.Service{},
		&model.PassengerOrder{},
		&model.OrderFeeDetail{},
		&model.PassengerAddressBook{},
		&model.PassengerMemberBenefit{},
		&model.PassengerTripSafetyLog{},
		&model.PassengerWalletFlow{},
		//优惠卷
		&model.CouponGrantTask{},
		&model.CouponTemplate{},
		&model.CouponUseLog{},
		&model.UserCoupon{},
		&model.UserCouponLimit{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("数据库迁移成功")
	sqlDB, err := global.DB.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	fmt.Println("连接池链接成功")
}
