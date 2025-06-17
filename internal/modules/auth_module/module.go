package authmodule

type IAuthModule interface {
	GetCore() IAuthCore
}

type AuthModule struct {
	AuthCore IAuthCore
}

func (a *AuthModule) GetCore() IAuthCore {
	return a.AuthCore
}
