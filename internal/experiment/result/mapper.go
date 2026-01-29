package result

func MapResultToDto(r ExperimentResult) ExperimentResultResponseDto {
	return ExperimentResultResponseDto{
		ID:              r.ID.String(),
		ExperimentID:    r.ExperimentID.String(),
		Outcome:         string(r.Outcome),
		Summary:         r.Summary,
		OutcomeReason:   r.OutcomeReason,
		ConfidenceLevel: string(r.ConfidenceLevel),
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
