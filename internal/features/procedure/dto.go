package procedure

import (
	"time"

	"github.com/google/uuid"
)

type StepsResponseDto struct {
	ID          string    `json:"step_id"`
	Index       int       `json:"step_order"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	IsOptional  bool      `json:"is_optional"`
	StepType    string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProcedureListResponseDto struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProcedureResponseDto struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Steps       []StepsResponseDto `json:"steps"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type CreateStepDto struct {
	Index       int     `json:"step_order" binding:"required"`
	Title       string  `json:"title" binding:"required,min=3,max=200"`
	Description *string `json:"description" binding:"omitempty,min=5,max=255"`
	IsOptional  bool    `json:"is_optional"`
	StepType    string  `json:"type" binding:"required,oneof=action wait decision observe cleanup"`
	WaitTime    *int    `json:"wait_time" binding:"omitempty,gt=0"`
}

type CreateProcedureDto struct {
	Title       string          `json:"title" binding:"required,min=3,max=200"`
	Description *string         `json:"description" binding:"omitempty,min=5,max=255"`
	Steps       []CreateStepDto `json:"steps" binding:"omitempty,dive"`
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
	IsOptional  *bool   `json:"is_optional" binding:"omitempty"`
	WaitTime    *int    `json:"wait_time" binding:"omitempty,gt=0"`
	Version     int     `json:"version" binding:"required"`
}

type UpdateProcedureStepInput struct {
	ID          uuid.UUID
	Index       int
	Title       string
	Description string
	StepType    StepType
	IsOptional  bool
	WaitTime    int
	Version     int
}
