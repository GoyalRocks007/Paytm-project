package modules

import (
	authmodule "paytm-project/internal/modules/auth_module"
	paymentsmodule "paytm-project/internal/modules/payments_module"

	"gorm.io/gorm"
)

var (
	registry *Registry
)

type Registry struct {
	AuthModule     authmodule.IAuthModule
	PaymentsModule paymentsmodule.IPaymentsModule
}

func GetRegistry() *Registry {
	if registry == nil {
		registry = &Registry{}
	}
	return registry
}

func (r *Registry) WithAuthModule(db *gorm.DB) *Registry {
	r.AuthModule = authmodule.GetAuthModule(db)
	return r
}

func (r *Registry) WithPaymentsModule(db *gorm.DB) *Registry {
	r.PaymentsModule = paymentsmodule.GetPaymentsModule(db, r.AuthModule)
	return r
}
