package config

import (
	"os"

	"github.com/ghofaralhasyim/be-appointment-system/pkg/utils"
	"github.com/go-redis/redis/v8"
)

func NewRedisConfig() (redis.Options, error) {

	var client redis.Options

	connUrl, err := utils.ConnURLBuilder("redis")
	if err != nil {
		return client, err
	}

	client = redis.Options{
		Addr:     connUrl,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	return client, nil
}
