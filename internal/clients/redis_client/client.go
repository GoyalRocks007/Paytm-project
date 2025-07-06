package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClientInstance IRedisClient
	GetRedisClient      = func(rdb *redis.Client) IRedisClient {
		if redisClientInstance != nil {
			return redisClientInstance
		}

		redisClientInstance = New(rdb)
		return redisClientInstance
	}
)

type IRedisClient interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, bool, error)
	Del(ctx context.Context, key string) error
}
