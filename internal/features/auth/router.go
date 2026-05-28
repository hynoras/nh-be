package auth

import (
	"nh-be/internal/app"
	"nh-be/internal/features/user"
	"nh-be/internal/middleware"
	"nh-be/internal/platform/email"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, deps *app.SharedDeps) {
	authGroup := rg.Group("/auth")
	authGroup.Use(middleware.WithService("auth-service"))

	authRepo := NewRepository(deps.DB)
	userRepo := user.NewRepository(deps.DB)
	emailPublisher := email.NewEmailPublisher(deps.PubCh)
	authService := NewService(deps.SessionStore, authRepo, userRepo, deps.PermissionService, emailPublisher)

	authGroup.POST("/signup", SignUpHandler(authService))
	authGroup.POST("/login", LoginHandler(authService))
	authGroup.POST("/logout", LogoutHandler(authService))
	authGroup.Use(middleware.RequireAuth(deps.SessionStore)).PUT("change-password/:id", ChangePasswordHandler((authService)))
	authGroup.GET("verify/:token", VerifyTokenHandler(authService))
}
