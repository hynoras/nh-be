package result

import (
	"context"
	"time"

	"nh-be/internal/constant"
	"nh-be/internal/features/experiment/result"

	"github.com/google/uuid"
)

func ContextWithUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), constant.CtxUserId, userID)
}

func CreateTestResult(experimentID uuid.UUID) *result.ExperimentResult {
	return &result.ExperimentResult{
		ID:              uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		ExperimentID:    experimentID,
		Outcome:         result.OutcomeSuccess,
		Summary:         "Test summary",
		OutcomeReason:   "Test reason",
		ConfidenceLevel: result.ConfidenceHigh,
		Version:         1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func CreateTestDto() *result.CreateResultDto {
	return &result.CreateResultDto{
		Outcome:         "success",
		Summary:         "Test summary for experiment",
		OutcomeReason:   "Test outcome reason",
		ConfidenceLevel: "high",
	}
}

func CreateInvalidTestDto() *result.CreateResultDto {
	return &result.CreateResultDto{
		Outcome:         "not a valid outcome",
		Summary:         "Test summary for experiment",
		OutcomeReason:   "Test outcome reason",
		ConfidenceLevel: "high",
	}
}

func CreateTestUpdateDto() *result.UpdateResultDto {
	outcome := "failure"
	summary := "Updated summary"
	return &result.UpdateResultDto{
		Version: 1,
		Outcome: &outcome,
		Summary: &summary,
	}
}

func CreateEmptyUpdateDto() *result.UpdateResultDto {
	return &result.UpdateResultDto{
		Version: 1,
	}
}

func CreateValidExperimentResult() *result.ExperimentResult {
	return &result.ExperimentResult{
		ID:              uuid.New(),
		ExperimentID:    uuid.New(),
		Outcome:         result.OutcomeSuccess,
		Summary:         "Test summary for experiment result",
		OutcomeReason:   "Test outcome reason for the experiment",
		ConfidenceLevel: result.ConfidenceHigh,
		Version:         1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func CreateUpdateFields() *result.UpdateFields {
	outcome := result.OutcomeFailure
	summary := "Updated summary"
	outcomeReason := "Updated reason"
	confidenceLevel := result.ConfidenceMedium

	return &result.UpdateFields{
		Outcome:         &outcome,
		Summary:         &summary,
		OutcomeReason:   &outcomeReason,
		ConfidenceLevel: &confidenceLevel,
	}
}
