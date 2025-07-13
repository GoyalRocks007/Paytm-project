package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	ctx                      = context.Background()
	RDB                      *redis.Client
	GetRedisClientConnection = func() *redis.Client {
		if RDB != nil {
			return RDB
		}
		return nil
	}
)

func InitRedis() error {
	// addr := os.Getenv("REDIS_HOST")         // e.g., "localhost:6379"
	// password := os.Getenv("REDIS_PASSWORD") // "" if no password
	// db := 0                                 // can make this env-controlled too
	// log.Println("twadi nu!!!")
	// RDB = redis.NewClient(&redis.Options{
	// 	Addr:      addr,
	// 	Password:  password,
	// 	DB:        db,
	// 	TLSConfig: &tls.Config{},
	// })

	redisURL := os.Getenv("REDIS_CONNECTION_URL")

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("failed to parse redis url: %v", err)
	}

	RDB := redis.NewClient(opt)

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %s", err.Error())
		return nil
	}

	log.Println("✅ Redis connected")
	return nil
}
