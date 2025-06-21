package paymentsmodule

import (
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"

	"gorm.io/gorm"
)

var (
	paymentsModuleInstance IPaymentsModule
	GetPaymentsModule      = func(db *gorm.DB, authModule authmodule.IAuthModule) IPaymentsModule {
		if paymentsModuleInstance != nil {
			return paymentsModuleInstance
		}
		return initPaymentsModule(db, authModule)
	}
	initPaymentsModule = func(db *gorm.DB, authModule authmodule.IAuthModule) IPaymentsModule {
		paymentsModuleInstance = &PaymentsModule{
			PaymentsCore: &PaymentsCore{
				AuthModule: authModule,
				PaymentsRepo: &PaymentsRepo{
					BaseRepo: models.BaseRepo{
						Db: db,
					},
				},
			},
		}
		return paymentsModuleInstance
	}
)

type IPaymentsModule interface {
	GetCore() IPaymentsCore
}

type PaymentsModule struct {
	PaymentsCore IPaymentsCore
}

func (a *PaymentsModule) GetCore() IPaymentsCore {
	return a.PaymentsCore
}
