package permission

import (
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, s Service) {

	// Permissions
	permissions := r.Group("/permissions")
	{
		permissions.GET("", GetAllPermissionsHandler(s))
		permissions.GET("/:id", GetPermissionHandler(s))
	}

	permissions.Use(middleware.WithService("permission-service"))

	// Permission Groups
	groups := r.Group("/permission-groups")
	groups.Use(middleware.WithService("permission-service"))
	{
		groups.GET("", GetAllPermissionGroupsHandler(s))
		groups.POST("", CreatePermissionGroupHandler(s))
		groups.GET("/:id", GetPermissionGroupHandler(s))
		groups.PUT("/:id", UpdatePermissionGroupHandler(s))
		groups.DELETE("/:id", DeletePermissionGroupHandler(s))
	}
}
