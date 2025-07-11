package structs

type TransactionDetailResponse struct {
	Id            uint   `json:"id"`
	TransactionId uint   `json:"transaction_id"`
	ProductId     uint   `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      uint   `json:"quantity"`
	Price         uint   `json:"price"`
	TotalPrice    uint   `json:"total_price"`
	Notes         string `json:"notes"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type TransactionDetailCreateRequest struct {
	ProductId uint   `json:"product_id" binding:"required"`
	Quantity  uint   `json:"quantity" binding:"required,min=1"`
	Notes     string `json:"notes"`
}
