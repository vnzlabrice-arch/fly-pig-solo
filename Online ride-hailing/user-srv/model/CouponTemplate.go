package model

import "time"

type CouponTemplate struct {
	ID           int64     `Gorm:"primaryKey;autoIncrement;column:id"`
	Name         string    `Gorm:"size:64;not null;column:name"`                      // 优惠券名称
	Type         int8      `Gorm:"not null;column:type"`                              // 券类型 1-满减 2-折扣 3-立减 4-新人券
	Discount     float64   `Gorm:"type:decimal(5,2);column:discount"`                 // 折扣（仅折扣券）
	MinAmount    float64   `Gorm:"type:decimal(10,2);default:0;column:min_amount"`    // 使用门槛
	ReduceAmount float64   `Gorm:"type:decimal(10,2);default:0;column:reduce_amount"` // 减免金额
	MaxReduce    float64   `Gorm:"type:decimal(10,2);column:max_reduce"`              // 折扣券最高减免
	ValidType    int8      `Gorm:"not null;column:valid_type"`                        // 有效期类型 1-固定日期 2-领取后N天
	ValidStart   time.Time `Gorm:"column:valid_start"`                                // 固定有效期开始
	ValidEnd     time.Time `Gorm:"column:valid_end"`                                  // 固定有效期结束
	ValidDays    int       `Gorm:"column:valid_days"`                                 // 领取后有效天数
	Total        int       `Gorm:"default:0;column:total"`                            // 发放总数量 0=不限
	Received     int       `Gorm:"default:0;column:received"`                         // 已领取数量
	PerLimit     int       `Gorm:"default:1;column:per_limit"`                        // 每人限领张数
	CityCode     string    `Gorm:"size:32;column:city_code"`                          // 可用城市编码
	StartRegion  string    `Gorm:"size:128;column:start_region"`                      // 可用起点区域
	EndRegion    string    `Gorm:"size:128;column:end_region"`                        // 可用终点区域
	CarType      string    `Gorm:"size:32;column:car_type"`                           // 可用车型
	UseTime      string    `Gorm:"size:255;column:use_time"`                          // 可用时段
	Status       int8      `Gorm:"default:1;column:status"`                           // 状态 1-正常 2-下架 3-已领完
	CreatedAt    time.Time `Gorm:"autoCreateTime;column:created_at"`
	UpdatedAt    time.Time `Gorm:"autoUpdateTime;column:updated_at"`
}
