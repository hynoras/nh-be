package root

import (
	"nh-be/internal/permission"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	experimentsGroup := rg.Group("/experiments")

	experimentRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	experimentService := NewService(experimentRepo, permissionService)

	experimentsGroup.GET("", GetAllExperimentsHandler(experimentService))
	experimentsGroup.GET("/:id", GetExperimentByIDHandler(experimentService))
	experimentsGroup.POST("", CreateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:id", UpdateExperimentHandler(experimentService))
	experimentsGroup.PUT("/:id/status", UpdateExperimentStatusHandler(experimentService))
	experimentsGroup.DELETE("/:id", DeleteExperimentHandler(experimentService))
}
