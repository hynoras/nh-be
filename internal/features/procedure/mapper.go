package procedure

import (
	"nh-be/internal/utils/timeutil"
	"time"
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
			StepType:    s.StepType,
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
		Steps:       MapStepsToDto(p.Steps),
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
			StepType:    step.StepType,
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
