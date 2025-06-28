package controllers

import (
	"net/http"
	adminmodule "paytm-project/internal/modules/admin_module"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	AdminModule adminmodule.IAdminModule
}

func NewAdminController(adminmodule adminmodule.IAdminModule) *AdminController {
	return &AdminController{
		AdminModule: adminmodule,
	}
}

func (ac *AdminController) MakeAdmin(c *gin.Context) {

	var req adminmodule.UpdateRoleRequestDto

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	res, err := ac.AdminModule.GetCore().UpdateRole(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
