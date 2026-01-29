package result

import (
	"nh-be/internal/experiment/root"
	"nh-be/internal/permission"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	resultsGroup := rg.Group("/experiment-results")

	resultRepo := NewRepository(db)
	experimentRepo := root.NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo)
	resultService := NewService(resultRepo, experimentRepo, permissionService)

	resultsGroup.GET("/:experimentId", GetResultByExperimentIDHandler(resultService))
	resultsGroup.POST("", CreateResultHandler(resultService))
	resultsGroup.PUT("/:id/experiment/:experimentId", UpdateResultHandler(resultService))
}
