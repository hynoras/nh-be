package root

func MapExperimentToDto(e Experiment) ExperimentResponseDto {
	return ExperimentResponseDto{
		ID:          e.ID.String(),
		Title:       e.Title,
		Objective:   e.Objective,
		Status:      string(e.Status),
		Type:        string(e.Type),
		CreatedBy:   e.CreatedBy.ID.String(),
		CreatedAt:   e.CreatedAt,
		StartedAt:   e.StartedAt,
		CompletedAt: e.CompletedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func MapExperimentsToDto(experiments []Experiment) []ExperimentsResponseDto {
	var result []ExperimentsResponseDto
	for _, e := range experiments {
		result = append(result, ExperimentsResponseDto{
			ID:        e.ID.String(),
			Title:     e.Title,
			Objective: e.Objective,
			Status:    string(e.Status),
			Type:      string(e.Type),
			CreatedBy: e.CreatedBy.ID.String(),
			CreatedAt: e.CreatedAt,
		})
	}
	return result
}
