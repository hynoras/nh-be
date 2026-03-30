package procedure

import (
	"nh-be/internal/features/permission"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	proceduresGroup := rg.Group("/procedures")

	// Setup shared dependencies
	permissionRepo := permission.NewRepository(db)
	permissionCache := permission.NewPermissionCache(rdb)
	permissionService := permission.NewService(permissionRepo, permissionCache)
	procedureRepo := NewRepository(db)
	procedureService := NewService(procedureRepo, permissionService)

	proceduresGroup.GET("", GetAllProceduresHandler(procedureService))
	proceduresGroup.GET("/:procedureId", GetProcedureByIDHandler(procedureService))
	proceduresGroup.POST("", CreateProcedureHandler(procedureService))
	proceduresGroup.PUT("/:procedureId", UpdateProcedureHandler(procedureService))
	proceduresGroup.GET("/:procedureId/procedure-steps", GetProcedureStepsHandler(procedureService))
	proceduresGroup.PUT("/:procedureId/procedure-steps", UpdateProcedureStepHandler(procedureService))
	// proceduresGroup.PUT("/:procedureId/status", UpdateProcedureStatusHandler(experimentService))
	proceduresGroup.DELETE("/:procedureId", DeleteProcedureHandler(procedureService))

}
