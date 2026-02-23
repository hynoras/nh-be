package procedure

import (
	"net/http"
	"nh-be/internal/constant"
	"nh-be/internal/utils/httputil"

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
// @Success 200 {object} httputil.SuccessResponse{data=[]ProcedureListResponseDto} "Procedures fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid pagination parameters"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get procedures"
// @Security SessionAuth
// @Router /procedures [get]
func GetAllProceduresHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")

		pageInt, pageSizeInt, err := httputil.ParsePaginationParams(c)
		if err != nil {
			return
		}

		procedures, length, serviceErr := s.GetAllProcedures(c.Request.Context(), search, pageInt, pageSizeInt)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrGetAllProceduresFailed) {
			return
		}

		httputil.MakeSuccessResponse(c, http.StatusOK, "Procedures fetched successfully", procedures, length)
	}
}

// GetProcedureByIDHandler godoc
// @Summary Get procedure by ID
// @Description Retrieve a procedure by ID
// @Tags Procedures
// @Accept json
// @Produce json
// @Param procedureId path string true "Procedure ID"
// @Success 200 {object} httputil.SuccessResponse{data=ProcedureResponseDto} "Procedure fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid procedure ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Procedure not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get procedure"
// @Security SessionAuth
// @Router /procedures/:procedureId [get]
func GetProcedureByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		procedureId, idErr := httputil.ValidateUUID(c, c.Param("procedureId"))
		if idErr != nil {
			return
		}
		procedure, serviceErr := s.GetProcedureByID(c.Request.Context(), *procedureId)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrGetProcedureDetailFailed) {
			return
		}

		httputil.MakeSuccessResponse(c, http.StatusOK, "Procedures fetched successfully", procedure)
	}
}

// CreateProcedureHandler godoc
// @Summary Create a new procedure
// @Description Create a new procedure
// @Tags Procedures
// @Accept json
// @Produce json
// @Param request body CreateProcedureDto true "Procedure creation details"
// @Success 201 {object} httputil.SuccessResponse "Procedure created successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to create procedure"
// @Security SessionAuth
// @Router /procedures [post]
func CreateProcedureHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreateProcedureDto
		if err := httputil.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		serviceErr := s.CreateProcedure(c.Request.Context(), &dto)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrCreateProcedureFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusCreated, "Procedure created successfully", nil)
	}
}

// UpdateProcedureHandler godoc
// @Summary Update a procedure
// @Description Update a procedure
// @Tags Procedures
// @Accept json
// @Produce json
// @Param request body UpdateProcedureDto true "Procedure update details"
// @Success 200 {object} httputil.SuccessResponse "Procedure updated successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Procedure not found"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to update procedure"
// @Security SessionAuth
// @Router /procedures/:procedureId [put]
func UpdateProcedureHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		procedureId, idErr := httputil.ValidateUUID(c, c.Param("procedureId"))
		if idErr != nil {
			return
		}

		var dto UpdateProcedureDto
		if err := httputil.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		serviceErr := s.UpdateProcedure(c.Request.Context(), *procedureId, &dto)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrUpdateProcedureFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "Procedure updated successfully", nil)
	}
}

// DeleteProcedureHandler godoc
// @Summary Delete a procedure
// @Description Delete a procedure
// @Tags Procedures
// @Accept json
// @Produce json
// @Param procedureId path string true "Procedure ID"
// @Success 200 {object} httputil.SuccessResponse "Procedure deleted successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Procedure not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to delete procedure"
// @Security SessionAuth
// @Router /procedures/:procedureId [delete]
func DeleteProcedureHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		procedureId, idErr := httputil.ValidateUUID(c, c.Param("procedureId"))
		if idErr != nil {
			return
		}

		serviceErr := s.DeleteProcedure(c.Request.Context(), *procedureId)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrDeleteProcedureFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "Procedure deleted successfully", nil)
	}
}

// GetProcedureStepsHandler godoc
// @Summary Get procedure steps
// @Description Get procedure steps
// @Tags Procedures
// @Accept json
// @Produce json
// @Param procedureId path string true "Procedure ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} httputil.SuccessResponse "Procedure steps fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Procedure not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get procedure steps"
// @Security SessionAuth
// @Router /procedures/:procedureId/procedure-steps [get]
func GetProcedureStepsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		procedureId, idErr := httputil.ValidateUUID(c, c.Param("procedureId"))
		if idErr != nil {
			return
		}

		pageInt, pageSizeInt, err := httputil.ParsePaginationParams(c)
		if err != nil {
			return
		}

		steps, length, serviceErr := s.GetProcedureSteps(c.Request.Context(), *procedureId, pageInt, pageSizeInt)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrGetProcedureStepsFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "Procedure steps fetched successfully", steps, length)
	}
}

// UpdateProcedureStepHandler godoc
// @Summary Update procedure steps
// @Description Update procedure steps (create, update, delete)
// @Tags Procedures
// @Accept json
// @Produce json
// @Param procedureId path string true "Procedure ID"
// @Param request body []UpdateProcedureStepDto true "Procedure steps update details"
// @Success 200 {object} httputil.SuccessResponse "Procedure steps updated successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Procedure not found"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to update procedure steps"
// @Security SessionAuth
// @Router /procedures/:procedureId/procedure-steps [put]
func UpdateProcedureStepHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		procedureId, idErr := httputil.ValidateUUID(c, c.Param("procedureId"))
		if idErr != nil {
			return
		}

		var dto []UpdateProcedureStepDto
		if err := httputil.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		cleanInput := MapUpdateProcStepDtoToProcStepInputs(dto)

		serviceErr := s.UpdateProcedureStep(c.Request.Context(), *procedureId, cleanInput)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrUpdateProcedureFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "Procedure updated successfully", nil)
	}
}
