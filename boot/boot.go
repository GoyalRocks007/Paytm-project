package boot

import (
	"paytm-project/internal/clients"
	"paytm-project/internal/modules"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitModuleRegistry(db *gorm.DB) *modules.Registry {
	return modules.GetRegistry().
		WithAuthModule(db).
		WithPaymentsModule(db).
		WithAdminModule(db)
}

func InitClientRegistery(db *gorm.DB, rdb *redis.Client) *clients.ClientRegistry {
	return clients.GetRegistry().
		WithEmailClient(db).
		WithRedisClient(rdb)
}
