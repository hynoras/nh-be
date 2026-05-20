package procedure

import (
	"errors"
	"net/http"
	"nh-be/internal/utils/httputil"
)

const (
	ErrGetAllProceduresFailed    = "Failed to get procedures"
	ErrGetProcedureDetailFailed   = "Failed to get procedure detail"
	ErrGetProcedureStepsFailed    = "Failed to get procedure steps"
	ErrCreateProcedureFailed      = "Failed to create procedure"
	ErrUpdateProcedureFailed      = "Failed to update procedure"
	ErrUpdateProcedureStepFailed  = "Failed to update procedure step"
	ErrDeleteProcedureFailed      = "Failed to delete procedure"
	ErrAuthorizationFailed        = "Authorization failed"
)

var (
	ErrForbidViewProcedure   = errors.New("you do not have permission to view this procedure")
	ErrForbidCreateProcedure = errors.New("you do not have permission to create this procedure")
	ErrForbidUpdateProcedure = errors.New("you do not have permission to update this procedure")
	ErrForbidDeleteProcedure = errors.New("you do not have permission to delete this procedure")
	ErrProcedureNotFound      = errors.New("procedure not found")
	ErrProcedureAlreadyExists = errors.New("procedure already exists")
	ErrProcedureConflict      = errors.New("This procedure is modified by another request, please retry")
	ErrProcedureStepNotFound  = errors.New("procedure step not found")
	ErrProcedureStepConflict  = errors.New("This procedure step is modified by another request, please retry")

	// Proactive named errors
	ErrForbidManageProcedure = errors.New("you do not have permission to manage this procedure")
)

func init() {
	httputil.RegisterError(ErrForbidViewProcedure, http.StatusForbidden, ErrAuthorizationFailed)
	httputil.RegisterError(ErrForbidCreateProcedure, http.StatusForbidden, ErrAuthorizationFailed)
	httputil.RegisterError(ErrForbidUpdateProcedure, http.StatusForbidden, ErrAuthorizationFailed)
	httputil.RegisterError(ErrForbidDeleteProcedure, http.StatusForbidden, ErrAuthorizationFailed)
	httputil.RegisterError(ErrProcedureNotFound, http.StatusNotFound, "Procedure not found")
	httputil.RegisterError(ErrProcedureAlreadyExists, http.StatusConflict, "Procedure already exists")
	httputil.RegisterError(ErrProcedureConflict, http.StatusConflict, ErrUpdateProcedureFailed)
	httputil.RegisterError(ErrProcedureStepNotFound, http.StatusNotFound, "Procedure step not found")
	httputil.RegisterError(ErrProcedureStepConflict, http.StatusConflict, ErrUpdateProcedureStepFailed)
	httputil.RegisterError(ErrForbidManageProcedure, http.StatusForbidden, "you do not have permission to manage this procedure")
}
