package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis *redis.Client
}

func Load() *Config {
	addr := os.Getenv("REDIS_ADDR")
	pass := os.Getenv("REDIS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}

	log.Println("Redis connected")

	return &Config{
		Redis: rdb,
	}
}
