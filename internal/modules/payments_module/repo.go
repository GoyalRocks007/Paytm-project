package paymentsmodule

import (
	"paytm-project/internal/models"

	"gorm.io/gorm"
)

type IPaymentsRepo interface {
	CreatePaymentTransactional(tx *gorm.DB, payment *Payment) error
	GetDb() *gorm.DB
}

type PaymentsRepo struct {
	models.BaseRepo
}

func (pr *PaymentsRepo) CreatePaymentTransactional(tx *gorm.DB, payment *Payment) error {
	return tx.Create(payment).Error
}

func (pr *PaymentsRepo) GetDb() *gorm.DB {
	return pr.Db
}
