package nosql

import (
	"context"
	"fmt"

	"circle-2.0/app/config"
	"circle-2.0/lib/logger"

	"github.com/go-redis/redis/v8"
)

var RedisConnect *redis.Client

func InitRedis() {
	redisConfig := config.Redis
	redisConn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Username: redisConfig.Username,
		Password: redisConfig.Password,
	})

	err := redisConn.Ping(context.Background()).Err()
	if err != nil {
		err = fmt.Errorf("Failed to connect to Redis! err: %v", err)
		logger.Error.Print(err.Error())
		panic(err)
	}

	RedisConnect = redisConn
}
