package structs

type TransactionResponse struct {
	Id                 uint                        `json:"id"`
	OrderNumber        string                      `json:"order_number"`
	SubTotal           uint                        `json:"sub_total"`
	TotalAmount        uint                        `json:"total_amount"`
	PaymentStatus      string                      `json:"payment_status"`
	PaymentMethod      string                      `json:"payment_method"`
	BuyerName          string                      `json:"buyer_name"`
	Phone              string                      `json:"phone"`
	Notes              string                      `json:"notes"`
	PaidAt             *string                     `json:"paid_at"`
	ExpiredAt          *string                     `json:"expired_at"`
	CreatedAt          string                      `json:"created_at"`
	UpdatedAt          string                      `json:"updated_at"`
	TransactionDetails []TransactionDetailResponse `json:"transaction_details,omitempty"`
}

type TransactionCreateRequest struct {
	BuyerName string                           `json:"buyer_name" binding:"required"`
	Phone     string                           `json:"phone" binding:"required"`
	Notes     string                           `json:"notes"`
	Items     []TransactionDetailCreateRequest `json:"items" binding:"required,dive"`
}
