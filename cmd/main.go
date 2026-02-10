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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	docs "nh-be/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"nh-be/config"
	experimentResult "nh-be/internal/experiment/result"
	experiment "nh-be/internal/experiment/root"
	"nh-be/internal/permission"
	"nh-be/internal/user"
	"nh-be/router"
	"nh-be/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	if env != "prod" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: No .env file found")
		}
	}

	log.Printf("Starting app in %s mode\n", env)

	db := config.ConnectDatabase()
	db.AutoMigrate(
		&user.User{},
		&permission.Permission{},
		&permission.PermissionGroup{},
		&user.UserPermission{},
		&experiment.Experiment{},
		&experimentResult.ExperimentResult{},
	)
	config.ConnectRabbitMQ()
	sqlDB, _ := db.DB()
	defer func() {
		if sqlDB != nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}
	}()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	secret := utils.MustEnv("SESSION_SECRET")

	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		MaxAge:   8 * 60 * 60, // 8 hours
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("auth_session", store))

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.SetupRoutes(r, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Println("Server is running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
