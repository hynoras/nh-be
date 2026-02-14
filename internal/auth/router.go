package auth

import (
	"nh-be/internal/middleware"
	"nh-be/internal/permission"
	"nh-be/internal/user"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, ch *amqp.Channel) {
	authGroup := rg.Group("/auth")
	authRepo := NewRepository(db)
	userRepo := user.NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	authPublisher := NewAuthPublisher(ch)
	authService := NewService(authRepo, userRepo, permissionService, authPublisher)

	authGroup.POST("/signup", SignUpHandler(authService))
	authGroup.POST("/login", LoginHandler(authService))
	authGroup.POST("/logout", LogoutHandler(authService))
	authGroup.Use(middleware.RequireAuth()).PUT("change-password/:id", ChangePasswordHandler((authService)))
}
