package services

import (
	"crypto/sha512"
	"deck/models"
	"deck/structs"
	"fmt"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"os"
	"strings"
)

type MidtransService struct {
	client snap.Client
}

func NewMidtransService() *MidtransService {
	var client snap.Client

	if os.Getenv("MIDTRANS_PRODUCTION") == "true" {
		client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Production)
	} else {
		client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	}

	return &MidtransService{
		client: client,
	}
}

func (s *MidtransService) CreatePayment(transaction *models.Transaction) (*structs.MidtransPaymentResponse, error) {
	var itemDetails []midtrans.ItemDetails
	for _, detail := range transaction.TransactionDetails {
		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:    fmt.Sprintf("%d", detail.ProductId),
			Name:  detail.ProductName,
			Price: int64(detail.Price),
			Qty:   int32(detail.Quantity),
		})
	}

	// Customer details
	customerDetails := midtrans.CustomerDetails{
		FName: transaction.BuyerName,
		Phone: transaction.Phone,
	}

	// Create snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.OrderNumber,
			GrossAmt: int64(transaction.TotalAmount),
		},
		CustomerDetail: &customerDetails,
		Items:          &itemDetails,
		Expiry: &snap.ExpiryDetails{
			Duration: 30,
			Unit:     "minutes",
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeGopay,
			snap.PaymentTypeShopeepay,
			snap.PaymentTypeBankTransfer,
		},
	}

	// Create transaction
	snapResp, err := s.client.CreateTransaction(snapReq)
	if err != nil {
		return nil, err
	}

	return &structs.MidtransPaymentResponse{
		Token:       snapResp.Token,
		RedirectUrl: snapResp.RedirectURL,
	}, nil
}

func (s *MidtransService) VerifySignature(notification *structs.MidtransNotification) bool {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	signatureString := fmt.Sprintf("%s%s%s%s",
		notification.OrderID,
		notification.StatusCode,
		notification.GrossAmount,
		serverKey)

	hash := sha512.New()
	hash.Write([]byte(signatureString))
	calculatedSignature := fmt.Sprintf("%x", hash.Sum(nil))

	return strings.EqualFold(calculatedSignature, notification.SignatureKey)
}

func (s *MidtransService) GetPaymentStatus(orderID string) (string, error) {
	return "pending", nil
}
