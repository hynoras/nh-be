package procedure

type StepType string

const (
	StepTypeAction   StepType = "action"
	StepTypeWait     StepType = "wait"
	StepTypeDecision StepType = "decision"
	StepTypeObserve  StepType = "observe"
	StepTypeCleanup  StepType = "cleanup"
)
