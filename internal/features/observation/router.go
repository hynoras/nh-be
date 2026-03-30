package observation

import (
	"nh-be/internal/features/permission"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	observationsGroup := rg.Group("/observations")

	// Setup shared dependencies
	observationRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionCache := permission.NewPermissionCache(rdb)
	permissionService := permission.NewService(permissionRepo, permissionCache)
	observationService := NewService(observationRepo, permissionService)

	// Observation routes
	observationsGroup.GET("/:experimentId/:procedureStepId", GetAllObservationsHandler(observationService))
	observationsGroup.POST("/:experimentId/:procedureStepId", CreateObservationHandler(observationService))
}
