package observation

import (
	"nh-be/internal/app"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, deps *app.SharedDeps) {
	observationsGroup := rg.Group("/observations")
	observationsGroup.Use(middleware.WithService("observation-service"))

	// Setup shared dependencies
	observationRepo := NewRepository(deps.DB)
	observationService := NewService(observationRepo, deps.PermissionService)

	// Observation routes
	observationsGroup.GET("/:experimentId/:procedureStepId", GetAllObservationsHandler(observationService))
	observationsGroup.POST("/:experimentId/:procedureStepId", CreateObservationHandler(observationService))
}
