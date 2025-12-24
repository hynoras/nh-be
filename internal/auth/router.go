package auth

import (
	"nh-be/internal/middleware"
	"nh-be/internal/permission"
	"nh-be/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	authGroup := rg.Group("/auth")
	userRepo := user.NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	authService := NewService(userRepo, permissionService)

	authGroup.POST("/login", LoginHandler(authService))
	authGroup.POST("/logout", LogoutHandler(authService))
	authGroup.Use(middleware.RequireAuth()).PUT("change-password/:id", ChangePasswordHandler((authService)))
}
