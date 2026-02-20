package procedure

import (
	"time"

	"github.com/google/uuid"
)

type UsedByExperiment struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Objective string `json:"objective"`
}

type Steps struct {
	ID          string `json:"step_id"`
	ProcedureID string `json:"procedure_id"`

	Index       int    `json:"step_order"`
	Title       string `json:"title"`
	Description string `json:"description"`

	StepType string `json:"type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProcedureListResponseDto struct {
	ID                string             `json:"id"`
	Title             string             `json:"title"`
	Description       string             `json:"description"`
	UsedByExperiments []UsedByExperiment `json:"used_by_experiments"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

type ProcedureResponseDto struct {
	ID                string             `json:"id"`
	Title             string             `json:"title"`
	Description       string             `json:"description"`
	UsedByExperiments []UsedByExperiment `json:"used_by_experiments"`
	Steps             []Steps            `json:"steps"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

type CreateStepDto struct {
	Index       int    `json:"step_order" binding:"required"`
	Title       string `json:"title" binding:"required,min=3,max=200"`
	Description string `json:"description" binding:"omitempty,min=5,max=255"`
	StepType    string `json:"type" binding:"required,oneof=action wait decision observe cleanup"`
}

type CreateExperimentAssignmentDto struct {
	ID string `json:"id" binding:"required,uuid4"`
}

type CreateProcedureDto struct {
	Title                 string                          `json:"title" binding:"required,min=3,max=200"`
	Description           string                          `json:"description" binding:"omitempty,min=5,max=255"`
	Steps                 []CreateStepDto                 `json:"steps" binding:"omitempty,dive"`
	ExperimentAssignments []CreateExperimentAssignmentDto `json:"assigned_experiments" binding:"omitempty,dive"`
}

type UpdateProcedureDto struct {
	Title       string `json:"title" binding:"omitempty,gt=0,min=3,max=200"`
	Description string `json:"description" binding:"omitempty,min=5,max=255"`
	Version     int    `json:"version" binding:"required"`
}

type UpdateProcedureStepDto struct {
	ID          *string `json:"id" binding:"omitempty,uuid4"`
	Index       *int    `json:"step_order" binding:"omitempty"`
	Title       *string `json:"title" binding:"omitempty,min=3,max=200"`
	Description *string `json:"description" binding:"omitempty,min=5,max=255"`
	StepType    *string `json:"type" binding:"omitempty,oneof=action wait decision observe cleanup"`
	Version     int     `json:"version" binding:"required"`
}

type UpdateProcedureStepInput struct {
	ID          uuid.UUID
	Index       int
	Title       string
	Description string
	StepType    StepType
	Version     int
}
