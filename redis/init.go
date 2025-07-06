package redis

import (
	"context"
	"fmt"
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
	addr := os.Getenv("REDIS_HOST")         // e.g., "localhost:6379"
	password := os.Getenv("REDIS_PASSWORD") // "" if no password
	db := 0                                 // can make this env-controlled too

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	fmt.Println("✅ Redis connected")
	return nil
}
