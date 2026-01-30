package result

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func createUpdateFields() *UpdateFields {
	outcome := OutcomeFailure
	summary := "Updated summary"
	outcomeReason := "Updated reason"
	confidenceLevel := ConfidenceMedium

	return &UpdateFields{
		Outcome:         &outcome,
		Summary:         &summary,
		OutcomeReason:   &outcomeReason,
		ConfidenceLevel: &confidenceLevel,
	}
}

func TestRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, result *ExperimentResult)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock, result *ExperimentResult) {
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
			},
			expectedError: nil,
		},
		{
			name: "fk_constraint_violation",
			setupMock: func(mock sqlmock.Sqlmock, result *ExperimentResult) {
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
					WillReturnError(errors.New(`pq: insert or update on table "experiment_results" violates foreign key constraint "experiment_results_experiment_id_fkey"`))
				mock.ExpectRollback()
			},
			expectedError: ErrExperimentNotFound,
		},
		{
			name: "duplicate_experiment_id",
			setupMock: func(mock sqlmock.Sqlmock, result *ExperimentResult) {
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
					WillReturnError(errors.New(`pq: duplicate key value violates unique constraint "experiment_results_experiment_id_key"`))
				mock.ExpectRollback()
			},
			expectedError: ErrExperimentResultAlreadyExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			repo := NewRepository(db)
			ctx := context.Background()

			result := createValidExperimentResult()

			tc.setupMock(mock, result)

			err := repo.Create(ctx, result)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_Update(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiment_results" SET`).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "optimistic_lock_non_existent_experiment",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiment_results" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
				mock.ExpectCommit()
			},
			expectedError: ErrOptimisticLockingConflict,
		},
		{
			name: "optimistic_lock_non_existent_result",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiment_results" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
				mock.ExpectCommit()
			},
			expectedError: ErrOptimisticLockingConflict,
		},
		{
			name: "optimistic_lock_stale_version",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiment_results" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
				mock.ExpectCommit()
			},
			expectedError: ErrOptimisticLockingConflict,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			repo := NewRepository(db)
			ctx := context.Background()

			resultID := uuid.New()
			experimentID := uuid.New()
			currentVersion := 1
			fields := createUpdateFields()

			tc.setupMock(mock)

			err := repo.Update(ctx, resultID, experimentID, fields, currentVersion)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
