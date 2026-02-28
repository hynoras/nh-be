package observation

import (
	"net/http"
	"nh-be/internal/constant"
	"nh-be/internal/utils/httputil"

	"github.com/gin-gonic/gin"
)

// GetAllObservationsHandler godoc
// @Summary Get all observations
// @Description Retrieve a paginated list of observations for a specific experiment and procedure step
// @Tags Observation
// @Accept json
// @Produce json
// @Param experimentId path string true "Experiment ID"
// @Param procedureStepId path string true "Procedure Step ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param sortBy query string false "Sort by"
// @Param sortOrder query string false "Sort order"
// @Success 200 {object} httputil.SuccessResponse{data=[]ObservationsResponseDto} "Observations fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid pagination or sorting parameters"
// @Failure 404 {object} httputil.ErrorResponse "Experiment or procedure not found"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get observations"
// @Security BearerAuth
// @Router /observations/{experimentId}/{procedureStepId} [get]
func GetAllObservationsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedExpId, expIdErr := httputil.ValidateUUID(c, c.Param("experimentId"))
		if expIdErr != nil {
			return
		}

		parsedProcStepId, procStepIdErr := httputil.ValidateUUID(c, c.Param("procedureStepId"))
		if procStepIdErr != nil {
			return
		}

		pageInt, pageSizeInt, err := httputil.ParsePaginationParams(c)
		if err != nil {
			return
		}

		sortBy, sortOrder, err := httputil.ParseSortParams(c)
		if err != nil {
			return
		}

		observations, length, serviceErr := s.GetAllObservations(
			c.Request.Context(),
			*parsedExpId,
			*parsedProcStepId,
			pageInt,
			pageSizeInt,
			&sortBy,
			&sortOrder,
		)

		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrGetAllObservationsFailed) {
			return
		}

		httputil.MakeSuccessResponse(c, http.StatusOK, "Observations fetched successfully", observations, length)
	}
}
