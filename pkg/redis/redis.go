package redis

import (
	"context"

	"github.com/CollabTED/CollabTed-Backend/pkg/logger"
	"github.com/redis/go-redis/v9"
)

var ctx = context.TODO()
var client *redis.Client

func Connect() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		// #nosec G402
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Logger.Err(err)
	}

	logger.Logger.Info().Msg("Connected to Redis successfully!")
}

func GetClient() *redis.Client {
	return client
}
