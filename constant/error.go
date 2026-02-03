package constant

import "errors"

const (
	ErrAuthorizationFailed = "Authorization failed"
)

var (
	ErrInvalidIDFormat           = errors.New("invalid id format")
	ErrProcedureNotFound         = errors.New("procedure not found")
	ErrProcedureAlreadyExists    = errors.New("procedure already exists")
	ErrForbidCreateProcedure     = errors.New("you do not have permission to create procedure")
	ErrForbidViewProcedure       = errors.New("you do not have permission to view this procedure")
	ErrForbidUpdateProcedure     = errors.New("you do not have permission to update this procedure")
	ErrForbidDeleteProcedure     = errors.New("you do not have permission to delete this procedure")
	ErrOptimisticLockingConflict = errors.New("the procedure was modified by another request, please retry")
)
