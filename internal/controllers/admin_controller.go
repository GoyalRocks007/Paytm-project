package controllers

import (
	"net/http"

	emailclient "paytm-project/internal/clients/email_client"
	adminmodule "paytm-project/internal/modules/admin_module"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	AdminModule adminmodule.IAdminModule
	EmailClient emailclient.IEmailClient
}

func NewAdminController(adminmodule adminmodule.IAdminModule, emailClient emailclient.IEmailClient) *AdminController {
	return &AdminController{
		AdminModule: adminmodule,
		EmailClient: emailClient,
	}
}

func (ac *AdminController) UpdateRole(c *gin.Context) {

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

func (ac *AdminController) StartAuth(c *gin.Context) {
	ctx := c.Request.Context()
	url, err := ac.EmailClient.StartAuth(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "OAuth config error: " + err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, url)
}

func (ac *AdminController) HandleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing authorization code",
		})
		return
	}
	err := ac.EmailClient.HandleCallback(ctx, code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token exchange error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful! Token saved to database",
	})
}
