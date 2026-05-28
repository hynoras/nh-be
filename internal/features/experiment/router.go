package experiment

import (
	"nh-be/internal/app"
	"nh-be/internal/features/procedure"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, deps *app.SharedDeps) {
	experimentsGroup := rg.Group("/experiments")
	experimentsGroup.Use(middleware.WithService("experiment-service"))

	// Setup shared dependencies
	experimentRepo := NewRepository(deps.DB)
	procedureRepo := procedure.NewRepository(deps.DB)
	procedureService := procedure.NewService(procedureRepo, deps.PermissionService)
	experimentService := NewService(experimentRepo, deps.PermissionService, procedureService)

	// Experiment routes
	experimentsGroup.GET("", GetAllExperimentsHandler(experimentService))
	experimentsGroup.GET("/:experimentId", GetExperimentByIDHandler(experimentService))
	experimentsGroup.POST("", CreateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:experimentId", UpdateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:experimentId/status", UpdateExperimentStatusHandler(experimentService))
	experimentsGroup.PUT("/:experimentId/procedures/:procedureId", AssignProcedureToExperimentHandler(experimentService))
	experimentsGroup.DELETE("/:experimentId", DeleteExperimentHandler(experimentService))
}
