package experiment

import (
	"nh-be/internal/features/permission"
	"nh-be/internal/features/procedure"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	experimentsGroup := rg.Group("/experiments")

	// Setup shared dependencies
	experimentRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	procedureRepo := procedure.NewRepository(db)
	procedureService := procedure.NewService(procedureRepo, permissionService)
	experimentService := NewService(experimentRepo, permissionService, procedureService)

	// Experiment routes
	experimentsGroup.GET("", GetAllExperimentsHandler(experimentService))
	experimentsGroup.GET("/:experimentId", GetExperimentByIDHandler(experimentService))
	experimentsGroup.POST("", CreateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:experimentId", UpdateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:experimentId/status", UpdateExperimentStatusHandler(experimentService))
	experimentsGroup.PUT("/:experimentId/procedures/:procedureId", AssignProcedureToExperimentHandler(experimentService))
	experimentsGroup.DELETE("/:experimentId", DeleteExperimentHandler(experimentService))
}
