package controllers

import (
	"net/http"
	otpmodule "paytm-project/internal/modules/otp_module"

	"github.com/gin-gonic/gin"
)

type OtpController struct {
	otpModule otpmodule.IOtpModule
}

func NewOtpController(otpmodule otpmodule.IOtpModule) *OtpController {
	return &OtpController{
		otpModule: otpmodule,
	}
}

func (ac *OtpController) GenerateOtp(c *gin.Context) {
	var req otpmodule.GenerateOtpRequest
	ctx := c.Request.Context()

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	res, err := ac.otpModule.GetCore().GenerateOtp(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ac *OtpController) VerifyOtp(c *gin.Context) {
	var req otpmodule.VerifyOtpRequest
	ctx := c.Request.Context()

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	res, token, err := ac.otpModule.GetCore().VerifyOtp(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Header("two-fac-auth", token)

	c.JSON(http.StatusOK, res)
}
