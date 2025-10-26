package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	usersGroup := rg.Group("/users")
	userRepo := NewRepository(db)
	userService := NewService(userRepo)

	usersGroup.GET("", GetAllUsersHandler(userService))
	usersGroup.GET("/:id", GetUserByIDHandler(userService))
	usersGroup.GET("/me", GetMeHandler(userService))
}
