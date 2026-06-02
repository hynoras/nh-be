package experiment

import (
	"errors"
	"net/http"
	"nh-be/internal/utils/httputil"
)

const (
	ErrGetAllExperimentFailed            = "Failed to get experiments"
	ErrGetExperimentDetailFailed         = "Failed to get experiment detail"
	ErrCreateExperimentFailed            = "Failed to create experiment"
	ErrUpdateExperimentFailed            = "Failed to update experiment"
	ErrDeleteExperimentFailed            = "Failed to delete experiment"
	ErrInvalidStatusTransition           = "Invalid status transition"
	ErrAssignProcedureToExperimentFailed = "Failed to assign procedure to experiment"
	ErrAuthorizationFailed               = "Authorization failed"
	ErrInvalidExperimentStatus           = "Invalid experiment status"
	ErrInvalidExperimentType             = "Invalid experiment type"
)

var (
	ErrExperimentNotFound                              = errors.New("experiment not found")
	ErrForbidCreateExperiment                          = errors.New("you do not have permission to create experiments")
	ErrForbidViewExperiments                           = errors.New("you do not have permission to view experiments")
	ErrForbidViewExperiment                            = errors.New("you do not have permission to view this experiment")
	ErrForbidUpdateExperiment                          = errors.New("you do not have permission to update this experiment")
	ErrForbidDeleteExperiment                          = errors.New("you do not have permission to delete this experiment")
	ErrForbidAssignProcedureToExperiment               = errors.New("you do not have permission to assign procedure to experiment")
	ErrStatusTransitionFromDraftToPlanning             = errors.New("Invalid status transition, only draft can be transition to planning")
	ErrStatusTransitionFromPlanningToRunning           = errors.New("Invalid status transition, only planning can be transition to running")
	ErrStatusTransitionFromRunningToCompletedOrAborted = errors.New("Invalid status transition, only running can be transition to completed or aborted")
	ErrExperimentConflict                              = errors.New("the experiment was modified by another request, please retry")
	ErrExperimentAlreadyInTargetState                  = errors.New("experiment is already in target state")
	ErrDuplicateProcedureAssignment                    = errors.New("This procedure is already assigned to the experiment")
	ErrMustBeOneOfExperimentStatus                     = errors.New("Status must be one of: draft, planning, running, completed, aborted")
	ErrMustBeOneOfExperimentType                       = errors.New("Type must be one of: exploratory, confirmatory")
	ErrForbidManageExperiment                          = errors.New("you do not have permission to manage this experiment")
)

func init() {
	httputil.RegisterError(ErrForbidViewExperiments, http.StatusForbidden, ErrGetAllExperimentFailed)
	httputil.RegisterError(ErrForbidViewExperiment, http.StatusForbidden, ErrGetExperimentDetailFailed)
	httputil.RegisterError(ErrForbidUpdateExperiment, http.StatusForbidden, ErrUpdateExperimentFailed)
	httputil.RegisterError(ErrForbidDeleteExperiment, http.StatusForbidden, ErrAuthorizationFailed)
	httputil.RegisterError(ErrStatusTransitionFromDraftToPlanning, http.StatusBadRequest, ErrInvalidStatusTransition)
	httputil.RegisterError(ErrStatusTransitionFromPlanningToRunning, http.StatusBadRequest, ErrInvalidStatusTransition)
	httputil.RegisterError(ErrStatusTransitionFromRunningToCompletedOrAborted, http.StatusBadRequest, ErrInvalidStatusTransition)
	httputil.RegisterError(ErrExperimentConflict, http.StatusConflict, ErrUpdateExperimentFailed)
	httputil.RegisterError(ErrExperimentAlreadyInTargetState, http.StatusBadRequest, ErrUpdateExperimentFailed)
	httputil.RegisterError(ErrExperimentNotFound, http.StatusNotFound, "Experiment not found")
	httputil.RegisterError(ErrDuplicateProcedureAssignment, http.StatusConflict, ErrAssignProcedureToExperimentFailed)
	httputil.RegisterError(ErrForbidManageExperiment, http.StatusForbidden, "you do not have permission to manage this experiment")
	httputil.RegisterError(ErrMustBeOneOfExperimentStatus, http.StatusBadRequest, ErrInvalidExperimentStatus)
	httputil.RegisterError(ErrMustBeOneOfExperimentType, http.StatusBadRequest, ErrInvalidExperimentType)
}
