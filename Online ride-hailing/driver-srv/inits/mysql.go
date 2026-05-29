package init

import (
	"driver-srv/global"
	user3 "driver-srv/model/driver"
	"fmt"
	"sync"
	"time"

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
		&user3.DriverLocation{},
		&user3.DriverAuth{},
		&user3.DriverFeedback{},
		&user3.DriverIncome{},
		&user3.DriverMessage{},
		&user3.DriverUser{},
		&user3.DriverViolation{},
		&user3.DriverWallet{},
		&user3.DriverWithdraw{},
		&user3.DriverCar{},
		&user3.OrderComment{},
		&user3.SysRule{},
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
