package procedure

import (
	"net/http"
	"nh-be/utils"

	"github.com/gin-gonic/gin"
)

// GetAllProceduresHandler godoc
// @Summary Get all procedures
// @Description Retrieve a paginated list of procedures with optional search filter
// @Tags Procedures
// @Accept json
// @Produce json
// @Param search query string false "Search term to filter procedures by title"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} utils.SuccessResponse{data=[]ProcedureListResponseDto} "Procedures fetched successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid pagination parameters"
// @Failure 403 {object} utils.ErrorResponse "Authorization failed"
// @Failure 500 {object} utils.ErrorResponse "Failed to get procedures"
// @Security SessionAuth
// @Router /procedures [get]
func GetAllProceduresHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")

		pageInt, pageSizeInt, err := utils.ParsePaginationParams(c)
		if err != nil {
			return
		}

		procedures, length, serviceErr := s.GetAllProcedures(c.Request.Context(), search, pageInt, pageSizeInt)
		if utils.MakeServiceErrorResponse(c, serviceErr) {
			return
		}

		utils.MakeSuccessResponse(c, http.StatusOK, "Procedures fetched successfully", procedures, length)
	}
}
