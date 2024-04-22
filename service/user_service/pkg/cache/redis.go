package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Redis *redis.Client

var Ctx = context.Background()

func InitRedis() {
	addr := viper.GetString("redis.address")
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}
