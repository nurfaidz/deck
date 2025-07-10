package models

type Transaction struct {
	GormModel
	OrderNumber        string              `json:"order_number" gorm:"not null;unique"`
	SubTotal           uint                `json:"sub_total" gorm:"not null"`
	TotalAmount        uint                `json:"total_amount" gorm:"not null"`
	PaymentStatus      string              `json:"payment_status" gorm:"not null;default:pending"`
	BuyerName          string              `json:"buyer_name" gorm:"not null"`
	Phone              string              `json:"phone" gorm:"not null"`
	TransactionDetails []TransactionDetail `json:"transaction_details" gorm:"foreignKey:TransactionId;references:Id"`
}
