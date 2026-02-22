package result

import (
	"context"
	"errors"
	"testing"

	"nh-be/internal/constant"
	"nh-be/internal/features/experiment/result"
	"nh-be/internal/utils/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, res *result.ExperimentResult)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock, res *result.ExperimentResult) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "experiment_results"`).
					WithArgs(
						res.ID,
						res.ExperimentID,
						res.Outcome,
						res.Summary,
						res.OutcomeReason,
						res.ConfidenceLevel,
						res.Version,
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
			setupMock: func(mock sqlmock.Sqlmock, res *result.ExperimentResult) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "experiment_results"`).
					WithArgs(
						res.ID,
						res.ExperimentID,
						res.Outcome,
						res.Summary,
						res.OutcomeReason,
						res.ConfidenceLevel,
						res.Version,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnError(errors.New(`pq: insert or update on table "experiment_results" violates foreign key constraint "experiment_results_experiment_id_fkey"`))
				mock.ExpectRollback()
			},
			expectedError: constant.ErrExperimentNotFound,
		},
		{
			name: "duplicate_experiment_id",
			setupMock: func(mock sqlmock.Sqlmock, res *result.ExperimentResult) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "experiment_results"`).
					WithArgs(
						res.ID,
						res.ExperimentID,
						res.Outcome,
						res.Summary,
						res.OutcomeReason,
						res.ConfidenceLevel,
						res.Version,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnError(errors.New(`pq: duplicate key value violates unique constraint "experiment_results_experiment_id_key"`))
				mock.ExpectRollback()
			},
			expectedError: constant.ErrExperimentResultAlreadyExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := result.NewRepository(db)
			ctx := context.Background()

			res := CreateValidExperimentResult()

			tc.setupMock(mock, res)

			err := repo.Create(ctx, res)

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
			expectedError: constant.ErrExperimentResultConflict,
		},
		{
			name: "optimistic_lock_non_existent_result",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiment_results" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
				mock.ExpectCommit()
			},
			expectedError: constant.ErrExperimentResultConflict,
		},
		{
			name: "optimistic_lock_stale_version",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiment_results" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
				mock.ExpectCommit()
			},
			expectedError: constant.ErrExperimentResultConflict,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := result.NewRepository(db)
			ctx := context.Background()

			resultID := uuid.New()
			experimentID := uuid.New()
			currentVersion := 1
			fields := CreateUpdateFields()

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
