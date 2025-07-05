package internal

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// Shared context for Redis operations
var Ctx = context.Background()

// NewRedisClient returns a new Redis client
func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
