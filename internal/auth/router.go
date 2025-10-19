package auth

import (
	"nh-be/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	authGroup := rg.Group("/auth")
	userRepo := user.NewRepository(db)
	authService := NewService(userRepo)

	authGroup.POST("/login", LoginHandler(authService))
	// authGroup.POST("/register", RegisterHandler(service))
}
