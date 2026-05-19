package user

import (
	"nh-be/internal/features/permission"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	usersGroup := rg.Group("/users")
	usersGroup.Use(middleware.WithService("user-service"))

	userRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionCache := permission.NewPermissionCache(rdb)
	permissionService := permission.NewService(permissionRepo, permissionCache)
	userService := NewService(userRepo, permissionRepo, permissionService)

	usersGroup.GET("", GetAllUsersHandler(userService))
	usersGroup.GET("/:id", GetUserByIDHandler(userService))
	usersGroup.GET("/me", GetMeHandler(userService))
	usersGroup.POST("", CreateUserHandler(userService))
	usersGroup.PUT("/:id", UpdateUserHandler(userService))
	usersGroup.DELETE("", DeleteUsersHandler(userService))
}
