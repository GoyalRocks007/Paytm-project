package paymentsmodule

import (
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"
)

// enums
type PaymentMedium string
type PaymentStatus string

const (
	NEFT PaymentMedium = "NEFT"
	UPI  PaymentMedium = "UPI"

	PaymentStatusSuccess    PaymentStatus = "SUCCESS"
	PaymentStatusInitiated  PaymentStatus = "INITIATED"
	PaymentStatusAuthorized PaymentStatus = "AUTHORIZED"
	PaymentStatusFailed     PaymentStatus = "FAIL"
)

type Payment struct {
	models.BaseModel
	Amount     float64
	SenderID   string
	ReceiverID string
	Sender     authmodule.User `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Receiver   authmodule.User `gorm:"foreignKey:ReceiverID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Mode       PaymentMedium
	Status     PaymentStatus
}

type InitiatePaymentRequestDto struct {
	Sender   string
	Receiver string
	Amount   float64
	Mode     PaymentMedium
}

type InitiatePaymentResponseDto struct {
	models.BaseSuccessResponse
	PaymentId string `json:"paymentId"`
}

type ExecutePaymentResponseDto struct {
	models.BaseSuccessResponse
}

type ExecutePaymentRequest struct {
	PaymentId string
	Sender    string
	Receiver  string
	Amount    float64
	Mode      PaymentMedium
}
