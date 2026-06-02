package experiment

func MapExperimentToDto(e Experiment) ExperimentResponseDto {
	return ExperimentResponseDto{
		ID:          e.ID.String(),
		Title:       e.Title,
		Objective:   e.Objective,
		Status:      string(e.Status),
		Type:        string(e.Type),
		Version:     e.Version,
		CreatedBy:   e.CreatedBy.ID.String(),
		CreatedAt:   e.CreatedAt,
		StartedAt:   e.StartedAt,
		CompletedAt: e.CompletedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func MapExperimentsQueryToDto(experiments []ExperimentsQueryDto) []ExperimentsResponseDto {
	mappedExperiments := make([]ExperimentsResponseDto, len(experiments))
	for i, e := range experiments {
		mappedExperiments[i] = ExperimentsResponseDto{
			Identifier:    e.Identifier,
			Title:         e.Title,
			Objective:     e.Objective,
			Status:        string(e.Status),
			Type:          string(e.Type),
			Creator:       e.Creator,
			Updater:       e.Updater,
			CreatedAt:     e.CreatedAt,
			UpdatedAt:     e.UpdatedAt,
			ProcedureName: e.ProcedureName,
		}
	}
	return mappedExperiments
}
