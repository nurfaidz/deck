package services

import (
	"deck/models"
	"deck/structs"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		db: db,
	}
}

// CreateTransaction
func (ts *TransactionService) CreateTransaction(req *structs.TransactionCreateRequest) (*models.Transaction, error) {
	tx := ts.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate order number
	orderNumber := ts.generateOrderNumber()

	transaction := models.Transaction{
		OrderNumber:   orderNumber,
		BuyerName:     req.BuyerName,
		Phone:         req.Phone,
		PaymentStatus: models.PaymentStatusPending,
		PaymentMethod: "midtrans",
	}

	// Calculate totals dan create transaction details
	var subTotal uint = 0
	var transactionDetails []models.TransactionDetail

	for _, item := range req.Items {
		var product models.Product
		if err := tx.Where("id = ? AND is_available = ?", item.ProductId, true).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("product not found or not available: %d", item.ProductId)
		}

		// Calculate item total
		itemTotal := product.Price * item.Quantity
		subTotal += itemTotal

		detail := models.TransactionDetail{
			ProductId:   product.Id,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Price:       product.Price,
			TotalPrice:  itemTotal,
			Notes:       item.Notes,
		}
		transactionDetails = append(transactionDetails, detail)
	}

	totalAmount := subTotal

	transaction.SubTotal = subTotal
	transaction.TotalAmount = totalAmount

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range transactionDetails {
		transactionDetails[i].TransactionId = transaction.Id
		if err := tx.Create(&transactionDetails[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load transaction details untuk response
	ts.db.Preload("TransactionDetails").First(&transaction, transaction.Id)

	return &transaction, nil
}

// generateOrderNumber
func (ts *TransactionService) generateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("ORD-%s", timestamp)
}

// GetTransactionByID
func (ts *TransactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := ts.db.Preload("TransactionDetails").First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetAllTransactions
func (ts *TransactionService) GetAllTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := ts.db.Preload("TransactionDetails").Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// UpdateTransactionStatus
func (ts *TransactionService) UpdateTransactionStatus(orderNumber string, status string) error {
	var transaction models.Transaction
	if err := ts.db.Where("order_number = ?", orderNumber).First(&transaction).Error; err != nil {
		return err
	}

	transaction.PaymentStatus = status
	if status == models.PaymentStatusPaid {
		now := time.Now()
		transaction.PaidAt = &now
	}

	return ts.db.Save(&transaction).Error
}
