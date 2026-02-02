package experiment

import (
	"nh-be/internal/experiment/result"
	"nh-be/internal/experiment/root"
	"nh-be/internal/permission"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	experimentsGroup := rg.Group("/experiments")

	// Setup shared dependencies
	experimentRepo := root.NewRepository(db)
	resultRepo := result.NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	experimentService := root.NewService(experimentRepo, permissionService)
	resultService := result.NewService(resultRepo, permissionService)

	// Experiment routes
	experimentsGroup.GET("", root.GetAllExperimentsHandler(experimentService))
	experimentsGroup.GET("/:experimentId", root.GetExperimentByIDHandler(experimentService))
	experimentsGroup.POST("", root.CreateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:experimentId", root.UpdateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:experimentId/status", root.UpdateExperimentStatusHandler(experimentService))
	experimentsGroup.DELETE("/:experimentId", root.DeleteExperimentHandler(experimentService))

	// Result routes nested under /:experimentId/result
	experimentsGroup.GET("/:experimentId/result", result.GetResultByExperimentIDHandler(resultService))
	experimentsGroup.POST("/:experimentId/result", result.CreateResultHandler(resultService))
	experimentsGroup.PUT("/:experimentId/result/:resultId", result.UpdateResultHandler(resultService))
}
