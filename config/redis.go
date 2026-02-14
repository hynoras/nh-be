package config

import (
	"log"
	"nh-be/utils"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	host := utils.MustEnv("REDIS_HOST")
	port := utils.MustEnv("REDIS_PORT")
	pass := utils.MustEnv("REDIS_PASSWORD")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pass,
		DB:       0,
	})

	log.Println("Redis client created")
	return redisClient
}
