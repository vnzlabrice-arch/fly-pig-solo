package driver

import "time"

// PassengerOrder 订单主表
type PassengerOrder struct {
	OrderID              string     `Gorm:"primaryKey;size:32;column:order_id"`
	PassengerID          int64      `Gorm:"not null;column:passenger_id"`
	PassengerName        string     `Gorm:"size:50;column:passenger_name"`
	PassengerPhone       string     `Gorm:"size:20;column:passenger_phone"`
	CouponID             int64      `Gorm:"column:coupon_id"`
	CouponName           string     `Gorm:"size:64;column:coupon_name"`
	DriverID             int64      `Gorm:"column:driver_id"`
	OrderType            int8       `Gorm:"not null;column:order_type"`
	CarType              string     `Gorm:"size:20;not null;column:car_type"`
	Status               int8       `Gorm:"not null;column:status"`
	StartLng             float64    `Gorm:"type:decimal(10,6);not null;column:start_lng"`
	StartLat             float64    `Gorm:"type:decimal(10,6);not null;column:start_lat"`
	StartAddress         string     `Gorm:"size:255;not null;column:start_address"`
	EndLng               float64    `Gorm:"type:decimal(10,6);not null;column:end_lng"`
	EndLat               float64    `Gorm:"type:decimal(10,6);not null;column:end_lat"`
	EndAddress           string     `Gorm:"size:255;not null;column:end_address"`
	PassRemark           string     `Gorm:"size:255;column:pass_remark"`
	DriverRemark         string     `Gorm:"size:255;column:driver_remark"`
	PassengerRiskLevel   int8       `Gorm:"default:0;column:passenger_risk_level"`
	EstimatedPrice       float64    `Gorm:"type:decimal(10,2);not null;column:estimated_price"`
	SurgePrice           float64    `Gorm:"type:decimal(10,2);default:0;column:surge_price"`
	FinalPrice           float64    `Gorm:"type:decimal(10,2);column:final_price"`
	PayStatus            int8       `Gorm:"default:0;column:pay_status"`
	PayType              string     `Gorm:"size:32;column:pay_type"`
	BookTime             time.Time  `Gorm:"not null;column:book_time"`
	AppointTime          *time.Time `Gorm:"column:appoint_time"`
	EstimatedArrivalTime *time.Time `Gorm:"column:estimated_arrival_time"`
	PickupTime           *time.Time `Gorm:"column:pickup_time"`
	StartTime            *time.Time `Gorm:"column:start_time"`
	EndTime              *time.Time `Gorm:"column:end_time"`
	PayTime              *time.Time `Gorm:"column:pay_time"`
	CouponDeduction      float64    `Gorm:"type:decimal(10,2);default:0;column:coupon_deduction"`
	ComplaintStatus      int8       `Gorm:"default:0;column:complaint_status"`
	PassengerTags        string     `Gorm:"size:255;column:passenger_tags"`
	ArrivalAccuracy      float64    `Gorm:"type:decimal(10,2);default:0;column:arrival_accuracy"`
	ArrivalPhoto         string     `Gorm:"size:255;column:arrival_photo"`
	PickupPhoto          string     `Gorm:"size:255;column:pickup_photo"`
	StartOdometer        float64    `Gorm:"type:decimal(10,2);default:0;column:start_odometer"`
	PassengerConfirmCode string     `Gorm:"size:16;column:passenger_confirm_code"`
	CancelFee            float64    `Gorm:"type:decimal(10,2);default:0;column:cancel_fee"`
	CancelResponsibility string     `Gorm:"size:32;column:cancel_responsibility"`
	EvidencePhoto        string     `Gorm:"size:255;column:evidence_photo"`
	AppealSupport        bool       `Gorm:"default:false;column:appeal_support"`
	DriverLng            float64    `Gorm:"type:decimal(10,6);default:0;column:driver_lng"`
	DriverLat            float64    `Gorm:"type:decimal(10,6);default:0;column:driver_lat"`
	CancelReason         string     `Gorm:"size:255;column:cancel_reason"`
	CreatedAt            time.Time  `Gorm:"autoCreateTime;column:created_at"`
	UpdatedAt            time.Time  `Gorm:"autoUpdateTime;column:updated_at"`
}

func (PassengerOrder) TableName() string {
	return "passenger_orders"
}
