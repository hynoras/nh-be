package config

import (
	"fmt"
	"log/slog"
	infradb "nh-be/infra/db"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?statement_cache_capacity=0&default_query_exec_mode=exec",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	if cfg.DBSslMode != "" {
		dsn += fmt.Sprintf("&sslmode=%s", cfg.DBSslMode)
	} else {
		dsn += "&sslmode=disable"
	}

	slog.Info("Database SSL Mode", "sslmode", cfg.DBSslMode)

	pgxCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	sqlDB := stdlib.OpenDB(*pgxCfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Use(&infradb.DbMetricsPlugin{}); err != nil {
		return nil, fmt.Errorf("failed to register db metrics plugin: %w", err)
	}

	if err := db.Use(tracing.NewPlugin(
		tracing.WithoutMetrics(),
		tracing.WithoutQueryVariables(),
	)); err != nil {
		return nil, fmt.Errorf("failed to register db tracing plugin: %w", err)
	}

	// Connection Pool
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	slog.Info("Successfully connected to PostgreSQL")
	return db, nil
}
