package redis

import (
	"context"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *config.AppConfig) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       0,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil

}
