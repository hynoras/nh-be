package root

import "errors"

type ExperimentStatus string

const (
	ExperimentDraft     ExperimentStatus = "draft"
	ExperimentPlanning  ExperimentStatus = "planning"
	ExperimentRunning   ExperimentStatus = "running"
	ExperimentCompleted ExperimentStatus = "completed"
	ExperimentAborted   ExperimentStatus = "aborted"
)

var (
	ErrExperimentNotFound     = errors.New("experiment not found")
	ErrForbidCreateExperiment = errors.New("you do not have permission to create experiments")
	ErrForbidViewExperiments  = errors.New("you do not have permission to view experiments")
	ErrForbidViewExperiment   = errors.New("you do not have permission to view this experiment")
	ErrForbidUpdateExperiment = errors.New("you do not have permission to update this experiment")
	ErrForbidDeleteExperiment = errors.New("you do not have permission to delete this experiment")
)
