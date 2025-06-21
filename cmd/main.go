package main

import (
	"paytm-project/boot"
	"paytm-project/internal/controllers"
	"paytm-project/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db.InitDb()
	// Register controller routes
	moduleRegistry := boot.InitModuleRegistry(db.GetDbConnection())
	authController := controllers.NewAuthController(moduleRegistry.AuthModule)
	paymentsController := controllers.NewPaymentsController(moduleRegistry.PaymentsModule)

	router.GET("/hello", authController.HelloWorld)
	router.POST("auth/signup", authController.CreateUser)
	router.POST("auth/signin", authController.UserLogin)

	router.POST("/payment/create", paymentsController.CreatePayment)
	// Start server
	router.Run(":8080")
}
