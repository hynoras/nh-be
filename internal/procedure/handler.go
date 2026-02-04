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

// GetProcedureByIDHandler godoc
// @Summary Get procedure by ID
// @Description Retrieve a procedure by ID
// @Tags Procedures
// @Accept json
// @Produce json
// @Param procedureId path string true "Procedure ID"
// @Success 200 {object} utils.SuccessResponse{data=ProcedureResponseDto} "Procedure fetched successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid procedure ID format"
// @Failure 403 {object} utils.ErrorResponse "Authorization failed"
// @Failure 404 {object} utils.ErrorResponse "Procedure not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to get procedure"
// @Security SessionAuth
// @Router /procedures/:procedureId [get]
func GetProcedureByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		procedureId, idErr := utils.ValidateUUID(c, c.Param("procedureId"))
		if idErr != nil {
			return
		}
		procedure, serviceErr := s.GetProcedureByID(c.Request.Context(), *procedureId)
		if utils.MakeServiceErrorResponse(c, serviceErr) {
			return
		}

		utils.MakeSuccessResponse(c, http.StatusOK, "Procedures fetched successfully", procedure)
	}
}
