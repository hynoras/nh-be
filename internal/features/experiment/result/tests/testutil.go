package result

import (
	"context"
	"testing"
	"time"

	"nh-be/internal/constant"
	"nh-be/internal/features/experiment/result"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func SetupTestRouter(method, path string, handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(method, path, handler)
	return router
}

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	return gormDB, mock
}
