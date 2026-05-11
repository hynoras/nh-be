package auth

import (
	"nh-be/internal/email"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/internal/middleware"
	"nh-be/internal/platform/session"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, rdb *redis.Client, ch *amqp.Channel) {
	authGroup := rg.Group("/auth")
	authGroup.Use(middleware.WithService("auth-service"))

	authRepo := NewRepository(db)
	userRepo := user.NewRepository(db)
	emailPublisher := email.NewEmailPublisher(ch)
	permissionRepo := permission.NewRepository(db)
	permissionCache := permission.NewPermissionCache(rdb)
	permissionService := permission.NewService(permissionRepo, permissionCache)
	sessionStore := session.NewSessionStore(rdb)
	authService := NewService(sessionStore, authRepo, userRepo, permissionService, emailPublisher)

	authGroup.POST("/signup", SignUpHandler(authService))
	authGroup.POST("/login", LoginHandler(authService))
	authGroup.POST("/logout", LogoutHandler(authService))
	authGroup.Use(middleware.RequireAuth(sessionStore)).PUT("change-password/:id", ChangePasswordHandler((authService)))
	authGroup.GET("verify/:token", VerifyTokenHandler(authService))
}
