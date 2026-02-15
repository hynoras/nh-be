package config

import (
	"log"
	"nh-be/pkg/env"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	host := env.MustEnv("REDIS_HOST")
	port := env.MustEnv("REDIS_PORT")
	pass := env.MustEnv("REDIS_PASSWORD")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
	})

	log.Println("Redis client created")
	return redisClient
}
