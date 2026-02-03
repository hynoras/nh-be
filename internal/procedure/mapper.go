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
