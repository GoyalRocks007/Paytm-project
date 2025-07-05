package adminmodule

import (
	"paytm-project/internal/models"
)

type IAdminRepo interface {
}

type AdminRepo struct {
	models.BaseRepo
}
