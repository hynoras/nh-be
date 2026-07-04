package router

import (
	"nh-be/internal/app"
	"nh-be/internal/features/auth"
	"nh-be/internal/features/experiment"
	"nh-be/internal/features/experiment/result"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/procedure"
	"nh-be/internal/features/user"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all application routes
func SetupRoutes(r *gin.Engine, deps *app.SharedDeps) {
	// API version 1 group
	v1 := r.Group("/api/v1")
	v1.Use(middleware.SetCSRFToken())

	// Register auth routes (public)
	auth.RegisterRoutes(v1, deps)

	// Protected routes group
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(deps.SessionStore))
	protected.Use(middleware.CSRFProtection())

	user.RegisterRoutes(protected, deps)
	permission.RegisterRoutes(protected, deps.PermissionService)
	experiment.RegisterRoutes(protected, deps)
	result.RegisterRoutes(protected, deps)
	procedure.RegisterRoutes(protected, deps)

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
