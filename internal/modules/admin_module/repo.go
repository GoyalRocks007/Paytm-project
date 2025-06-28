package adminmodule

import "paytm-project/internal/models"

type IAdminRepo interface {
	Dummy() bool
}

type AdminRepo struct {
	models.BaseRepo
}

func (ar *AdminRepo) Dummy() bool {
	return true
}
