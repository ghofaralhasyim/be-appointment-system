package database

import (
	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func NewRedisClient(option redis.Options) *redis.Client {
	client = redis.NewClient(&option)

	return client
}
