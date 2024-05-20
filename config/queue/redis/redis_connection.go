package redis

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func InitQueue() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})
}
