package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis *redis.Client
}

func Load() *Config {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}

	return &Config{
		Redis: rdb,
	}
}
