package redis_client

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0})
	if err := client.Ping(context.Background()).Err(); err != nil {

		panic(err.Error())
	}
	return client
}
