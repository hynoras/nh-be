package experiment

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
