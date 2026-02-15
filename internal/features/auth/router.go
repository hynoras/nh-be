package auth

import (
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/internal/infra"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, rdb *redis.Client, ch *amqp.Channel) {
	authGroup := rg.Group("/auth")
	authRepo := NewRepository(db)
	userRepo := user.NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	authPublisher := NewAuthPublisher(ch)
	sessionStore := infra.NewSessionStore(rdb)
	authService := NewService(sessionStore, authRepo, userRepo, permissionService, authPublisher)

	authGroup.POST("/signup", SignUpHandler(authService))
	authGroup.POST("/login", LoginHandler(authService))
	authGroup.POST("/logout", LogoutHandler(authService))
	authGroup.Use(middleware.RequireAuth(sessionStore)).PUT("change-password/:id", ChangePasswordHandler((authService)))
	authGroup.GET("verify/:token", VerifyTokenHandler(authService))
}
