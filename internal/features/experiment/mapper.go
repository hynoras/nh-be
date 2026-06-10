package experiment

func MapExperimentToDto(e Experiment) ExperimentResponseDto {
	dto := ExperimentResponseDto{
		ID:             e.ID.String(),
		Identifier:     e.Identifier,
		Title:          e.Title,
		Objective:      e.Objective,
		Status:         string(e.Status),
		Type:           string(e.Type),
		Version:        e.Version,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		StartedAt:      e.StartedAt,
		CompletedAt:    e.CompletedAt,
		CreatedByID:    e.CreatedByID.String(),
		CreatedBy:      e.CreatedBy.ID.String(),
		PlannedStartAt: e.PlannedStartAt,
		PlannedEndAt:   e.PlannedEndAt,
	}

	if e.UpdatedByID != nil {
		dto.UpdatedByID = e.UpdatedByID.String()
	}
	if e.UpdatedBy != nil {
		dto.UpdatedBy = e.UpdatedBy.ID.String()
	}

	return dto
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
