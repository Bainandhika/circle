package nosql

import (
	"circle-fiber/app/config"
	"circle-fiber/lib/logger"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func RedisConnect() *redis.Client {
	redisConfig := config.Redis
	redisConn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Username: redisConfig.Username,
		Password: redisConfig.Password,
	})

	err := redisConn.Ping(context.Background()).Err()
	if err != nil {
		logger.Error.Printf("Failed to connect to Redis! err: %v", err)
		panic(err)
	}

	return redisConn
}
