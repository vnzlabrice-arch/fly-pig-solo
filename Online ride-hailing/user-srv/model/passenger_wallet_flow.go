package model

import "time"

// PassengerWalletFlow 钱包流水表 - 存储用户账户的资金变动记录
type PassengerWalletFlow struct {
	ID          int64     `Gorm:"primaryKey;autoIncrement;column:id"`         // 流水ID，主键，自增
	PassengerID int64     `Gorm:"not null;column:passenger_id"`               // 用户ID，关联用户表
	OrderID     string    `Gorm:"size:32;column:order_id"`                    // 关联订单ID（如有）
	FlowType    int8      `Gorm:"not null;column:flow_type"`                  // 流水类型：1-充值，2-消费，3-退款，4-提现，5-奖励，6-扣款
	Amount      float64   `Gorm:"type:decimal(10,2);not null;column:amount"`  // 变动金额（正数为收入，负数为支出）
	Balance     float64   `Gorm:"type:decimal(10,2);not null;column:balance"` // 变动后的账户余额
	PayChannel  string    `Gorm:"size:20;column:pay_channel"`                 // 支付渠道：如"支付宝"、"微信支付"、"银行卡"
	TradeNo     string    `Gorm:"size:64;column:trade_no"`                    // 第三方交易单号
	Remark      string    `Gorm:"size:255;column:remark"`                     // 备注说明
	CreatedAt   time.Time `Gorm:"autoCreateTime;column:created_at"`           // 创建时间（交易时间）
}
