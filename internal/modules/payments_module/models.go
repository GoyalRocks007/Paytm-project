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

	SUCCESS PaymentStatus = "SUCCESS"
	FAIL    PaymentStatus = "FAIL"
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

type CreatePaymentRequestDto struct {
	Sender   string
	Receiver string
	Amount   float64
	Mode     PaymentMedium
}

type CreatePaymentResponseDto struct {
	models.BaseSuccessResponse
}
