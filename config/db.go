package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USERNAME")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

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
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to PostgreSQL using pgx-stdlib (IPv4 mode)")
	return db
}
