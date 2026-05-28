package procedure

import (
	"nh-be/internal/app"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, deps *app.SharedDeps) {
	proceduresGroup := rg.Group("/procedures")
	proceduresGroup.Use(middleware.WithService("procedure-service"))

	// Setup shared dependencies
	procedureRepo := NewRepository(deps.DB)
	procedureService := NewService(procedureRepo, deps.PermissionService)

	proceduresGroup.GET("", GetAllProceduresHandler(procedureService))
	proceduresGroup.GET("/:procedureId", GetProcedureByIDHandler(procedureService))
	proceduresGroup.POST("", CreateProcedureHandler(procedureService))
	proceduresGroup.PUT("/:procedureId", UpdateProcedureHandler(procedureService))
	proceduresGroup.GET("/:procedureId/procedure-steps", GetProcedureStepsHandler(procedureService))
	proceduresGroup.PUT("/:procedureId/procedure-steps", UpdateProcedureStepHandler(procedureService))
	// proceduresGroup.PUT("/:procedureId/status", UpdateProcedureStatusHandler(experimentService))
	proceduresGroup.DELETE("/:procedureId", DeleteProcedureHandler(procedureService))
}
