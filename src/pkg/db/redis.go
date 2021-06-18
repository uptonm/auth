package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/uptonm/auth/src/common"
)

var RedisConn *redis.Client

// InitRedis initializes the db connection to redis, currently being used to store state codes
func InitRedis() error {
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", common.Config.RedisHost, common.Config.RedisPort),
		Password: common.Config.RedisPass,
		DB:       0,
	})

	return RedisConn.Ping(context.Background()).Err()
}
