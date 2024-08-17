package storage

import (
	"time"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Redis) *redis.Client {
	redisHost := cfg.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		Password:     cfg.RedisPass,
		DB:           cfg.RedisDB,
	})

	return client
}
