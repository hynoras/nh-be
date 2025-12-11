package user

import (
	"nh-be/internal/permission"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	usersGroup := rg.Group("/users")
	userRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	userService := NewService(userRepo, permissionRepo, permissionService)

	usersGroup.GET("", GetAllUsersHandler(userService))
	usersGroup.GET("/:id", GetUserByIDHandler(userService))
	usersGroup.GET("/me", GetMeHandler(userService))
	usersGroup.POST("", CreateUserHandler(userService))
	usersGroup.PUT("/:id", UpdateUserHander(userService))
	usersGroup.DELETE("", DeleteUsersHandler(userService))
}
