package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	permissionRepo := NewRepository(db)
	permissionCache := NewPermissionCache(rdb)
	s := NewService(permissionRepo, permissionCache)

	// Permissions
	permissions := r.Group("/permissions")
	{
		permissions.GET("", GetAllPermissionsHandler(s))
		permissions.GET("/:id", GetPermissionHandler(s))
	}

	// Permission Groups
	groups := r.Group("/permission-groups")
	{
		groups.GET("", GetAllPermissionGroupsHandler(s))
		groups.POST("", CreatePermissionGroupHandler(s))
		groups.GET("/:id", GetPermissionGroupHandler(s))
		groups.PUT("/:id", UpdatePermissionGroupHandler(s))
		groups.DELETE("/:id", DeletePermissionGroupHandler(s))
	}
}
