package observation

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/observation"
	"nh-be/internal/utils/testutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetAllObsByExpIDAndProcID(t *testing.T) {
	expId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	procId := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name              string
		expId             uuid.UUID
		procId            uuid.UUID
		offset            int
		limit             int
		sortBy            *string
		sortOrder         *constant.Order
		ctx               func() context.Context
		setupMock         func(mock sqlmock.Sqlmock)
		expectedError     error
		checkError        func(t *testing.T, err error)
		expectedCount     int64
		expectedDataCount int
	}{
		{
			name:   "experiment_not_found",
			expId:  expId,
			procId: procId,
			offset: 0,
			limit:  10,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedError: constant.ErrExperimentNotFound,
		},
		{
			name:   "procedure_not_found",
			expId:  expId,
			procId: procId,
			offset: 0,
			limit:  10,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedError: constant.ErrProcedureNotFound,
		},
		{
			name:   "database_error_on_experiment_check",
			expId:  expId,
			procId: procId,
			offset: 0,
			limit:  10,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnError(assert.AnError)
			},
			expectedError: assert.AnError,
		},
		{
			name:   "database_error_on_procedure_check",
			expId:  expId,
			procId: procId,
			offset: 0,
			limit:  10,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnError(assert.AnError)
			},
			expectedError: assert.AnError,
		},
		{
			name:   "empty_result_returns_zero_count",
			expId:  expId,
			procId: procId,
			offset: 0,
			limit:  10,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(procId))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2`).
					WithArgs(procId, expId).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectQuery(`SELECT .+ FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2 LIMIT \$3`).
					WithArgs(procId, expId, 10).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "observed_at", "title", "notes", "previous_observation_id", "created_by", "created_at",
					}))
			},
			expectedError:     nil,
			expectedCount:     0,
			expectedDataCount: 0,
		},
		{
			name:   "returns_observations_with_correct_count",
			expId:  expId,
			procId: procId,
			offset: 0,
			limit:  10,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				observations := TestObservationMetadataList(3)

				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(procId))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2`).
					WithArgs(procId, expId).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

				rows := sqlmock.NewRows([]string{
					"id", "observed_at", "title", "notes", "previous_observation_id", "created_by", "created_at",
				})
				for _, o := range observations {
					rows.AddRow(o.ID, o.ObservedAt, o.Title, o.Notes, o.PreviousObservationID, o.CreatedBy, o.CreatedAt)
				}

				mock.ExpectQuery(`SELECT .+ FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2 LIMIT \$3`).
					WithArgs(procId, expId, 10).
					WillReturnRows(rows)
			},
			expectedError:     nil,
			expectedCount:     3,
			expectedDataCount: 3,
		},
		{
			name:   "pagination_applied_correctly",
			expId:  expId,
			procId: procId,
			offset: 2,
			limit:  1,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				observations := TestObservationMetadataList(1)

				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(procId))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2`).
					WithArgs(procId, expId).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

				rows := sqlmock.NewRows([]string{
					"id", "observed_at", "title", "notes", "previous_observation_id", "created_by", "created_at",
				})
				for _, o := range observations {
					rows.AddRow(o.ID, o.ObservedAt, o.Title, o.Notes, o.PreviousObservationID, o.CreatedBy, o.CreatedAt)
				}

				mock.ExpectQuery(`SELECT .+ FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2 LIMIT \$3 OFFSET \$4`).
					WithArgs(procId, expId, 1, 1).
					WillReturnRows(rows)
			},
			expectedError:     nil,
			expectedCount:     5,
			expectedDataCount: 1,
		},
		{
			name:      "sort_by_created_at_asc",
			expId:     expId,
			procId:    procId,
			offset:    0,
			limit:     10,
			sortBy:    func() *string { s := "created_at"; return &s }(),
			sortOrder: func() *constant.Order { o := constant.ASC; return &o }(),
			ctx:       func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				observations := TestObservationMetadataList(2)

				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(procId))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2`).
					WithArgs(procId, expId).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				rows := sqlmock.NewRows([]string{
					"id", "observed_at", "title", "notes", "previous_observation_id", "created_by", "created_at",
				})
				for _, o := range observations {
					rows.AddRow(o.ID, o.ObservedAt, o.Title, o.Notes, o.PreviousObservationID, o.CreatedBy, o.CreatedAt)
				}

				mock.ExpectQuery(`SELECT .+ FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2 ORDER BY created_at ASC LIMIT \$3`).
					WithArgs(procId, expId, 10).
					WillReturnRows(rows)
			},
			expectedError:     nil,
			expectedCount:     2,
			expectedDataCount: 2,
		},
		{
			name:      "sort_by_created_at_desc",
			expId:     expId,
			procId:    procId,
			offset:    0,
			limit:     10,
			sortBy:    func() *string { s := "created_at"; return &s }(),
			sortOrder: func() *constant.Order { o := constant.DESC; return &o }(),
			ctx:       func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				observations := TestObservationMetadataList(2)

				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(procId))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2`).
					WithArgs(procId, expId).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				rows := sqlmock.NewRows([]string{
					"id", "observed_at", "title", "notes", "previous_observation_id", "created_by", "created_at",
				})
				for _, o := range observations {
					rows.AddRow(o.ID, o.ObservedAt, o.Title, o.Notes, o.PreviousObservationID, o.CreatedBy, o.CreatedAt)
				}

				mock.ExpectQuery(`SELECT .+ FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2 ORDER BY created_at DESC LIMIT \$3`).
					WithArgs(procId, expId, 10).
					WillReturnRows(rows)
			},
			expectedError:     nil,
			expectedCount:     2,
			expectedDataCount: 2,
		},
		{
			name:      "invalid_sort_column_ignored",
			expId:     expId,
			procId:    procId,
			offset:    0,
			limit:     10,
			sortBy:    func() *string { s := "title"; return &s }(),
			sortOrder: func() *constant.Order { o := constant.ASC; return &o }(),
			ctx:       func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				observations := TestObservationMetadataList(1)

				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(procId))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2`).
					WithArgs(procId, expId).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{
					"id", "observed_at", "title", "notes", "previous_observation_id", "created_by", "created_at",
				})
				for _, o := range observations {
					rows.AddRow(o.ID, o.ObservedAt, o.Title, o.Notes, o.PreviousObservationID, o.CreatedBy, o.CreatedAt)
				}

				// No ORDER BY clause since "title" is not in allowed columns
				mock.ExpectQuery(`SELECT .+ FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2 LIMIT \$3`).
					WithArgs(procId, expId, 10).
					WillReturnRows(rows)
			},
			expectedError:     nil,
			expectedCount:     1,
			expectedDataCount: 1,
		},
		{
			name:   "count_query_error",
			expId:  expId,
			procId: procId,
			offset: 0,
			limit:  10,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(procId))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2`).
					WithArgs(procId, expId).
					WillReturnError(assert.AnError)
			},
			expectedError: assert.AnError,
		},
		{
			name:   "find_query_error",
			expId:  expId,
			procId: procId,
			offset: 0,
			limit:  10,
			ctx:    func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT "id" FROM "experiments" WHERE id = \$1`).
					WithArgs(expId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expId))

				mock.ExpectQuery(`SELECT "id" FROM "procedure_steps" WHERE id = \$1`).
					WithArgs(procId, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(procId))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2`).
					WithArgs(procId, expId).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				mock.ExpectQuery(`SELECT .+ FROM "observations" WHERE procedure_step_id = \$1 AND experiment_id = \$2 LIMIT \$3`).
					WithArgs(procId, expId, 10).
					WillReturnError(assert.AnError)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			tc.setupMock(mock)
			repo := observation.NewRepository(db)
			ctx := tc.ctx()

			result, count, err := repo.GetAllObsByExpIDAndProcID(
				ctx, tc.expId, tc.procId,
				tc.offset, tc.limit,
				tc.sortBy, tc.sortOrder,
			)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tc.expectedDataCount)
				assert.Equal(t, tc.expectedCount, count)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateObservation(t *testing.T) {
	expId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	procId := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, obs *observation.Observation)
		expectedError error
	}{
		{
			name: "success",
			setupMock: func(mock sqlmock.Sqlmock, obs *observation.Observation) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "observations"`).
					WithArgs(
						obs.ObservedAt,
						obs.Title,
						obs.Notes,
						obs.CreatedBy,
						obs.ExperimentID,
						obs.ProcedureStepID,
						obs.ID,
						obs.CreatedAt,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(obs.ID, obs.CreatedAt))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "experiment_not_found",
			setupMock: func(mock sqlmock.Sqlmock, obs *observation.Observation) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "observations"`).
					WithArgs(
						obs.ObservedAt,
						obs.Title,
						obs.Notes,
						obs.CreatedBy,
						obs.ExperimentID,
						obs.ProcedureStepID,
						obs.ID,
						obs.CreatedAt,
					).
					WillReturnError(&pgconn.PgError{
						Code:           "23503",
						ConstraintName: "observations_experiment_id_fkey",
					})
				mock.ExpectRollback()
			},
			expectedError: constant.ErrExperimentNotFound,
		},
		{
			name: "procedure_step_not_found",
			setupMock: func(mock sqlmock.Sqlmock, obs *observation.Observation) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "observations"`).
					WithArgs(
						obs.ObservedAt,
						obs.Title,
						obs.Notes,
						obs.CreatedBy,
						obs.ExperimentID,
						obs.ProcedureStepID,
						obs.ID,
						obs.CreatedAt,
					).
					WillReturnError(&pgconn.PgError{
						Code:           "23503",
						ConstraintName: "observations_procedure_step_id_fkey",
					})
				mock.ExpectRollback()
			},
			expectedError: constant.ErrProcedureNotFound,
		},
		{
			name: "database_error",
			setupMock: func(mock sqlmock.Sqlmock, obs *observation.Observation) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "observations"`).
					WithArgs(
						obs.ObservedAt,
						obs.Title,
						obs.Notes,
						obs.CreatedBy,
						obs.ExperimentID,
						obs.ProcedureStepID,
						obs.ID,
						obs.CreatedAt,
					).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			expectedError: assert.AnError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := observation.NewRepository(db)
			ctx := context.Background()

			obs := TestObservation()

			tc.setupMock(mock, &obs)

			result, err := repo.CreateObservation(ctx, expId, procId, obs)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Equal(t, observation.Observation{}, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, obs, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
