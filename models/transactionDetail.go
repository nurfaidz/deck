package models

type TransactionDetail struct {
	GormModel
	TransactionId uint        `json:"transaction_id" gorm:"not null"`
	ProductId     uint        `json:"product_id" gorm:"not null"`
	ProductName   string      `json:"product_name" gorm:"not null"`
	Quantity      uint        `json:"quantity" gorm:"not null"`
	Price         uint        `json:"price" gorm:"not null"`
	TotalPrice    uint        `json:"total_price" gorm:"not null"`
	Notes         string      `json:"notes"`
	Transaction   Transaction `json:"transaction" gorm:"foreignKey:TransactionId;references:Id"`
	Product       Product     `json:"product" gorm:"foreignKey:ProductId;references:Id"`
}
