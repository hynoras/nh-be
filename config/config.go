package config

import (
	"nh-be/pkg/env"
)

type Config struct {
	AppEnv string
	Port   string

	// Database
	DBHost     string
	DBPort     int
	DBUsername string
	DBName     string
	DBPassword string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// RabbitMQ
	RabbitMQHost     string
	RabbitMQPort     int
	RabbitMQUsername string
	RabbitMQPassword string

	// Email / Frontend
	FrontendURL          string
	ResendAPIKey         string
	VerifyEmailSuffixURL string
}

func LoadConfig() *Config {
	cfg := &Config{
		AppEnv: env.MustEnv("APP_ENV"),
		Port:   env.MustEnv("PORT"),

		DBHost:     env.MustEnv("DB_HOST"),
		DBPort:     env.MustEnvInt("DB_PORT"),
		DBUsername: env.MustEnv("DB_USERNAME"),
		DBName:     env.MustEnv("DB_NAME"),
		DBPassword: env.MustEnv("DB_PASSWORD"),

		RedisHost:     env.MustEnv("REDIS_HOST"),
		RedisPort:     env.MustEnv("REDIS_PORT"),
		RedisPassword: env.MustEnv("REDIS_PASSWORD"),

		RabbitMQHost:     env.MustEnv("RABBITMQ_HOST"),
		RabbitMQPort:     env.MustEnvInt("RABBITMQ_PORT"),
		RabbitMQUsername: env.MustEnv("RABBITMQ_USERNAME"),
		RabbitMQPassword: env.MustEnv("RABBITMQ_PASSWORD"),
	}

	if cfg.AppEnv == "prod" {
		cfg.FrontendURL = env.MustEnv("FRONTEND_URL")
		cfg.ResendAPIKey = env.MustEnv("RESEND_API_KEY")
		cfg.VerifyEmailSuffixURL = env.MustEnv("VERIFY_EMAIL_SUFFIX_URL")
	}

	return cfg
}
