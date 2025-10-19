package main

import (
	"nh-be/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes initializes all application routes
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// API version 1 group
	v1 := r.Group("/api/v1")
	
	// Register auth routes
	auth.RegisterRoutes(v1, db)
	
	// Add other route registrations here as needed
	// user.RegisterRoutes(v1, db)
	// product.RegisterRoutes(v1, db)
}
