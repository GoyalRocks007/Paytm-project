package main

import (
	"os"
	"paytm-project/boot"
	"paytm-project/internal/controllers"
	"paytm-project/internal/db"
	"paytm-project/middlewares"
	"paytm-project/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load(".env")
		if err != nil {
			panic("env file not found")
		}
	}
	router := gin.Default()
	db.InitDb()
	redis.InitRedis()
	// Register controller routes
	clientRegistery := boot.InitClientRegistery(db.GetDbConnection(), redis.GetRedisClientConnection())
	moduleRegistry := boot.InitModuleRegistry(db.GetDbConnection(), clientRegistery)
	authController := controllers.NewAuthController(moduleRegistry.AuthModule)
	paymentsController := controllers.NewPaymentsController(moduleRegistry.PaymentsModule)
	adminController := controllers.NewAdminController(moduleRegistry.AdminModule, clientRegistery.EmailClient)
	otpController := controllers.NewOtpController(moduleRegistry.OtpModule)

	adminRoutes := router.Group("/admin", middlewares.AuthMiddleware(), middlewares.CheckAdmin())
	{
		adminRoutes.PUT("/role", adminController.UpdateRole)
		adminRoutes.GET("/auth", adminController.StartAuth)
	}

	otpRoutes := router.Group("/otp")
	{
		otpRoutes.POST("/generate", otpController.GenerateOtp)
		otpRoutes.POST("/verify", otpController.VerifyOtp)
	}
	router.GET("/oauth2callback", adminController.HandleCallback)
	router.GET("/hello", authController.HelloWorld)
	router.POST("auth/signup", authController.CreateUser)
	router.POST("auth/signin", authController.UserLogin)

	router.POST("/payment/create", paymentsController.CreatePayment)
	// Start server
	router.Run(":8080")
}
