package models

import "time"

type Transaction struct {
	GormModel
	OrderNumber        string              `json:"order_number" gorm:"not null;unique"`
	SubTotal           uint                `json:"sub_total" gorm:"not null"`
	TotalAmount        uint                `json:"total_amount" gorm:"not null"`
	PaymentStatus      string              `json:"payment_status" gorm:"not null;default:pending"`
	PaymentMethod      string              `json:"payment_method" gorm:"default:midtrans"`
	MidtransToken      string              `json:"midtrans_token,omitempty" gorm:"column:midtrans_token"`
	MidtransOrderID    string              `json:"midtrans_order_id,omitempty" gorm:"column:midtrans_order_id"`
	BuyerName          string              `json:"buyer_name" gorm:"not null"`
	Phone              string              `json:"phone" gorm:"not null"`
	PaidAt             *time.Time          `json:"paid_at"`
	ExpiredAt          *time.Time          `json:"expired_at"`
	TransactionDetails []TransactionDetail `json:"transaction_details" gorm:"foreignKey:TransactionId;references:Id"`
}

const (
	PaymentStatusPending   = "pending"
	PaymentStatusPaid      = "paid"
	PaymentStatusFailed    = "failed"
	PaymentStatusExpired   = "expired"
	PaymentStatusCancelled = "cancelled"
)
