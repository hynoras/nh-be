package procedure

func MapProcedureListToDto(p Procedure) ProcedureListResponseDto {
	return ProcedureListResponseDto{
		ID:                p.ID.String(),
		Title:             p.Title,
		Description:       p.Description,
		UsedByExperiments: MapUsedByExperimentsToDto(p.Experiments),
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

func MapUsedByExperimentsToDto(experiments []ProcedureExperimentAssignment) []UsedByExperiment {
	result := []UsedByExperiment{}
	for _, e := range experiments {
		result = append(result, UsedByExperiment{
			ID:        e.ExperimentID.String(),
			Title:     e.Experiment.Title,
			Objective: e.Experiment.Objective,
		})
	}
	return result
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
			UpdatedAt:   s.UpdatedAt,
		})
	}
	return result
}

func MapProcedureToDto(p *Procedure) ProcedureResponseDto {
	return ProcedureResponseDto{
		ID:                p.ID.String(),
		Title:             p.Title,
		Description:       p.Description,
		UsedByExperiments: MapUsedByExperimentsToDto(p.Experiments),
		Steps:             MapStepsToDto(p.Steps),
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

func MapProceduresToDto(procedures []Procedure) []ProcedureListResponseDto {
	result := []ProcedureListResponseDto{}
	for _, p := range procedures {
		result = append(result, ProcedureListResponseDto{
			ID:                p.ID.String(),
			Title:             p.Title,
			Description:       p.Description,
			UsedByExperiments: MapUsedByExperimentsToDto(p.Experiments),
			CreatedAt:         p.CreatedAt,
			UpdatedAt:         p.UpdatedAt,
		})
	}
	return result
}
