package observation

import (
	"nh-be/internal/features/permission"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	observationsGroup := rg.Group("/observations")

	// Setup shared dependencies
	observationRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	observationService := NewService(observationRepo, permissionService)

	// Observation routes
	observationsGroup.GET("/:experimentId/:procedureStepId", GetAllObservationsHandler(observationService))
}
