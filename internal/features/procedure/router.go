package procedure

import (
	"nh-be/internal/features/permission"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	proceduresGroup := rg.Group("/procedures")

	// Setup shared dependencies
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	procedureRepo := NewRepository(db)
	procedureService := NewService(procedureRepo, permissionService)

	proceduresGroup.GET("", GetAllProceduresHandler(procedureService))
	proceduresGroup.GET("/:procedureId", GetProcedureByIDHandler(procedureService))
	proceduresGroup.POST("", CreateProcedureHandler(procedureService))
	proceduresGroup.PUT("/:procedureId", UpdateProcedureHandler(procedureService))
	// proceduresGroup.PUT("/:procedureId/status", UpdateProcedureStatusHandler(experimentService))
	// proceduresGroup.DELETE("/:procedureId", DeleteProcedureHandler(experimentService))

}
