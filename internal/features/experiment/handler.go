package experiment

import (
	"net/http"
	"nh-be/constant"
	"nh-be/internal/httputil"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllExperimentsHandler godoc
// @Summary Get all experiments
// @Description Retrieve a paginated list of experiments with optional search filter
// @Tags Experiments
// @Accept json
// @Produce json
// @Param search query string false "Search term to filter experiments by title"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} httputil.SuccessResponse{data=[]ExperimentsResponseDto} "Experiments fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid pagination parameters"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get experiments"
// @Security SessionAuth
// @Router /experiments [get]
func GetAllExperimentsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")

		pageInt, pageSizeInt, err := httputil.ParsePaginationParams(c)
		if err != nil {
			return
		}

		experiments, length, serviceErr := s.GetAllExperiments(c.Request.Context(), search, pageInt, pageSizeInt)
		switch serviceErr {
		case ErrForbidViewExperiments:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case nil:
			experimentResp := MapExperimentsToDto(experiments)
			httputil.MakeSuccessResponse(c, http.StatusOK, "Experiments fetched successfully", experimentResp, length)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get experiments", serviceErr.Error())
			return
		}
	}
}

// GetExperimentByIDHandler godoc
// @Summary Get experiment by ID
// @Description Retrieve a single experiment by its UUID
// @Tags Experiments
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Success 200 {object} httputil.SuccessResponse{data=ExperimentResponseDto} "Experiment fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Experiment not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get experiment"
// @Security SessionAuth
// @Router /experiments/{experimentId} [get]
func GetExperimentByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("experimentId"))
		if idErr != nil {
			return
		}

		experiment, serviceErr := s.GetExperimentByID(c.Request.Context(), *parsedId)
		switch serviceErr {
		case ErrForbidViewExperiment:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "Experiment not found", serviceErr.Error())
			return
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusOK, "Experiment fetched successfully", MapExperimentToDto(*experiment))
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get experiment", serviceErr.Error())
			return
		}
	}
}

// CreateExperimentHandler godoc
// @Summary Create a new experiment
// @Description Create a new experiment with title and objective
// @Tags Experiments
// @Accept json
// @Produce json
// @Param request body CreateExperimentDto true "Experiment creation details"
// @Success 201 {object} httputil.SuccessResponse "Experiment created successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to create experiment"
// @Security SessionAuth
// @Router /experiments [post]
func CreateExperimentHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreateExperimentDto
		if err := httputil.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		serviceErr := s.CreateExperiment(c.Request.Context(), &dto)
		switch serviceErr {
		case ErrForbidCreateExperiment:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusCreated, "Experiment created successfully", nil)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to create experiment", serviceErr.Error())
			return
		}
	}
}

// UpdateExperimentHandler godoc
// @Summary Update an experiment
// @Description Update an existing experiment by its UUID
// @Tags Experiments
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Param request body UpdateExperimentDto true "Updated experiment details"
// @Success 200 {object} httputil.SuccessResponse "Experiment updated successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid experiment ID"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Experiment not found"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to update experiment"
// @Security SessionAuth
// @Router /experiments/{experimentId} [put]
func UpdateExperimentHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("experimentId"))
		if idErr != nil {
			return
		}

		var dto UpdateExperimentDto
		if err := httputil.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		serviceErr := s.UpdateExperiment(c.Request.Context(), *parsedId, &dto)
		switch serviceErr {
		case ErrForbidUpdateExperiment:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "Experiment not found", serviceErr.Error())
			return
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusOK, "Experiment updated successfully", nil)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update experiment", serviceErr.Error())
			return
		}
	}
}

// UpdateExperimentStatusHandler godoc
// @Summary Update experiment status
// @Description Update the status of an existing experiment by its UUID
// @Tags Experiments
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Param request body UpdateExperimentStatusDto true "New experiment status"
// @Success 200 {object} httputil.SuccessResponse "Experiment status updated successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid experiment ID or status transition"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Experiment not found"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to update experiment status"
// @Security SessionAuth
// @Router /experiments/{experimentId}/status [patch]
func UpdateExperimentStatusHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("experimentId"))
		if idErr != nil {
			return
		}

		var dto UpdateExperimentStatusDto
		if err := httputil.ValidateRequestFormat(c, &dto); err != nil {
			return
		}

		status := ExperimentStatus(dto.Status)

		serviceErr := s.UpdateExperimentStatus(c.Request.Context(), *parsedId, status)
		switch serviceErr {
		case ErrForbidUpdateExperiment:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "Experiment not found", nil)
			return
		case ErrAlreadyInTargetState:
			httputil.MakeErrorResponse(c, http.StatusOK, "Experiment is already in target state", nil)
			return
		//for some reasons, stacked case is not working here so i use multiple case
		case ErrStatusTransitionFromDraftToPlanning:
			httputil.MakeErrorResponse(c, http.StatusConflict, "Invalid status transition", serviceErr.Error())
			return
		case ErrStatusTransitionFromPlanningToRunning:
			httputil.MakeErrorResponse(c, http.StatusConflict, "Invalid status transition", serviceErr.Error())
			return
		case ErrStatusTransitionFromRunningToCompletedOrAborted:
			httputil.MakeErrorResponse(c, http.StatusConflict, "Invalid status transition", serviceErr.Error())
			return
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusOK, "Experiment status updated successfully", nil)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update experiment status", serviceErr.Error())
			return
		}
	}
}

// DeleteExperimentHandler godoc
// @Summary Delete an experiment
// @Description Delete a single experiment by its UUID
// @Tags Experiments
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID (UUID format)"
// @Success 200 {object} httputil.SuccessResponse "Experiment deleted successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid experiment ID"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Experiment not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to delete experiment"
// @Security SessionAuth
// @Router /experiments/{experimentId} [delete]
func DeleteExperimentHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("experimentId"))
		if idErr != nil {
			return
		}

		serviceErr := s.DeleteExperiment(c.Request.Context(), *parsedId)
		switch serviceErr {
		case ErrForbidDeleteExperiment:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "Experiment not found", serviceErr.Error())
			return
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusOK, "Experiment deleted successfully", nil)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to delete experiment", serviceErr.Error())
			return
		}
	}
}
