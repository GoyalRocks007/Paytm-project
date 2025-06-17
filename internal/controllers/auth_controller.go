package controllers

import (
	"net/http"
	authmodule "paytm-project/internal/modules/auth_module"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authModule authmodule.IAuthModule
}

func NewAuthController(authmodule authmodule.IAuthModule) *AuthController {
	return &AuthController{
		authModule: authmodule,
	}
}

func (ac *AuthController) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from Uddyan!",
	})
}

func (ac *AuthController) CreateUser(c *gin.Context) {

	var req authmodule.SignupRequestDto

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	res, err := ac.authModule.GetCore().UserSignup(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ac *AuthController) UserLogin(c *gin.Context) {
	var req authmodule.LoginRequestDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	token, err := ac.authModule.GetCore().UserLogin(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("X-Auth-Token", token)
	c.JSON(http.StatusOK, authmodule.LoginResponseDto{Success: true})
}
