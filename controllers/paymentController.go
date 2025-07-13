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

// CreatePayment - Membuat payment link dari Midtrans
// Endpoint ini digunakan setelah transaksi dibuat untuk generate payment link
func (pc *PaymentController) CreatePayment(c *gin.Context) {
	var req structs.MidtransPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Cari transaksi berdasarkan order number
	var transaction models.Transaction
	if err := pc.db.Preload("TransactionDetails").Where("order_number = ?", req.OrderNumber).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Transaction not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Check apakah transaksi sudah dibayar
	if transaction.PaymentStatus == models.PaymentStatusPaid {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Transaction already paid",
			Errors:  nil,
		})
		return
	}

	// Create payment via Midtrans
	paymentResp, err := pc.midtransService.CreatePayment(&transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create payment",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Update transaksi dengan token dan expired time
	transaction.MidtransToken = paymentResp.Token
	transaction.MidtransOrderID = transaction.OrderNumber
	expiredAt := time.Now().Add(30 * time.Minute)
	transaction.ExpiredAt = &expiredAt

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

// MidtransNotification - Webhook endpoint untuk notifikasi dari Midtrans
// Endpoint ini akan dipanggil otomatis oleh Midtrans ketika ada perubahan status pembayaran
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

	// Verify signature untuk keamanan
	if !pc.midtransService.VerifySignature(&notification) {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Invalid signature",
			Errors:  nil,
		})
		return
	}

	// Cari transaksi berdasarkan order ID
	var transaction models.Transaction
	if err := pc.db.Where("order_number = ?", notification.OrderID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Transaction not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Update status berdasarkan notifikasi Midtrans
	switch notification.TransactionStatus {
	case "capture", "settlement":
		transaction.PaymentStatus = models.PaymentStatusPaid
		now := time.Now()
		transaction.PaidAt = &now
		transaction.PaymentMethod = notification.PaymentType
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
		Data: map[string]interface{}{
			"order_number":   transaction.OrderNumber,
			"payment_status": transaction.PaymentStatus,
		},
	})
}

// GetPaymentStatus - Mendapatkan status pembayaran transaksi
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

	var paidAt, expiredAt *string
	if transaction.PaidAt != nil {
		paidAtStr := transaction.PaidAt.Format("2006-01-02 15:04:05")
		paidAt = &paidAtStr
	}
	if transaction.ExpiredAt != nil {
		expiredAtStr := transaction.ExpiredAt.Format("2006-01-02 15:04:05")
		expiredAt = &expiredAtStr
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Payment status retrieved successfully",
		Data: map[string]interface{}{
			"order_number":   transaction.OrderNumber,
			"payment_status": transaction.PaymentStatus,
			"payment_method": transaction.PaymentMethod,
			"paid_at":        paidAt,
			"expired_at":     expiredAt,
		},
	})
}
