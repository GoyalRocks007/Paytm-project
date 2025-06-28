package adminmodule

import (
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"
)

type UpdateRoleRequestDto struct {
	Email string          `json:"email"`
	Role  authmodule.Role `json:"role"`
}

type UpdateRoleResponseDto struct {
	models.BaseSuccessResponse
}
