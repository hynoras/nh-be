package user

import (
	"nh-be/internal/app"
	"nh-be/internal/features/permission"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, deps *app.SharedDeps) {
	usersGroup := rg.Group("/users")
	usersGroup.Use(middleware.WithService("user-service"))

	userRepo := NewRepository(deps.DB)
	permissionRepo := permission.NewRepository(deps.DB)
	userService := NewService(userRepo, permissionRepo, deps.PermissionService)

	usersGroup.GET("", GetAllUsersHandler(userService))
	usersGroup.GET("/:id", GetUserByIDHandler(userService))
	usersGroup.GET("/me", GetMeHandler(userService))
	usersGroup.POST("", CreateUserHandler(userService))
	usersGroup.PUT("/:id", UpdateUserHandler(userService))
	usersGroup.DELETE("", DeleteUsersHandler(userService))
}
