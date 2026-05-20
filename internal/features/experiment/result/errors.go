package result

import (
	"errors"
	"net/http"
	"nh-be/internal/utils/httputil"
)

const (
	ErrGetExperimentResultDetailFailed = "Failed to get experiment result"
	ErrCreateExperimentResultFailed    = "Failed to create experiment result"
	ErrUpdateExperimentResultFailed    = "Failed to update experiment result"
)

var (
	ErrExperimentResultNotFound      = errors.New("experiment result not found")
	ErrExperimentResultAlreadyExists = errors.New("experiment result already exists for this experiment")
	ErrForbidCreateExperimentResult  = errors.New("you do not have permission to create experiment results")
	ErrForbidViewExperimentResult    = errors.New("you do not have permission to view this experiment result")
	ErrForbidUpdateExperimentResult  = errors.New("you do not have permission to update this experiment result")
	ErrInvalidOutcome                = errors.New("invalid outcome value")
	ErrInvalidConfidenceLevel        = errors.New("invalid confidence level value")
	ErrExperimentResultConflict      = errors.New("the experiment result was modified by another request, please retry")
)

func init() {
	httputil.RegisterError(ErrExperimentResultNotFound, http.StatusNotFound, "Experiment result not found")
	httputil.RegisterError(ErrExperimentResultAlreadyExists, http.StatusConflict, "Experiment result already exists")
	httputil.RegisterError(ErrForbidCreateExperimentResult, http.StatusForbidden, ErrCreateExperimentResultFailed)
	httputil.RegisterError(ErrForbidViewExperimentResult, http.StatusForbidden, ErrGetExperimentResultDetailFailed)
	httputil.RegisterError(ErrForbidUpdateExperimentResult, http.StatusForbidden, ErrUpdateExperimentResultFailed)
	httputil.RegisterError(ErrInvalidOutcome, http.StatusBadRequest, "Invalid outcome value")
	httputil.RegisterError(ErrInvalidConfidenceLevel, http.StatusBadRequest, "Invalid confidence level value")
	httputil.RegisterError(ErrExperimentResultConflict, http.StatusConflict, ErrUpdateExperimentResultFailed)
}
