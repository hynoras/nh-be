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
	"nh-be/mq"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"nh-be/config"
	"nh-be/internal/auth"
	experimentResult "nh-be/internal/experiment/result"
	experiment "nh-be/internal/experiment/root"
	"nh-be/internal/permission"
	"nh-be/internal/user"
	"nh-be/router"

	"github.com/gin-contrib/cors"
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
		&auth.VerificationToken{},
		&user.User{},
		&permission.Permission{},
		&permission.PermissionGroup{},
		&user.UserPermission{},
		&experiment.Experiment{},
		&experimentResult.ExperimentResult{},
	)

	rdb := config.NewRedisClient()

	conn, err := mq.NewRabbitMQConnection()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	pubCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a publisher channel: %v", err)
	}

	conCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a consumer channel: %v", err)
	}

	dqErr := mq.DeclareQueues(pubCh, auth.SendVerificationEmailQueue)
	if dqErr != nil {
		log.Fatalf("Failed to declare queue: %v", dqErr)
	}

	deErr := mq.DeclareExchange(conCh, auth.AuthExchangeName)
	if deErr != nil {
		log.Fatalf("Failed to declare exchange: %v", deErr)
	}

	bqErr := mq.BindQueue(
		conCh,
		auth.SendVerificationEmailQueue,
		auth.UserRegisteredRoutingKey,
		auth.AuthExchangeName,
	)
	if bqErr != nil {
		log.Fatalf("Failed to bind queue: %v", bqErr)
	}

	conCtx, conCancel := context.WithCancel(context.Background())
	authConsumer := auth.NewAuthConsumer(conCh)
	go authConsumer.ConsumeSendVerificationEmail(conCtx)
	defer conCancel()

	sqlDB, _ := db.DB()
	defer func() {
		if pubCh != nil {
			pubCh.Close()
			log.Println("RabbitMQ publisher channel closed")
		}
		if conCh != nil {
			conCh.Close()
			log.Println("RabbitMQ consumer channel closed")
		}
		if conn != nil {
			conn.Close()
			log.Println("RabbitMQ connection closed")
		}
		if sqlDB != nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}
		if rdb != nil {
			rdb.Close()
			log.Println("Redis connection closed")
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

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.SetupRoutes(r, db, rdb, pubCh)

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
