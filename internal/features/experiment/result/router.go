package result

import (
	"nh-be/internal/features/permission"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	experimentsGroup := rg.Group("/experiments")
	experimentsGroup.Use(middleware.WithService("experiment-service"))

	// Setup shared dependencies
	resultRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionCache := permission.NewPermissionCache(rdb)
	permissionService := permission.NewService(permissionRepo, permissionCache)
	resultService := NewService(resultRepo, permissionService)

	// Result routes nested under /:experimentId/result
	experimentsGroup.GET("/:experimentId/result", GetResultByExperimentIDHandler(resultService))
	experimentsGroup.POST("/:experimentId/result", CreateResultHandler(resultService))
	experimentsGroup.PUT("/:experimentId/result/:resultId", UpdateResultHandler(resultService))
}
