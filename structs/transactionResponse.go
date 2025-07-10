package structs

type TransactionResponse struct {
	Id            uint   `json:"id"`
	OrderNumber   string `json:"order_number"`
	SubTotal      uint   `json:"sub_total"`
	TotalAmount   uint   `json:"total_amount"`
	PaymentStatus string `json:"payment_status"`
	BuyerName     string `json:"buyer_name"`
	Phone         string `json:"phone"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	//TransactionDetails []TransactionDetail `json:"transaction_details"`
}

type TransactionCreateRequest struct {
	BuyerName string `json:"buyer_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}
