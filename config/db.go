package config

import (
	"fmt"
	"log"
	"nh-be/pkg/env"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() *gorm.DB {
	host := env.MustEnv("DB_HOST")
	port := env.MustEnvInt("DB_PORT")
	user := env.MustEnv("DB_USERNAME")
	dbname := env.MustEnv("DB_NAME")
	pass := env.MustEnv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=require",
		user, pass, host, port, dbname,
	)

	pgxCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse pgx config: %v", err)
	}

	sqlDB := stdlib.OpenDB(*pgxCfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Connection Pool
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Successfully connected to PostgreSQL")
	return db
}
