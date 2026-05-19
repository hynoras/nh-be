package result

import (
	"nh-be/internal/app"
	"nh-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, deps *app.SharedDeps) {
	experimentsGroup := rg.Group("/experiments")
	experimentsGroup.Use(middleware.WithService("experiment-service"))

	// Setup shared dependencies
	resultRepo := NewRepository(deps.DB)
	resultService := NewService(resultRepo, deps.PermissionService)

	// Result routes nested under /:experimentId/result
	experimentsGroup.GET("/:experimentId/result", GetResultByExperimentIDHandler(resultService))
	experimentsGroup.POST("/:experimentId/result", CreateResultHandler(resultService))
	experimentsGroup.PUT("/:experimentId/result/:resultId", UpdateResultHandler(resultService))
}
