package serve

import (
	"context"
	"fyoukuapi/config"
	"github.com/go-redis/redis/v8"
)

var (
	Rdb  *redis.Client
	Rctx context.Context
)

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
	})
	Rctx = context.Background()
}
