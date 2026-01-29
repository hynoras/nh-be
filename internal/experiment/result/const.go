package result

import "errors"

var (
	ErrExperimentResultNotFound      = errors.New("experiment result not found")
	ErrExperimentResultAlreadyExists = errors.New("experiment result already exists for this experiment")
	ErrForbidCreateExperimentResult  = errors.New("you do not have permission to create experiment results")
	ErrForbidViewExperimentResult    = errors.New("you do not have permission to view this experiment result")
	ErrForbidUpdateExperimentResult  = errors.New("you do not have permission to update this experiment result")
	ErrInvalidOutcome                = errors.New("invalid outcome value")
	ErrInvalidConfidenceLevel        = errors.New("invalid confidence level value")
	ErrOptimisticLockingConflict     = errors.New("the experiment result was modified by another request, please retry")
)
