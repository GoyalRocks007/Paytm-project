package adminmodule

import (
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"

	"gorm.io/gorm"
)

var (
	adminModuleInstance IAdminModule
	GetAdminModule      = func(db *gorm.DB, authModule authmodule.IAuthModule) IAdminModule {
		if adminModuleInstance != nil {
			return adminModuleInstance
		}
		return initAdminModule(db, authModule)
	}
	initAdminModule = func(db *gorm.DB, authModule authmodule.IAuthModule) IAdminModule {
		adminModuleInstance = &AdminModule{
			AdminCore: &AdminCore{
				AuthModule: authModule,
				AdminRepo: &AdminRepo{
					BaseRepo: models.BaseRepo{
						Db: db,
					},
				},
			},
		}
		return adminModuleInstance
	}
)

type IAdminModule interface {
	GetCore() IAdminCore
}

type AdminModule struct {
	AdminCore IAdminCore
}

func (am *AdminModule) GetCore() IAdminCore {
	return am.AdminCore
}
