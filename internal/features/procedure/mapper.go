package procedure

import (
	"nh-be/internal/utils/timeutil"
	"time"

	"github.com/google/uuid"
)

func MapProcedureListToDto(p Procedure) ProcedureListResponseDto {
	return ProcedureListResponseDto{
		ID:          p.ID.String(),
		Title:       p.Title,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   *p.UpdatedAt,
	}
}

func MapStepsToDto(step []ProcedureStep) []Steps {
	result := []Steps{}
	for _, s := range step {
		result = append(result, Steps{
			ID:          s.ID.String(),
			Title:       s.Title,
			Description: s.Description,
			Index:       s.Index,
			StepType:    string(s.StepType),
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   *s.UpdatedAt,
		})
	}
	return result
}

func MapProcedureToDto(p *Procedure) ProcedureResponseDto {
	return ProcedureResponseDto{
		ID:          p.ID.String(),
		Title:       p.Title,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   *p.UpdatedAt,
	}
}

func MapProceduresToDto(procedures []Procedure) []ProcedureListResponseDto {
	result := []ProcedureListResponseDto{}
	for _, p := range procedures {
		result = append(result, ProcedureListResponseDto{
			ID:          p.ID.String(),
			Title:       p.Title,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   *p.UpdatedAt,
		})
	}
	return result
}

func MapCreateDtoToProcedure(p *CreateProcedureDto) *Procedure {
	return &Procedure{
		Title:       p.Title,
		Description: p.Description,
		Steps:       MapCreateDtoToProcedureStep(p.Steps),
	}
}

func MapCreateDtoToProcedureStep(s []CreateStepDto) []ProcedureStep {
	result := make([]ProcedureStep, 0, len(s))
	for _, step := range s {
		result = append(result, ProcedureStep{
			Title:       step.Title,
			Description: step.Description,
			Index:       step.Index,
			StepType:    StepType(step.StepType),
		})
	}
	return result
}

func MapUpdateDtoToProcedure(p *UpdateProcedureDto) *Procedure {
	return &Procedure{
		Title:       p.Title,
		Description: p.Description,
		UpdatedAt:   timeutil.TimePtr(time.Now()),
		Version:     p.Version,
	}
}

func MapUpdateProcStepDtoToProcStepInput(step *UpdateProcedureStepDto) *UpdateProcedureStepInput {
	input := &UpdateProcedureStepInput{
		Version: step.Version,
	}
	if step.ID != nil {
		input.ID = uuid.MustParse(*step.ID)
	}
	if step.Title != nil {
		input.Title = *step.Title
	}
	if step.Description != nil {
		input.Description = *step.Description
	}
	if step.Index != nil {
		input.Index = *step.Index
	}
	if step.StepType != nil {
		input.StepType = StepType(*step.StepType)
	}
	return input
}

func MapUpdateProcStepDtoToProcStepInputs(step []UpdateProcedureStepDto) []UpdateProcedureStepInput {
	result := make([]UpdateProcedureStepInput, 0, len(step))
	for _, s := range step {
		result = append(result, *MapUpdateProcStepDtoToProcStepInput(&s))
	}
	return result
}

func MapCreateProcStepInputToProcStep(step *UpdateProcedureStepInput) *ProcedureStep {
	return &ProcedureStep{
		Title:       step.Title,
		Description: step.Description,
		Index:       step.Index,
		StepType:    StepType(step.StepType),
	}
}

func MapUpdateProcStepInputToProcStep(step *UpdateProcedureStepInput, updatedAt *time.Time) *ProcedureStep {
	return &ProcedureStep{
		Title:       step.Title,
		Description: step.Description,
		Index:       step.Index,
		StepType:    step.StepType,
		UpdatedAt:   updatedAt,
		Version:     step.Version,
	}
}

func MapUpdateProcStepInputsToProcSteps(steps []UpdateProcedureStepInput, updatedAt *time.Time) []ProcedureStep {
	result := make([]ProcedureStep, 0, len(steps))
	for _, s := range steps {
		result = append(result, *MapUpdateProcStepInputToProcStep(&s, updatedAt))
	}
	return result
}
