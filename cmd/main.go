package main

// @title NoHeir API
// @version 1.0
// @description API server for NoHeir application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@noheir.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey SessionAuth
// @in cookie
// @name auth_session

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nh-be/config"
	docs "nh-be/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"nh-be/internal/app"
	"nh-be/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("Warning: No .env file found")
	}

	cfg := config.LoadConfig()

	slog.Info("starting app", "env", cfg.AppEnv)

	service, err := app.InitializeServices(cfg)

	if err != nil {
		slog.Error("failed to initialize services", "error", err)
		os.Exit(1)
	}

	defer func() {
		if service == nil {
			return
		}
		service.ConCancel()
		service.WG.Wait()

		if service.PubCh != nil {
			service.PubCh.Close()
			slog.Info("RabbitMQ publisher channel closed")
		}
		if service.ConCh != nil {
			service.ConCh.Close()
			slog.Info("RabbitMQ consumer channel closed")
		}
		if service.RabbitMQ != nil {
			service.RabbitMQ.Close()
			slog.Info("RabbitMQ connection closed")
		}
		if service.SQLDB != nil {
			service.SQLDB.Close()
			slog.Info("Database connection closed")
		}
		if service.Redis != nil {
			service.Redis.Close()
			slog.Info("Redis connection closed")
		}
	}()

	var origin []string
	if cfg.AppEnv == "dev" {
		origin = []string{"http://localhost:3000"}
	} else {
		origin = []string{cfg.FrontendURL}
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origin,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	app.RegisterHealthRoutes(r, app.HealthDeps{
		SQLDB:    service.SQLDB,
		Redis:    service.Redis,
		RabbitMQ: service.RabbitMQ,
	})

	router.SetupRoutes(r, service.DB, service.Redis, service.PubCh)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		slog.Info("Server is running", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown: %v", err)
	}

	slog.Info("Server exited cleanly")
}
