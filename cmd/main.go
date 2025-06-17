package main

import (
	"paytm-project/internal/controllers"
	"paytm-project/internal/db"
	"paytm-project/internal/models"
	authmodule "paytm-project/internal/modules/auth_module"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initAuthModule(db *gorm.DB) authmodule.IAuthModule {
	return &authmodule.AuthModule{
		AuthCore: &authmodule.AuthCore{
			AuthRepo: &authmodule.AuthRepo{
				BaseRepo: models.BaseRepo{
					Db: db,
				},
			},
		},
	}
}
func main() {
	router := gin.Default()
	db.InitDb()
	// Register controller routes
	authController := controllers.NewAuthController(initAuthModule(db.GetDbConnection()))

	router.GET("/hello", authController.HelloWorld)
	router.POST("/signup", authController.CreateUser)

	// Start server
	router.Run(":8080")
}
