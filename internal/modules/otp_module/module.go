package otpmodule

import (
	"paytm-project/internal/clients"
)

var (
	otpModuleInstance IOtpModule
	GetOtpModule      = func(clients *clients.ClientRegistry) IOtpModule {
		if otpModuleInstance != nil {
			return otpModuleInstance
		}
		return initOtpModule(clients)
	}
	initOtpModule = func(clients *clients.ClientRegistry) IOtpModule {
		otpModuleInstance = &OtpModule{
			OtpCore: NewOtpCore(clients),
		}
		return otpModuleInstance
	}
)

type IOtpModule interface {
	GetCore() IOtpCore
}

type OtpModule struct {
	OtpCore IOtpCore
}

func (om *OtpModule) GetCore() IOtpCore {
	return om.OtpCore
}
