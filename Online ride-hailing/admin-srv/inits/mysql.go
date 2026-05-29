package inits

import (
	"admin-srv/global"
	"admin-srv/model/aftersale"
	"admin-srv/model/order"
	"admin-srv/model/system"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlInit() {
	data := global.Config.Mysql
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
		// 系统管理模块
		&system.AdminUser{},
		&system.AdminRole{},
		&system.AdminMenu{},
		&system.AdminRoleMenu{},
		&system.AdminLoginLog{},
		&system.AdminOperationLog{},
		&system.CityConfig{},
		&system.SystemConfig{},
		&system.SystemDict{},
		&system.CarTypeConfig{},
		&system.ArticleCategory{},
		&system.Article{},
		&system.AppUpdate{},
		&system.RiskConfig{},
		// 售后模块
		&aftersale.AftersaleOrder{},
		&aftersale.AftersaleAuditLog{},
		// 订单模块
		&order.InvoiceRecord{},
		&order.SettleRecord{},
		&order.BadDebtOrder{},
		&order.PaymentOrder{},
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
