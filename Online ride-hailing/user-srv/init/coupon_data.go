package init

import (
	"fmt"
	"time"
	"user-srv/model"

	"user-srv/global"
)

// InitCouponData 初始化优惠券测试数据
func InitCouponData() {
	initCouponTemplates()
	initUserCoupons()
	initCouponGrantTasks()
}

// initCouponTemplates 初始化优惠券模板数据
func initCouponTemplates() {
	var count int64
	global.DB.Model(&model.CouponTemplate{}).Count(&count)
	if count > 0 {
		return
	}

	templates := []model.CouponTemplate{
		{
			Name:         "新人专享券",
			Type:         4, // 新人券
			MinAmount:    10,
			ReduceAmount: 10,
			ValidType:    2, // 领取后N天
			ValidDays:    7,
			Total:        10000,
			Received:     500,
			PerLimit:     1,
			CityCode:     "",
			CarType:      "",
			Status:       1,
		},
		{
			Name:         "满30减5",
			Type:         1, // 满减券
			MinAmount:    30,
			ReduceAmount: 5,
			ValidType:    2,
			ValidDays:    14,
			Total:        5000,
			Received:     1200,
			PerLimit:     3,
			CityCode:     "",
			CarType:      "",
			Status:       1,
		},
		{
			Name:         "满50减10",
			Type:         1, // 满减券
			MinAmount:    50,
			ReduceAmount: 10,
			ValidType:    2,
			ValidDays:    14,
			Total:        3000,
			Received:     800,
			PerLimit:     2,
			CityCode:     "",
			CarType:      "",
			Status:       1,
		},
		{
			Name:      "8折折扣券",
			Type:      2, // 折扣券
			Discount:  0.8,
			MinAmount: 20,
			MaxReduce: 20,
			ValidType: 2,
			ValidDays: 10,
			Total:     2000,
			Received:  500,
			PerLimit:  2,
			CityCode:  "",
			CarType:   "",
			Status:    1,
		},
		{
			Name:         "立减3元",
			Type:         3, // 立减券
			ReduceAmount: 3,
			ValidType:    2,
			ValidDays:    7,
			Total:        10000,
			Received:     3000,
			PerLimit:     5,
			CityCode:     "",
			CarType:      "",
			Status:       1,
		},
		{
			Name:         "周末特惠券",
			Type:         1, // 满减券
			MinAmount:    40,
			ReduceAmount: 8,
			ValidType:    1, // 固定日期
			ValidStart:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			ValidEnd:     time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local),
			Total:        5000,
			Received:     2000,
			PerLimit:     2,
			CityCode:     "",
			UseTime:      "周六,周日",
			CarType:      "",
			Status:       1,
		},
		{
			Name:         "夜间出行券",
			Type:         1, // 满减券
			MinAmount:    25,
			ReduceAmount: 5,
			ValidType:    2,
			ValidDays:    30,
			Total:        3000,
			Received:     800,
			PerLimit:     4,
			CityCode:     "",
			UseTime:      "22:00-06:00",
			CarType:      "",
			Status:       1,
		},
		{
			Name:         "商务车专享券",
			Type:         1, // 满减券
			MinAmount:    100,
			ReduceAmount: 20,
			ValidType:    2,
			ValidDays:    14,
			Total:        1000,
			Received:     200,
			PerLimit:     1,
			CityCode:     "",
			CarType:      "商务型",
			Status:       1,
		},
	}

	err := global.DB.Create(&templates).Error
	if err != nil {
		panic("初始化优惠券模板数据失败: " + err.Error())
	}
}

// initUserCoupons 初始化用户优惠券数据
func initUserCoupons() {
	var count int64
	global.DB.Model(&model.UserCoupon{}).Count(&count)
	if count > 0 {
		return
	}

	now := time.Now()
	usedTime := now.Add(-24 * time.Hour)
	userCoupons := []model.UserCoupon{
		{
			UserID:     1,
			TemplateID: 1,
			CouponNo:   "CP" + generateCouponNo(),
			Status:     1,
			UsedTime:   nil,
			StartTime:  now,
			EndTime:    now.Add(7 * 24 * time.Hour),
		},
		{
			UserID:     1,
			TemplateID: 2,
			CouponNo:   "CP" + generateCouponNo(),
			Status:     1,
			UsedTime:   nil,
			StartTime:  now,
			EndTime:    now.Add(14 * 24 * time.Hour),
		},
		{
			UserID:     1,
			TemplateID: 3,
			CouponNo:   "CP" + generateCouponNo(),
			Status:     1,
			UsedTime:   nil,
			StartTime:  now,
			EndTime:    now.Add(14 * 24 * time.Hour),
		},
		{
			UserID:     2,
			TemplateID: 1,
			CouponNo:   "CP" + generateCouponNo(),
			Status:     1,
			UsedTime:   nil,
			StartTime:  now,
			EndTime:    now.Add(7 * 24 * time.Hour),
		},
		{
			UserID:     2,
			TemplateID: 4,
			CouponNo:   "CP" + generateCouponNo(),
			Status:     1,
			UsedTime:   nil,
			StartTime:  now,
			EndTime:    now.Add(10 * 24 * time.Hour),
		},
		{
			UserID:     2,
			TemplateID: 5,
			CouponNo:   "CP" + generateCouponNo(),
			Status:     2,
			UsedTime:   &usedTime,
			OrderNo:    "ORD202401010001",
			StartTime:  now.Add(-10 * 24 * time.Hour),
			EndTime:    now.Add(-3 * 24 * time.Hour),
		},
		{
			UserID:     3,
			TemplateID: 6,
			CouponNo:   "CP" + generateCouponNo(),
			Status:     1,
			UsedTime:   nil,
			StartTime:  now,
			EndTime:    now.Add(30 * 24 * time.Hour),
		},
		{
			UserID:     3,
			TemplateID: 7,
			CouponNo:   "CP" + generateCouponNo(),
			Status:     3,
			UsedTime:   nil,
			StartTime:  now.Add(-40 * 24 * time.Hour),
			EndTime:    now.Add(-10 * 24 * time.Hour),
		},
	}

	err := global.DB.Create(&userCoupons).Error
	if err != nil {
		panic("初始化用户优惠券数据失败: " + err.Error())
	}
}

// initCouponGrantTasks 初始化优惠券发放任务数据
func initCouponGrantTasks() {
	var count int64
	global.DB.Model(&model.CouponGrantTask{}).Count(&count)
	if count > 0 {
		return
	}

	tasks := []model.CouponGrantTask{
		{
			TemplateID: 1,
			GrantType:  1, // 新人发放
			GrantNum:   1,
			Status:     1,
		},
		{
			TemplateID: 5,
			GrantType:  3, // 活动发放
			GrantNum:   2,
			Status:     1,
		},
	}

	err := global.DB.Create(&tasks).Error
	if err != nil {
		panic("初始化优惠券发放任务数据失败: " + err.Error())
	}
}

// generateCouponNo 生成优惠券编号
func generateCouponNo() string {
	return fmt.Sprintf("%d%08d", time.Now().Unix(), int(time.Now().UnixNano())%100000000)
}
