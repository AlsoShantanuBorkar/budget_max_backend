package redis

import (
	"context"

	"github.com/AlsoShantanuBorkar/budget_max/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis() {

	Client = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisHost,
		Password: config.Config.RedisPassword,
		DB:       0,
	})

	_, err := Client.Ping(Ctx).Result()

	if err != nil {
		panic(err)
	}

}
