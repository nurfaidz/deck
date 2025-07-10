package structs

type TransactionDetailResponse struct {
	Id            uint   `json:"id"`
	TransactionId uint   `json:"transaction_id"`
	ProductId     uint   `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      uint   `json:"quantity"`
	Price         uint   `json:"price"`
	TotalPrice    uint   `json:"total_price"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type TransactionDetailCreateRequest struct {
	TransactionId uint   `json:"transaction_id" binding:"required"`
	ProductId     uint   `json:"product_id" binding:"required"`
	ProductName   string `json:"product_name" binding:"required"`
	Quantity      uint   `json:"quantity" binding:"required"`
	Price         uint   `json:"price" binding:"required"`
	TotalPrice    uint   `json:"total_price" binding:"required"`
}
