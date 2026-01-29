package result

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func createValidExperimentResult() *ExperimentResult {
	return &ExperimentResult{
		ID:              uuid.New(),
		ExperimentID:    uuid.New(),
		Outcome:         OutcomeSuccess,
		Summary:         "Test summary for experiment result",
		OutcomeReason:   "Test outcome reason for the experiment",
		ConfidenceLevel: ConfidenceHigh,
		Version:         1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func TestRepository_Create_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	result := createValidExperimentResult()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "experiment_results"`).
		WithArgs(
			result.ID,
			result.ExperimentID,
			result.Outcome,
			result.Summary,
			result.OutcomeReason,
			result.ConfidenceLevel,
			result.Version,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, result)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRepository_Create_FKConstraintViolation(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	result := createValidExperimentResult()
	result.ExperimentID = uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "experiment_results"`).
		WithArgs(
			result.ID,
			result.ExperimentID,
			result.Outcome,
			result.Summary,
			result.OutcomeReason,
			result.ConfidenceLevel,
			result.Version,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("pq: insert or update on table \"experiment_results\" violates foreign key constraint \"experiment_results_experiment_id_fkey\""))
	mock.ExpectRollback()

	err := repo.Create(ctx, result)

	if err == nil || !strings.Contains(err.Error(), "foreign key constraint") {
		t.Errorf("expected FK constraint error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRepository_Create_DuplicateExperimentID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	result := createValidExperimentResult()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "experiment_results"`).
		WithArgs(
			result.ID,
			result.ExperimentID,
			result.Outcome,
			result.Summary,
			result.OutcomeReason,
			result.ConfidenceLevel,
			result.Version,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"experiment_results_experiment_id_key\""))
	mock.ExpectRollback()

	err := repo.Create(ctx, result)

	if err == nil || !strings.Contains(err.Error(), "duplicate key") {
		t.Errorf("expected duplicate key error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// ==================== UPDATE TESTS ====================

func TestRepository_Update_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	resultID := uuid.New()
	experimentID := uuid.New()
	currentVersion := 1

	result := &ExperimentResult{
		Outcome:         OutcomeFailure,
		Summary:         "Updated summary",
		OutcomeReason:   "Updated reason",
		ConfidenceLevel: ConfidenceMedium,
		UpdatedAt:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "experiment_results" SET`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, resultID, experimentID, result, currentVersion)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRepository_Update_OptimisticLockingConflict(t *testing.T) {
	cases := []struct {
		name string
	}{
		{"non-existent experiment"},
		{"non-existent result"},
		{"stale version"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			repo := NewRepository(db)
			ctx := context.Background()

			result := &ExperimentResult{
				Outcome:         OutcomeFailure,
				Summary:         "Updated summary",
				OutcomeReason:   "Updated reason",
				ConfidenceLevel: ConfidenceMedium,
				UpdatedAt:       time.Now(),
			}

			mock.ExpectBegin()
			mock.ExpectExec(`UPDATE "experiment_results" SET`).
				WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
			mock.ExpectCommit()

			err := repo.Update(ctx, uuid.New(), uuid.New(), result, 1)

			if !errors.Is(err, ErrOptimisticLockingConflict) {
				t.Errorf("expected ErrOptimisticLockingConflict, got: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}
