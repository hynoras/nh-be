package result

import (
	"net/http"
	"nh-be/constant"
	"nh-be/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetResultByExperimentIDHandler godoc
// @Summary Get experiment result by experiment ID
// @Description Retrieve the result of a specific experiment by its experiment UUID
// @Tags Experiment Results
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Success 200 {object} utils.SuccessResponse{data=ExperimentResultResponseDto} "Experiment result fetched successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid experiment ID format"
// @Failure 403 {object} utils.ErrorResponse "Authorization failed"
// @Failure 404 {object} utils.ErrorResponse "Experiment or result not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to get experiment result"
// @Security SessionAuth
// @Router /experiment-results/experiment/{experimentId} [get]
func GetResultByExperimentIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := utils.ValidateUUID(c, c.Param("experimentId"))
		if idErr != nil {
			return
		}

		result, serviceErr := s.GetResultByExperimentID(c.Request.Context(), *parsedId)
		switch serviceErr {
		case ErrForbidViewExperimentResult:
			utils.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "Experiment result not found", serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusOK, "Experiment result fetched successfully", MapResultToDto(*result))
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get experiment result", serviceErr.Error())
			return
		}
	}
}

// CreateResultHandler godoc
// @Summary Create a new experiment result
// @Description Create a result for an experiment
// @Tags Experiment Results
// @Accept json
// @Produce json
// @Param request body CreateResultDto true "Experiment result creation details"
// @Success 201 {object} utils.SuccessResponse "Experiment result created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 403 {object} utils.ErrorResponse "Authorization failed"
// @Failure 404 {object} utils.ErrorResponse "Experiment not found"
// @Failure 409 {object} utils.ErrorResponse "Experiment result already exists"
// @Failure 422 {object} utils.ErrorResponse "Validation failed"
// @Failure 500 {object} utils.ErrorResponse "Failed to create experiment result"
// @Security SessionAuth
// @Router /experiment-results [post]
func CreateResultHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreateResultDto
		if err := utils.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		serviceErr := s.CreateResult(c.Request.Context(), &dto)
		switch serviceErr {
		case ErrForbidCreateExperimentResult:
			utils.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case ErrExperimentResultAlreadyExists:
			utils.MakeErrorResponse(c, http.StatusConflict, "Experiment result already exists", serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "Experiment not found", serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusCreated, "Experiment result created successfully", nil)
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to create experiment result", serviceErr.Error())
			return
		}
	}
}

// UpdateResultHandler godoc
// @Summary Update an experiment result
// @Description Update an existing experiment result by result ID and experiment ID
// @Tags Experiment Results
// @Accept json
// @Produce json
// @Param id path string true "Result ID (UUID format)"
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Param request body UpdateResultDto true "Updated experiment result details (includes version for optimistic locking)"
// @Success 200 {object} utils.SuccessResponse "Experiment result updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID format"
// @Failure 403 {object} utils.ErrorResponse "Authorization failed"
// @Failure 404 {object} utils.ErrorResponse "Experiment or result not found"
// @Failure 409 {object} utils.ErrorResponse "Optimistic locking conflict"
// @Failure 422 {object} utils.ErrorResponse "Validation failed"
// @Failure 500 {object} utils.ErrorResponse "Failed to update experiment result"
// @Security SessionAuth
// @Router /experiment-results/{id}/experiment/{experimentId} [put]
func UpdateResultHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedResultId, idErr := utils.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		parsedExperimentId, expIdErr := utils.ValidateUUID(c, c.Param("experimentId"))
		if expIdErr != nil {
			return
		}

		var dto UpdateResultDto
		if err := utils.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		serviceErr := s.UpdateResult(c.Request.Context(), *parsedResultId, *parsedExperimentId, &dto)
		switch serviceErr {
		case ErrForbidUpdateExperimentResult:
			utils.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "Experiment result not found", serviceErr.Error())
			return
		case ErrOptimisticLockingConflict:
			utils.MakeErrorResponse(c, http.StatusConflict, "Version conflict", serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusOK, "Experiment result updated successfully", nil)
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update experiment result", serviceErr.Error())
			return
		}
	}
}
