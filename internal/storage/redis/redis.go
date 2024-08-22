package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	return &RedisStore{}
}

func (s *RedisStore) Stop() error {
	if err := s.client.Close(); err != nil {
		return err
	}
	log.Info("🛑 Redis connection closed")
	return nil
}

func (s *RedisStore) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.adress"),
		Password: viper.GetString("redis.password"),
		Username: viper.GetString("redis.username"),

		// this should be set to true only for testing purposes
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		},
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.SetTimeFormat(time.Kitchen)
	log.Info("🔌 Redis Connected")
	s.client = client
	return nil
}
