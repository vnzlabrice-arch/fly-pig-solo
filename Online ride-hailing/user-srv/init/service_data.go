package init

import (
	"user-srv/global"
	"user-srv/model"
	"user-srv/pkg"
)

// InitServiceData 初始化服务数据
func InitServiceData() {
	var count int64
	global.DB.Model(&model.Service{}).Count(&count)
	if count > 0 {
		return
	}

	list := []model.Service{
		{Name: "经济型", Price: 12.50, CarType: 1, Status: 1, Sort: 1, Sales: 580},
		{Name: "舒适型", Price: 18.80, CarType: 2, Status: 1, Sort: 2, Sales: 420},
		{Name: "商务型", Price: 28.00, CarType: 3, Status: 1, Sort: 3, Sales: 260},
		{Name: "豪华型", Price: 45.00, CarType: 4, Status: 1, Sort: 4, Sales: 120},
		{Name: "拼车", Price: 9.90, CarType: 5, Status: 1, Sort: 5, Sales: 680},
		{Name: "顺风车", Price: 15.80, CarType: 6, Status: 1, Sort: 6, Sales: 310},
	}

	err := global.DB.Create(&list).Error
	if err != nil {
		panic("初始化服务数据失败: " + err.Error())
	}
	_ = pkg.ClearServiceListCache()
}
