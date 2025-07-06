package redisclient

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) IRedisClient {
	return &redisClient{rdb: rdb}
}

func (r *redisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, key, data, ttl).Err()
}

func (r *redisClient) Get(ctx context.Context, key string) (interface{}, bool, error) {
	var result interface{}
	data, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key not found — return zero value
		return nil, false, nil
	} else if err != nil {
		log.Println("some error occured while fetching key", key, err)
		return nil, false, err
	}
	err = json.Unmarshal([]byte(data), &result)
	return result, true, err
}

func (r *redisClient) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}
