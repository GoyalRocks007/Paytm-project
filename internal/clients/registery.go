package clients

import (
	emailclient "paytm-project/internal/clients/email_client"
	redisclient "paytm-project/internal/clients/redis_client"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	registry *ClientRegistry
)

type ClientRegistry struct {
	EmailClient emailclient.IEmailClient
	RedisClient redisclient.IRedisClient
}

func GetRegistry() *ClientRegistry {
	if registry == nil {
		registry = &ClientRegistry{}
	}
	return registry
}

func (cr *ClientRegistry) WithEmailClient(db *gorm.DB) *ClientRegistry {
	cr.EmailClient = emailclient.GetEmailClient(db)
	return cr
}

func (cr *ClientRegistry) WithRedisClient(rdb *redis.Client) *ClientRegistry {
	cr.RedisClient = redisclient.GetRedisClient(rdb)
	return cr
}
