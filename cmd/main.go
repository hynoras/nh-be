package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nh-be/config"
	"nh-be/internal/middleware"
	"nh-be/router"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env early
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: No .env file found")
	}

	// Initialize database
	db := config.ConnectDatabase()
	sqlDB, _ := db.DB()
	defer func() {
		if sqlDB != nil {
			sqlDB.Close()
			log.Println("🗃️  Database connection closed")
		}
	}()

	// Create router
	r := gin.Default()

	// Initialize session store
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "noheir_secret" // fallback for local dev
	}
	store := cookie.NewStore([]byte(secret))
	r.Use(sessions.Sessions("noheir_session", store))

	// Global middleware
	r.Use(middleware.ResponseFormatter())

	// Simple test route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Register routes
	router.SetupRoutes(r, db)

	// HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run server in goroutine
	go func() {
		log.Println("🚀 Server is running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server gracefully...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}
