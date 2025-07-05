package main

import (
	"context"
	"os"
	"paytm-project/boot"
	"paytm-project/internal/controllers"
	"paytm-project/internal/db"
	"paytm-project/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	db.InitDb()
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load(".env")
		if err != nil {
			panic("env file not find")
		}
	}
	// Register controller routes
	moduleRegistry := boot.InitModuleRegistry(db.GetDbConnection())
	clientRegistery := boot.InitClientRegistery(db.GetDbConnection())
	authController := controllers.NewAuthController(moduleRegistry.AuthModule)
	paymentsController := controllers.NewPaymentsController(moduleRegistry.PaymentsModule)
	adminController := controllers.NewAdminController(moduleRegistry.AdminModule, clientRegistery.EmailClient)

	clientRegistery.EmailClient.SendEmail(context.Background(), os.Getenv("GOOGLE_EMAIL_ADDRESS"), "abhi25goyal@gmail.com", "Jai Shri Ram", "Hi Abhi, This is to confirm the email client is working fine!")

	adminRoutes := router.Group("/admin", middlewares.AuthMiddleware(), middlewares.CheckAdmin())
	{
		adminRoutes.PUT("/role", adminController.UpdateRole)
		adminRoutes.GET("/auth", adminController.StartAuth)
	}
	router.GET("/oauth2callback", adminController.HandleCallback)
	router.GET("/hello", authController.HelloWorld)
	router.POST("auth/signup", authController.CreateUser)
	router.POST("auth/signin", authController.UserLogin)

	router.POST("/payment/create", paymentsController.CreatePayment)
	// Start server
	router.Run(":8080")
}
