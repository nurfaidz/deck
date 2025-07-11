package controllers

import (
	"deck/helpers"
	"deck/models"
	"deck/services"
	"deck/structs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type PaymentController struct {
	db              *gorm.DB
	midtransService *services.MidtransService
}

func NewPaymentController(db *gorm.DB, midtransService *services.MidtransService) *PaymentController {
	return &PaymentController{
		db:              db,
		midtransService: midtransService,
	}
}

func (pc *PaymentController) CreatePayment(c *gin.Context) {
	var req structs.MidtransPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	var transaction models.Transaction
	if err := pc.db.Preload("TransactionDetails").Where("order_number = ?", req.OrderNumber).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Transaction not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	if transaction.PaymentStatus == models.PaymentStatusPaid {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Transaction already paid",
			Errors:  nil,
		})

		return
	}

	paymentResp, err := pc.midtransService.CreatePayment(&transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create payment",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	transaction.MidtransToken = paymentResp.Token
	transaction.MidtransOrderID = transaction.OrderNumber
	transaction.ExpiredAt = &time.Time{}
	*transaction.ExpiredAt = time.Now().Add(30 * time.Minute)

	if err := pc.db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update transaction",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Payment created successfully",
		Data:    paymentResp,
	})
}

func (pc *PaymentController) MidtransNotification(c *gin.Context) {
	var notification structs.MidtransNotification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Invalid notification data",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	if !pc.midtransService.VerifySignature(&notification) {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Invalid signature",
			Errors:  nil,
		})

		return
	}

	var transaction models.Transaction
	if err := pc.db.Where("order_number = ?", notification.OrderID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Transaction not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	switch notification.TransactionStatus {
	case "capture", "settlement":
		transaction.PaymentStatus = models.PaymentStatusPaid
		now := time.Now()
		transaction.PaidAt = &now
	case "pending":
		transaction.PaymentStatus = models.PaymentStatusPending
	case "deny", "cancel":
		transaction.PaymentStatus = models.PaymentStatusCancelled
	case "expire":
		transaction.PaymentStatus = models.PaymentStatusExpired
	case "failure":
		transaction.PaymentStatus = models.PaymentStatusFailed
	}

	if err := pc.db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update transaction status",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Notification processed successfully",
	})
}

func (pc *PaymentController) GetPaymentStatus(c *gin.Context) {
	orderNumber := c.Param("order_number")

	var transaction models.Transaction
	if err := pc.db.Where("order_number = ?", orderNumber).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Transaction not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})

		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Transaction status retrieved successfully",
		Data: map[string]interface{}{
			"order_number":   transaction.OrderNumber,
			"payment_status": transaction.PaymentStatus,
			"payment_method": transaction.PaymentMethod,
			"paid_at":        transaction.PaidAt,
			"expired_at":     transaction.ExpiredAt,
		},
	})
}
