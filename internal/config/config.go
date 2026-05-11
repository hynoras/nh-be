package config

import (
	"log/slog"
	"nh-be/internal/platform/vault"
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
	DBSslMode  string

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
	appEnv := env.MustEnv("APP_ENV")
	useVault := env.GetEnvOrDefaultBool("USE_VAULT", false)

	if useVault {
		return loadConfigFromVault(appEnv)
	}

	return loadConfigFromEnv(appEnv)
}

func loadConfigFromVault(appEnv string) *Config {
	// Initialize Vault
	vaultClient := vault.NewVaultClient()
	vault.AuthenticateVault(vaultClient)

	// Fetch vaults from Vault
	basePath := "noheir/" + appEnv

	dbSecrets, err := vault.GetSecret(vaultClient, basePath+"/db")
	if err != nil {
		slog.Error("failed to fetch db vaults from vault", "error", err)
		panic(err)
	}

	redisSecrets, err := vault.GetSecret(vaultClient, basePath+"/redis")
	if err != nil {
		slog.Error("failed to fetch redis vaults from vault", "error", err)
		panic(err)
	}

	rabbitmqSecrets, err := vault.GetSecret(vaultClient, basePath+"/rabbitmq")
	if err != nil {
		slog.Error("failed to fetch rabbitmq vaults from vault", "error", err)
		panic(err)
	}

	cfg := &Config{
		AppEnv: appEnv,
		Port:   env.MustEnv("PORT"),

		// DB: sensitive from Vault, non-sensitive from env
		DBHost:     vault.MustGetSecretValue(dbSecrets, "DB_HOST"),
		DBPort:     env.MustEnvInt("DB_PORT"),
		DBUsername: vault.MustGetSecretValue(dbSecrets, "DB_USERNAME"),
		DBName:     env.MustEnv("DB_NAME"),
		DBPassword: vault.MustGetSecretValue(dbSecrets, "DB_PASSWORD"),
		DBSslMode:  env.MustEnv("DB_SSL_MODE"),

		// Redis: sensitive from Vault, non-sensitive from env
		RedisHost:     vault.MustGetSecretValue(redisSecrets, "REDIS_HOST"),
		RedisPort:     env.MustEnv("REDIS_PORT"),
		RedisPassword: vault.MustGetSecretValue(redisSecrets, "REDIS_PASSWORD"),

		// RabbitMQ: password from Vault, rest from env
		RabbitMQHost:     env.MustEnv("RABBITMQ_HOST"),
		RabbitMQPort:     env.MustEnvInt("RABBITMQ_PORT"),
		RabbitMQUsername: env.MustEnv("RABBITMQ_USERNAME"),
		RabbitMQPassword: vault.MustGetSecretValue(rabbitmqSecrets, "RABBITMQ_PASSWORD"),
	}

	if cfg.AppEnv == "prod" {
		resendSecrets, err := vault.GetSecret(vaultClient, basePath+"/resend")
		if err != nil {
			slog.Error("failed to fetch resend vaults from vault", "error", err)
			panic(err)
		}

		cfg.FrontendURL = env.MustEnv("FRONTEND_URL")
		cfg.ResendAPIKey = vault.MustGetSecretValue(resendSecrets, "RESEND_API_KEY")
		cfg.VerifyEmailSuffixURL = env.MustEnv("VERIFY_EMAIL_SUFFIX_URL")
	}

	return cfg
}

func loadConfigFromEnv(appEnv string) *Config {
	slog.Info("loading config from .env (Vault disabled)")

	cfg := &Config{
		AppEnv: appEnv,
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
