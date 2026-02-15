package result

import (
	"net/http"
	"nh-be/constant"
	"nh-be/internal/utils/httputil"

	"github.com/gin-gonic/gin"
)

// HandleServiceError maps domain errors to HTTP responses.
// Returns true if error was handled, false if nil error.
func HandleServiceError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	switch err {
	// Permission errors → 403
	case ErrForbidViewExperimentResult,
		ErrForbidCreateExperimentResult,
		ErrForbidUpdateExperimentResult:
		httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, err.Error())

	// Not found errors → 404
	case ErrExperimentResultNotFound:
		httputil.MakeErrorResponse(c, http.StatusNotFound, "Experiment result not found", err.Error())
	case ErrExperimentNotFound:
		httputil.MakeErrorResponse(c, http.StatusNotFound, "Experiment not found", err.Error())

	// Conflict errors → 409
	case ErrExperimentResultAlreadyExists:
		httputil.MakeErrorResponse(c, http.StatusConflict, "Experiment result already exists", err.Error())
	case ErrOptimisticLockingConflict:
		httputil.MakeErrorResponse(c, http.StatusConflict, "Version conflict", err.Error())

	// Unknown errors → 500
	default:
		httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Internal server error", err.Error())
	}
	return true
}
