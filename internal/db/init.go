package db

import (
	"log"
	"os"
	gmailclient "paytm-project/internal/clients/email_client/gmail_client"
	authmodule "paytm-project/internal/modules/auth_module"
	paymentsmodule "paytm-project/internal/modules/payments_module"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DbConnection    *gorm.DB
	GetDbConnection = func() *gorm.DB {
		if DbConnection != nil {
			return DbConnection
		}
		return nil
	}
)

func InitDb() {
	var dsn string = os.Getenv("DB_CONNECTION_URL")

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
