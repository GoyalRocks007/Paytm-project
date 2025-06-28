package main

import (
	"paytm-project/boot"
	"paytm-project/internal/controllers"
	"paytm-project/internal/db"
	"paytm-project/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db.InitDb()
	// Register controller routes
	moduleRegistry := boot.InitModuleRegistry(db.GetDbConnection())
	authController := controllers.NewAuthController(moduleRegistry.AuthModule)
	paymentsController := controllers.NewPaymentsController(moduleRegistry.PaymentsModule)
	adminController := controllers.NewAdminController(moduleRegistry.AdminModule)

	adminRoutes := router.Group("/admin", middlewares.AuthMiddleware(), middlewares.CheckAdmin())
	{
		adminRoutes.PUT("/role", adminController.MakeAdmin)
	}

	router.GET("/hello", authController.HelloWorld)
	router.POST("auth/signup", authController.CreateUser)
	router.POST("auth/signin", authController.UserLogin)

	router.POST("/payment/create", paymentsController.CreatePayment)
	// Start server
	router.Run(":8080")
}
