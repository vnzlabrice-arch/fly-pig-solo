package model

import "time"

// OrderFeeDetail 订单费用明细表 - 存储订单的各项费用明细
type OrderFeeDetail struct {
	ID          int64     `Gorm:"primaryKey;autoIncrement;column:id"`        // 费用明细ID，主键，自增
	OrderID     string    `Gorm:"size:32;not null;column:order_id"`          // 关联订单ID
	FeeType     string    `Gorm:"size:20;not null;column:fee_type"`          // 费用类型：如"里程费"、"时长费"、"起步价"、"燃油附加费"、"优惠券抵扣"、"平台服务费"
	Amount      float64   `Gorm:"type:decimal(10,2);not null;column:amount"` // 费用金额（正数为收入，负数为减免）
	Description string    `Gorm:"size:255;column:description"`               // 费用描述
	CreatedAt   time.Time `Gorm:"autoCreateTime;column:created_at"`          // 创建时间
}
