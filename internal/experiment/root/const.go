package root

import "errors"

type ExperimentStatus string
type ExperimentType string

const (
	ExperimentDraft     ExperimentStatus = "draft"
	ExperimentPlanning  ExperimentStatus = "planning"
	ExperimentRunning   ExperimentStatus = "running"
	ExperimentCompleted ExperimentStatus = "completed"
	ExperimentAborted   ExperimentStatus = "aborted"
)

const (
	ExperimentExploratoryType  ExperimentType = "exploratory"
	ExperimentConfirmatoryType ExperimentType = "confirmatory"
)

var (
	ErrExperimentNotFound                              = errors.New("experiment not found")
	ErrForbidCreateExperiment                          = errors.New("you do not have permission to create experiments")
	ErrForbidViewExperiments                           = errors.New("you do not have permission to view experiments")
	ErrForbidViewExperiment                            = errors.New("you do not have permission to view this experiment")
	ErrForbidUpdateExperiment                          = errors.New("you do not have permission to update this experiment")
	ErrForbidDeleteExperiment                          = errors.New("you do not have permission to delete this experiment")
	ErrStatusTransitionFromDraftToPlanning             = errors.New("Invalid status transition, only draft can be transition to planning")
	ErrStatusTransitionFromPlanningToRunning           = errors.New("Invalid status transition, only planning can be transition to running")
	ErrStatusTransitionFromRunningToCompletedOrAborted = errors.New("Invalid status transition, only running can be transition to completed or aborted")
)
