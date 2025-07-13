package controllers

import (
	"deck/helpers"
	"deck/models"
	"deck/services"
	"deck/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionController struct {
	db                 *gorm.DB
	transactionService *services.TransactionService
}

func NewTransactionController(db *gorm.DB, transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{
		db:                 db,
		transactionService: transactionService,
	}
}

// CreateTransaction - Create new transaction
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var req structs.TransactionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Create transaction
	transaction, err := tc.transactionService.CreateTransaction(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create transaction",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Convert to response
	response := tc.toTransactionResponse(transaction)

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Transaction created successfully",
		Data:    response,
	})
}

// GetTransaction - Get transaction by order number
func (tc *TransactionController) GetTransaction(c *gin.Context) {
	orderNumber := c.Param("order_number")

	transaction, err := tc.transactionService.GetTransaction(orderNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Transaction not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	response := tc.toTransactionResponse(transaction)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Transaction retrieved successfully",
		Data:    response,
	})
}

// GetTransactionByID - Get transaction by ID
func (tc *TransactionController) GetTransactionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Invalid transaction ID",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	transaction, err := tc.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Transaction not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	response := tc.toTransactionResponse(transaction)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Transaction retrieved successfully",
		Data:    response,
	})
}

// GetAllTransactions - Get all transactions (untuk admin)
func (tc *TransactionController) GetAllTransactions(c *gin.Context) {
	transactions, err := tc.transactionService.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve transactions",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var responses []structs.TransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, *tc.toTransactionResponse(&transaction))
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Transactions retrieved successfully",
		Data:    responses,
	})
}

// Helper method to convert model to response
func (tc *TransactionController) toTransactionResponse(transaction *models.Transaction) *structs.TransactionResponse {
	var details []structs.TransactionDetailResponse
	for _, detail := range transaction.TransactionDetails {
		details = append(details, structs.TransactionDetailResponse{
			Id:            detail.Id,
			TransactionId: detail.TransactionId,
			ProductId:     detail.ProductId,
			ProductName:   detail.ProductName,
			Quantity:      detail.Quantity,
			Price:         detail.Price,
			TotalPrice:    detail.TotalPrice,
			Notes:         detail.Notes,
			CreatedAt:     detail.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     detail.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	var paidAt, expiredAt *string
	if transaction.PaidAt != nil {
		paidAtStr := transaction.PaidAt.Format("2006-01-02 15:04:05")
		paidAt = &paidAtStr
	}
	if transaction.ExpiredAt != nil {
		expiredAtStr := transaction.ExpiredAt.Format("2006-01-02 15:04:05")
		expiredAt = &expiredAtStr
	}

	return &structs.TransactionResponse{
		Id:                 transaction.Id,
		OrderNumber:        transaction.OrderNumber,
		SubTotal:           transaction.SubTotal,
		TotalAmount:        transaction.TotalAmount,
		PaymentStatus:      transaction.PaymentStatus,
		PaymentMethod:      transaction.PaymentMethod,
		BuyerName:          transaction.BuyerName,
		Phone:              transaction.Phone,
		PaidAt:             paidAt,
		ExpiredAt:          expiredAt,
		CreatedAt:          transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          transaction.UpdatedAt.Format("2006-01-02 15:04:05"),
		TransactionDetails: details,
	}
}
