package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

// NewRedis: Redis bağlantısını başlatır
func NewRedis(addr string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test için ping atıyoruz
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

// Get: Veriyi çeker
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set: Veriyi kaydeder (TTL: ne kadar süre saklanacağı)
func (r *RedisClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}
