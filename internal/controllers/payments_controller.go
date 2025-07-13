package controllers

import (
	"net/http"
	otpmodule "paytm-project/internal/modules/otp_module"
	paymentsmodule "paytm-project/internal/modules/payments_module"

	"github.com/gin-gonic/gin"
)

type PaymentsController struct {
	paymentsModule paymentsmodule.IPaymentsModule
}

func NewPaymentsController(paymentsmodule paymentsmodule.IPaymentsModule) *PaymentsController {
	return &PaymentsController{
		paymentsModule: paymentsmodule,
	}
}

func (pc *PaymentsController) InitiatePayment(c *gin.Context) {

	ctx := c.Request.Context()

	var req paymentsmodule.InitiatePaymentRequestDto

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	res, err := pc.paymentsModule.GetCore().InitiatePayment(ctx, &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (pc *PaymentsController) GenerateOtpForPayment(c *gin.Context) {

	ctx := c.Request.Context()

	paymentId := c.Param("id")

	res, err := pc.paymentsModule.GetCore().GenerateOtpForPayment(ctx, paymentId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (pc *PaymentsController) VerifyOtpForPayment(c *gin.Context) {

	ctx := c.Request.Context()

	var req otpmodule.VerifyOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	res, err := pc.paymentsModule.GetCore().AuthorizePayment(ctx, &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("payment-id", res)

	c.JSON(http.StatusOK, nil)
}

func (pc *PaymentsController) ExecutePayment(c *gin.Context) {

	ctx := c.Request.Context()

	paymentId := c.Param("id")

	res, err := pc.paymentsModule.GetCore().ExecutePayment(ctx, paymentId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
