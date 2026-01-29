package result

func MapResultToDto(r ExperimentResult) ExperimentResultResponseDto {
	return ExperimentResultResponseDto{
		ID:              r.ID.String(),
		ExperimentID:    r.ExperimentID.String(),
		Outcome:         string(r.Outcome),
		Summary:         r.Summary,
		OutcomeReason:   r.OutcomeReason,
		ConfidenceLevel: string(r.ConfidenceLevel),
		Version:         r.Version,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
