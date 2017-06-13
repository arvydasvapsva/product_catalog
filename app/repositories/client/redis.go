package client

import (
	"github.com/go-redis/redis"
)

type RedisClient struct {
}

func (*RedisClient) GetClient() *redis.Client  {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}