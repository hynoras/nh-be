package experiment

import (
	"nh-be/internal/features/permission"
	"nh-be/internal/features/procedure"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	experimentsGroup := rg.Group("/experiments")

	// Setup shared dependencies
	experimentRepo := NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionCache := permission.NewPermissionCache(rdb)
	permissionService := permission.NewService(permissionRepo, permissionCache)
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
