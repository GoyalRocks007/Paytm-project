package authmodule

import (
	"paytm-project/internal/models"

	"gorm.io/gorm"
)

var (
	authModuleInstance IAuthModule
	GetAuthModule      = func(db *gorm.DB) IAuthModule {
		if authModuleInstance != nil {
			return authModuleInstance
		}
		return initAuthModule(db)
	}
	initAuthModule = func(db *gorm.DB) IAuthModule {
		authModuleInstance = &AuthModule{
			AuthCore: &AuthCore{
				AuthRepo: &AuthRepo{
					BaseRepo: models.BaseRepo{
						Db: db,
					},
				},
			},
		}
		return authModuleInstance
	}
)

type IAuthModule interface {
	GetCore() IAuthCore
}

type AuthModule struct {
	AuthCore IAuthCore
}

func (a *AuthModule) GetCore() IAuthCore {
	return a.AuthCore
}
