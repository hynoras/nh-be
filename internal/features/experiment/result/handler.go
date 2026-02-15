package result

import (
	"net/http"
	"nh-be/internal/utils/httputil"

	"github.com/gin-gonic/gin"
)

// GetResultByExperimentIDHandler godoc
// @Summary Get experiment result by experiment ID
// @Description Retrieve the result of a specific experiment by its experiment UUID
// @Tags Experiment Results
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Success 200 {object} httputil.SuccessResponse{data=ExperimentResultResponseDto} "Experiment result fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid experiment ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Experiment or result not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get experiment result"
// @Security SessionAuth
// @Router /experiments/{experimentId}/result [get]
func GetResultByExperimentIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		experimentID, idErr := httputil.ValidateUUID(c, c.Param("experimentId"))
		if idErr != nil {
			return
		}

		result, serviceErr := s.GetResultByExperimentID(c.Request.Context(), *experimentID)
		if HandleServiceError(c, serviceErr) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "Experiment result fetched successfully", MapResultToDto(*result))
	}
}

// CreateResultHandler godoc
// @Summary Create a new experiment result
// @Description Create a result for an experiment
// @Tags Experiment Results
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Param request body CreateResultDto true "Experiment result creation details"
// @Success 201 {object} httputil.SuccessResponse "Experiment result created successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Experiment not found"
// @Failure 409 {object} httputil.ErrorResponse "Experiment result already exists"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to create experiment result"
// @Security SessionAuth
// @Router /experiments/{experimentId}/result [post]
func CreateResultHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		experimentID, idErr := httputil.ValidateUUID(c, c.Param("experimentId"))
		if idErr != nil {
			return
		}

		var dto CreateResultDto
		if err := httputil.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		serviceErr := s.CreateResult(c.Request.Context(), *experimentID, &dto)
		if HandleServiceError(c, serviceErr) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusCreated, "Experiment result created successfully", nil)
	}
}

// UpdateResultHandler godoc
// @Summary Update an experiment result
// @Description Update an existing experiment result by result ID and experiment ID
// @Tags Experiment Results
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Param resultId path string true "Result ID (UUID format)"
// @Param request body UpdateResultDto true "Updated experiment result details (includes version for optimistic locking)"
// @Success 200 {object} httputil.SuccessResponse "Experiment result updated successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Experiment or result not found"
// @Failure 409 {object} httputil.ErrorResponse "Optimistic locking conflict"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to update experiment result"
// @Security SessionAuth
// @Router /experiments/{experimentId}/result/{resultId} [put]
func UpdateResultHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		experimentID, expIdErr := httputil.ValidateUUID(c, c.Param("experimentId"))
		if expIdErr != nil {
			return
		}

		resultID, idErr := httputil.ValidateUUID(c, c.Param("resultId"))
		if idErr != nil {
			return
		}

		var dto UpdateResultDto
		if err := httputil.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		serviceErr := s.UpdateResult(c.Request.Context(), *resultID, *experimentID, &dto)
		if HandleServiceError(c, serviceErr) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "Experiment result updated successfully", nil)
	}
}
