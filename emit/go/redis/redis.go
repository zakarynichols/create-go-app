package redis

import (
	"github.com/redis/go-redis/v9"
)

func Open() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return client
}
