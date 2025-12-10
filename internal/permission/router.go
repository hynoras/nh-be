package permission

import (
	"nh-be/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	permissionRepo := NewRepository(db)
	userRepo := user.NewRepository(db)
	s := NewService(permissionRepo, userRepo)

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
		
		// User Assignments
		groups.POST("/assign", AssignUserToGroupHandler(s))
		groups.DELETE("/:id/users/:userId", RemoveUserFromGroupHandler(s))
	}
}
