package adminmodule

import (
	"errors"
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"
)

type IAdminCore interface {
	UpdateRole(req *UpdateRoleRequestDto) (*UpdateRoleResponseDto, error)
}

type AdminCore struct {
	AdminRepo  IAdminRepo
	AuthModule authmodule.IAuthModule
}

func (c *AdminCore) UpdateRole(req *UpdateRoleRequestDto) (*UpdateRoleResponseDto, error) {
	role := req.Role
	if !role.ValidateRole() {
		return &UpdateRoleResponseDto{
			BaseSuccessResponse: models.BaseSuccessResponse{
				Success: false,
			},
		}, errors.New("invalid role")
	}
	err := c.AuthModule.GetCore().GetRepo().UpdateUserRole(&authmodule.User{
		Email: req.Email,
		Role:  req.Role,
	})
	if err != nil {
		return &UpdateRoleResponseDto{
			BaseSuccessResponse: models.BaseSuccessResponse{
				Success: false,
			}}, err
	}
	return &UpdateRoleResponseDto{
		BaseSuccessResponse: models.BaseSuccessResponse{
			Success: true,
		}}, nil
}
