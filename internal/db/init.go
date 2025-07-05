package db

import (
	"log"
	gmailclient "paytm-project/internal/clients/email_client/gmail_client"
	authmodule "paytm-project/internal/modules/auth_module"
	paymentsmodule "paytm-project/internal/modules/payments_module"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dsn             string = "root:12345678@tcp(localhost:3306)/paytm?charset=utf8mb4&parseTime=True&loc=Local"
	DbConnection    *gorm.DB
	GetDbConnection = func() *gorm.DB {
		if DbConnection != nil {
			return DbConnection
		}
		return nil
	}
)

func InitDb() {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DbConnection = db

	err = db.AutoMigrate(&authmodule.User{}, &authmodule.Wallet{}, &paymentsmodule.Payment{}, &gmailclient.TokenStorage{})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("Database migrated successfully 🚀")
}
