package paymentsmodule

import (
	"context"
	"paytm-project/internal/models"

	"gorm.io/gorm"
)

type IPaymentsRepo interface {
	UpdatePaymentTransactional(tx *gorm.DB, paymentId string, status PaymentStatus) error
	CreatePayment(ctx context.Context, payment *Payment) error
	UpdatePaymentStatusById(ctx context.Context, paymentId string, status PaymentStatus) error
	GetPaymentById(ctx context.Context, paymentId string) (*Payment, error)
	GetDb() *gorm.DB
}

type PaymentsRepo struct {
	models.BaseRepo
}

func (pr *PaymentsRepo) GetPaymentById(ctx context.Context, paymentId string) (*Payment, error) {
	var payment Payment
	if err := pr.Db.Where("id = ?", paymentId).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (pr *PaymentsRepo) UpdatePaymentTransactional(tx *gorm.DB, paymentId string, status PaymentStatus) error {
	return tx.Model(&Payment{}).Where("id = ?", paymentId).Update("status", status).Error
}

func (pr *PaymentsRepo) CreatePayment(ctx context.Context, payment *Payment) error {
	return pr.Db.Create(payment).Error
}

func (pr *PaymentsRepo) UpdatePaymentStatusById(ctx context.Context, paymentId string, status PaymentStatus) error {
	return pr.Db.Model(&Payment{}).Where("id = ?", paymentId).Update("status", status).Error
}

func (pr *PaymentsRepo) GetDb() *gorm.DB {
	return pr.Db
}
