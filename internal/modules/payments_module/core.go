package paymentsmodule

import (
	"context"
	"errors"
	"log"
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"
	otpmodule "paytm-project/internal/modules/otp_module"

	"gorm.io/gorm"
)

type IPaymentsCore interface {
	InitiatePayment(ctx context.Context, initiatePaymentRequest *InitiatePaymentRequestDto) (*InitiatePaymentResponseDto, error)
	GenerateOtpForPayment(ctx context.Context, paymentId string) (*otpmodule.GenerateOtpResponse, error)
	AuthorizePayment(ctx context.Context, verifyOtpRequest *otpmodule.VerifyOtpRequest) (string, error)
	ExecutePayment(ctx context.Context, paymentId string) (*ExecutePaymentResponseDto, error)
}

type PaymentsCore struct {
	AuthModule   authmodule.IAuthModule
	OtpModule    otpmodule.IOtpModule
	PaymentsRepo IPaymentsRepo
}

func (pc *PaymentsCore) InitiatePayment(ctx context.Context, initiatePaymentRequest *InitiatePaymentRequestDto) (*InitiatePaymentResponseDto, error) {

	senderId, err := pc.AuthModule.GetCore().GetRepo().GetUserIdByEmail(initiatePaymentRequest.Sender)
	if err != nil {
		return nil, err
	}
	receiverId, err := pc.AuthModule.GetCore().GetRepo().GetUserIdByEmail(initiatePaymentRequest.Receiver)
	if err != nil {
		return nil, err
	}
	payment := &Payment{
		SenderID:   senderId,
		ReceiverID: receiverId,
		Amount:     initiatePaymentRequest.Amount,
		Mode:       initiatePaymentRequest.Mode,
		Status:     PaymentStatusInitiated,
	}
	if err := pc.PaymentsRepo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return &InitiatePaymentResponseDto{
		BaseSuccessResponse: models.BaseSuccessResponse{
			Success: true,
		},
		PaymentId: payment.Id,
	}, nil

}

// GenerateOtpForPayment generates an otp for a given payment id, using the email stored in the context.
// It is called by the PaymentsController, and the generated otp is sent to the user to verify the payment.
func (pc *PaymentsCore) GenerateOtpForPayment(ctx context.Context, paymentId string) (*otpmodule.GenerateOtpResponse, error) {
	email := ctx.Value("email").(string)
	generateOtpRequest := &otpmodule.GenerateOtpRequest{
		Receiver: email,
		Claims: map[string]interface{}{
			"paymentId": paymentId,
		},
	}
	res, err := pc.OtpModule.GetCore().GenerateOtp(ctx, generateOtpRequest)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pc *PaymentsCore) AuthorizePayment(ctx context.Context, verifyOtpRequest *otpmodule.VerifyOtpRequest) (string, error) {
	_, _, claims, err := pc.OtpModule.GetCore().VerifyOtp(ctx, verifyOtpRequest)
	if err != nil {
		return "", err
	}

	paymentId := claims["paymentId"].(string)
	if err := pc.PaymentsRepo.UpdatePaymentStatusById(ctx, paymentId, PaymentStatusAuthorized); err != nil {
		return "", err
	}
	return paymentId, nil
}

func (pc *PaymentsCore) ExecutePayment(ctx context.Context, paymentId string) (*ExecutePaymentResponseDto, error) {

	payment, err := pc.PaymentsRepo.GetPaymentById(ctx, paymentId)
	if err != nil {
		return nil, err
	}
	if payment.Status != PaymentStatusAuthorized {
		return nil, errors.New("payment is not authorized")
	}

	executePaymentRequest := &ExecutePaymentRequest{
		Sender:   payment.SenderID,
		Receiver: payment.ReceiverID,
		Amount:   payment.Amount,
		Mode:     payment.Mode,
	}

	err = pc.TransactionalCreatePayment(executePaymentRequest)

	if err != nil {
		return nil, err
	}

	return &ExecutePaymentResponseDto{
		BaseSuccessResponse: models.BaseSuccessResponse{
			Success: true,
		},
	}, nil

}
func (pc *PaymentsCore) TransactionalCreatePayment(executePaymentRequest *ExecutePaymentRequest) error {
	return pc.PaymentsRepo.GetDb().Transaction(func(tx *gorm.DB) error {
		// Fetch and lock both wallets
		senderWallet, err := pc.AuthModule.GetCore().GetRepo().GetWalletWithLock(tx, executePaymentRequest.Sender)
		if err != nil {
			return err
		}
		_, err = pc.AuthModule.GetCore().GetRepo().GetWalletWithLock(tx, executePaymentRequest.Receiver)
		if err != nil {
			return err
		}

		if senderWallet.Balance < executePaymentRequest.Amount {
			log.Printf("the user %s, doesn't have enough balance", executePaymentRequest.Sender)
			return errors.New("insufficient Balance")
		}

		// Update balances
		if err := pc.AuthModule.GetCore().GetRepo().UpdateWalletBalance(tx, executePaymentRequest.Sender, -executePaymentRequest.Amount); err != nil {
			return err
		}

		if err := pc.AuthModule.GetCore().GetRepo().UpdateWalletBalance(tx, executePaymentRequest.Receiver, executePaymentRequest.Amount); err != nil {
			return err
		}

		// (Optional) execute payment record
		if err := pc.PaymentsRepo.UpdatePaymentTransactional(tx, executePaymentRequest.PaymentId, PaymentStatusSuccess); err != nil {
			return err
		}

		return nil
	})
}
