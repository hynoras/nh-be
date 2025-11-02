package router

import (
	"nh-be/internal/auth"
	"nh-be/internal/middleware"
	"nh-be/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes initializes all application routes
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// API version 1 group
	v1 := r.Group("/api/v1")
	
	
	// Register auth routes (public)
	auth.RegisterRoutes(v1, db)
	
	// Protected routes group
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth())
	
	// Register protected user routes
	user.RegisterRoutes(protected, db)
	
	// Add other protected route registrations here as needed
	// product.RegisterRoutes(protected, db)
	
	// 404 handler for undefined routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "Not Found",
			"message": "The requested endpoint does not exist",
		})
	})
	
	// 405 handler for unsupported methods on existing routes
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"success": false,
			"error":   "Method Not Allowed",
			"message": "The requested method is not supported for this endpoint",
		})
	})
}
