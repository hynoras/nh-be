package config

import (
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	log.Println("Redis client created")
	return redisClient
}
