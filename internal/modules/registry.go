package modules

import (
	"paytm-project/internal/clients"
	adminmodule "paytm-project/internal/modules/admin_module"
	authmodule "paytm-project/internal/modules/auth_module"
	otpmodule "paytm-project/internal/modules/otp_module"
	paymentsmodule "paytm-project/internal/modules/payments_module"

	"gorm.io/gorm"
)

var (
	registry *Registry
)

type Registry struct {
	AuthModule     authmodule.IAuthModule
	PaymentsModule paymentsmodule.IPaymentsModule
	AdminModule    adminmodule.IAdminModule
	OtpModule      otpmodule.IOtpModule
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
	r.PaymentsModule = paymentsmodule.GetPaymentsModule(db, r.AuthModule, r.OtpModule)
	return r
}

func (r *Registry) WithAdminModule(db *gorm.DB) *Registry {
	r.AdminModule = adminmodule.GetAdminModule(db, r.AuthModule)
	return r
}

func (r *Registry) WithOtpModule(clients *clients.ClientRegistry) *Registry {
	r.OtpModule = otpmodule.GetOtpModule(clients)
	return r
}
