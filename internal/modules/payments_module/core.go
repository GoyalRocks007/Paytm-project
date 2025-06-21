package paymentsmodule

import (
	"errors"
	"log"
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"

	"gorm.io/gorm"
)

type IPaymentsCore interface {
	CreatePayment(createPaymentRequest *CreatePaymentRequestDto) (*CreatePaymentResponseDto, error)
}

type PaymentsCore struct {
	AuthModule   authmodule.IAuthModule
	PaymentsRepo IPaymentsRepo
}

func (pc *PaymentsCore) CreatePayment(createPaymentRequest *CreatePaymentRequestDto) (*CreatePaymentResponseDto, error) {

	err := pc.TransactionalCreatePayment(createPaymentRequest)

	if err != nil {
		log.Println("sorry some problem occured while creating your payment")
		return nil, err
	}

	return &CreatePaymentResponseDto{
		BaseSuccessResponse: models.BaseSuccessResponse{
			Success: true,
		},
	}, nil

}
func (pc *PaymentsCore) TransactionalCreatePayment(createPaymentRequest *CreatePaymentRequestDto) error {
	return pc.PaymentsRepo.GetDb().Transaction(func(tx *gorm.DB) error {
		// Fetch and lock both wallets
		sender, err := pc.AuthModule.GetCore().GetRepo().GetUserByEmailTransactional(tx, createPaymentRequest.Sender)
		if err != nil {
			return err
		}
		receiver, err := pc.AuthModule.GetCore().GetRepo().GetUserByEmailTransactional(tx, createPaymentRequest.Receiver)
		if err != nil {
			return err
		}

		if sender.Wallet.Balance < createPaymentRequest.Amount {
			log.Printf("the user %s, doesn't have enough balance", sender.Id)
			return errors.New("insufficient Balance")
		}

		// Update balances
		if err := pc.AuthModule.GetCore().GetRepo().UpdateWalletBalance(tx, sender.Id, -createPaymentRequest.Amount); err != nil {
			return err
		}

		if err := pc.AuthModule.GetCore().GetRepo().UpdateWalletBalance(tx, receiver.Id, createPaymentRequest.Amount); err != nil {
			return err
		}

		// (Optional) create payment record
		if err := pc.PaymentsRepo.CreatePaymentTransactional(tx, &Payment{
			SenderID:   sender.Id,
			ReceiverID: receiver.Id,
			Amount:     createPaymentRequest.Amount,
			Mode:       createPaymentRequest.Mode,
		}); err != nil {
			return err
		}

		return nil
	})
}
