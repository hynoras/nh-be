package observation

import (
	"errors"
	"net/http"
	"nh-be/internal/utils/httputil"
)

const (
	ErrGetAllObservationsFailed = "Failed to get observations"
	ErrCreateObservationFailed  = "Failed to create observation"
)

var (
	ErrObservationNotFound     = errors.New("observation not found")
	ErrForbidViewObservation   = errors.New("you do not have permission to view observation")
	ErrForbidCreateObservation = errors.New("you do not have permission to create this observation")
)

func init() {
	httputil.RegisterError(ErrObservationNotFound, http.StatusNotFound, "Observation not found")
	httputil.RegisterError(ErrForbidViewObservation, http.StatusForbidden, ErrGetAllObservationsFailed)
	httputil.RegisterError(ErrForbidCreateObservation, http.StatusForbidden, ErrCreateObservationFailed)
}
