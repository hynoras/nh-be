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
	"nh-be/internal/user"
	"nh-be/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env early
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize database
	db := config.ConnectDatabase()
	db.AutoMigrate(&user.User{})
	sqlDB, _ := db.DB()
	defer func() {
		if sqlDB != nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}
	}()

	// Create router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
    	AllowOrigins:     []string{"http://localhost:3000"}, // Specify allowed origins
    	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    	ExposeHeaders:    []string{"Content-Length"}, // Headers the client can access
    	AllowCredentials: true,
    	MaxAge:           12 * time.Hour, 
    }))

	// Initialize session store
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "noheir_secret" // fallback for local dev
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		MaxAge:   60 * 60,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("auth_session", store))

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
		log.Println("Server is running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
