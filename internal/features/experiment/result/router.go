package result

import (
	"nh-be/internal/features/permission"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	experimentsGroup := rg.Group("/experiments")

	// Setup shared dependencies
	resultRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)

	permissionService := permission.NewService(permissionRepo)
	resultService := NewService(resultRepo, permissionService)

	// Result routes nested under /:experimentId/result
	experimentsGroup.GET("/:experimentId/result", GetResultByExperimentIDHandler(resultService))
	experimentsGroup.POST("/:experimentId/result", CreateResultHandler(resultService))
	experimentsGroup.PUT("/:experimentId/result/:resultId", UpdateResultHandler(resultService))
}
