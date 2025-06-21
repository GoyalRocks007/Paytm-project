package controllers

import (
	"net/http"
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

func (pc *PaymentsController) CreatePayment(c *gin.Context) {

	var req paymentsmodule.CreatePaymentRequestDto

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	res, err := pc.paymentsModule.GetCore().CreatePayment(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
