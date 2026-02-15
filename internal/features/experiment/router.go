package experiment

import (
	"nh-be/internal/features/permission"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	experimentsGroup := rg.Group("/experiments")

	// Setup shared dependencies
	experimentRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	experimentService := NewService(experimentRepo, permissionService)

	// Experiment routes
	experimentsGroup.GET("", GetAllExperimentsHandler(experimentService))
	experimentsGroup.GET("/:experimentId", GetExperimentByIDHandler(experimentService))
	experimentsGroup.POST("", CreateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:experimentId", UpdateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:experimentId/status", UpdateExperimentStatusHandler(experimentService))
	experimentsGroup.DELETE("/:experimentId", DeleteExperimentHandler(experimentService))
}
