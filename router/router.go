package router

import (
	"nh-be/internal/features/auth"
	"nh-be/internal/features/experiment"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/internal/infra"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetupRoutes initializes all application routes
func SetupRoutes(r *gin.Engine, db *gorm.DB, rdb *redis.Client, ch *amqp.Channel) {
	sessionStore := infra.NewSessionStore(rdb)
	// API version 1 group
	v1 := r.Group("/api/v1")

	// Register auth routes (public)
	auth.RegisterRoutes(v1, db, rdb, ch)

	// Protected routes group
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(sessionStore))

	user.RegisterRoutes(protected, db, rdb)
	permission.RegisterRoutes(protected, db)
	experiment.RegisterRoutes(protected, db)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "Not Found",
			"message": "The requested endpoint does not exist",
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"success": false,
			"error":   "Method Not Allowed",
			"message": "The requested method is not supported for this endpoint",
		})
	})
}
