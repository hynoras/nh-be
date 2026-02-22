package procedure

import (
	"context"
	"nh-be/internal/constant"
	"nh-be/internal/features/procedure"
	"nh-be/internal/utils/stringutil"
	"nh-be/internal/utils/testutil"
	"nh-be/internal/utils/timeutil"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type QueryParams struct {
	Limit  int
	Offset int
	Search string
}

func TestRepository_FindAll(t *testing.T) {
	tests := []struct {
		name                    string
		queryParams             QueryParams
		ctx                     func() context.Context
		setupMock               func(mock sqlmock.Sqlmock)
		expectedError           error
		checkError              func(t *testing.T, err error)
		expectedProceduresCount int
		expectedTotalCount      int64
	}{
		{
			name: "success",
			queryParams: QueryParams{
				Limit:  10,
				Offset: 0,
				Search: "",
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				procedures := TestProcedureList()

				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				})

				for _, p := range procedures {
					rows.AddRow(
						p.ID,
						p.Title,
						p.Description,
						p.Version,
						p.ParentID,
						p.CreatedAt,
						p.UpdatedAt,
					)
				}

				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures"`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				mock.ExpectQuery(`SELECT \* FROM "procedures" LIMIT \$1`).
					WithArgs(10).
					WillReturnRows(rows)

			},
			expectedError:           nil,
			checkError:              nil,
			expectedProceduresCount: 2,
			expectedTotalCount:      2,
		},
		{
			name: "success with search",
			queryParams: QueryParams{
				Limit:  10,
				Offset: 0,
				Search: "Test",
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				procedures := TestProcedureList()

				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				})

				for _, p := range procedures {
					rows.AddRow(
						p.ID,
						p.Title,
						p.Description,
						p.Version,
						p.ParentID,
						p.CreatedAt,
						p.UpdatedAt,
					)
				}

				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures" WHERE LOWER\(procedures.title\) LIKE \$1`).
					WithArgs("%test%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE LOWER\(procedures.title\) LIKE \$1 LIMIT \$2`).
					WithArgs("%test%", 10).
					WillReturnRows(rows)

			},
			expectedError:           nil,
			checkError:              nil,
			expectedProceduresCount: 2,
			expectedTotalCount:      2,
		},
		{
			name: "empty result",
			queryParams: QueryParams{
				Limit:  10,
				Offset: 0,
				Search: "NonExistent",
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures" WHERE LOWER\(procedures.title\) LIKE \$1`).
					WithArgs("%nonexistent%").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE LOWER\(procedures.title\) LIKE \$1 LIMIT \$2`).
					WithArgs("%nonexistent%", 10).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
					}))
			},
			expectedError:           nil,
			expectedProceduresCount: 0,
			expectedTotalCount:      0,
		},
		{
			name: "pagination",
			queryParams: QueryParams{
				Limit:  1,
				Offset: 2,
				Search: "",
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				procedures := TestProcedureList()[1:]

				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				})

				for _, p := range procedures {
					rows.AddRow(
						p.ID,
						p.Title,
						p.Description,
						p.Version,
						p.ParentID,
						p.CreatedAt,
						p.UpdatedAt,
					)
				}

				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures"`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				mock.ExpectQuery(`SELECT \* FROM "procedures" LIMIT \$1 OFFSET \$2`).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			expectedError:           nil,
			expectedProceduresCount: 1,
			expectedTotalCount:      2,
		},
		{
			name: "count error",
			queryParams: QueryParams{
				Limit:  10,
				Offset: 0,
				Search: "",
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures"`).
					WillReturnError(assert.AnError)
			},
			expectedError: assert.AnError,
		},
		{
			name: "find error",
			queryParams: QueryParams{
				Limit:  10,
				Offset: 0,
				Search: "",
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures"`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				mock.ExpectQuery(`SELECT \* FROM "procedures" LIMIT \$1`).
					WithArgs(10).
					WillReturnError(assert.AnError)
			},
			expectedError: assert.AnError,
		},
		{
			name: "context canceled",
			queryParams: QueryParams{
				Limit:  10,
				Offset: 0,
				Search: "",
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setupMock: func(mock sqlmock.Sqlmock) {
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, context.Canceled)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			tc.setupMock(mock)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			procedures, total, err := repo.FindAll(ctx, tc.queryParams.Search, tc.queryParams.Offset, tc.queryParams.Limit)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, procedures, tc.expectedProceduresCount)
				assert.Equal(t, tc.expectedTotalCount, total)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_FindByID(t *testing.T) {
	testID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	proc := TestProcedureDetail()

	tests := []struct {
		name            string
		id              uuid.UUID
		withExperiments bool
		ctx             func() context.Context
		setupMock       func(mock sqlmock.Sqlmock, id uuid.UUID)
		expectedError   error
		checkError      func(t *testing.T, err error)
		checkResult     func(t *testing.T, result *procedure.Procedure)
	}{
		{
			name:            "success_no_preloads",
			id:              testID,
			withExperiments: false,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				}).AddRow(
					proc.ID,
					proc.Title,
					proc.Description,
					proc.Version,
					proc.ParentID,
					proc.CreatedAt,
					proc.UpdatedAt,
				)

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(rows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.NotNil(t, result)
				assert.Equal(t, testID, result.ID)
				assert.Equal(t, "Test Procedure", result.Title)
			},
		},
		{
			name:            "procedure_not_found",
			id:              uuid.New(),
			withExperiments: false,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
					}))
			},
			expectedError: constant.ErrProcedureNotFound,
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.Nil(t, result)
			},
		},
		{
			name:            "db_error_unexpected",
			id:              testID,
			withExperiments: false,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnError(assert.AnError)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.NotErrorIs(t, err, constant.ErrProcedureNotFound)
			},
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.Nil(t, result)
			},
		},
		{
			name:            "context_canceled",
			id:              testID,
			withExperiments: false,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				// No expectations because context is canceled
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.Nil(t, result)
			},
		},
		{
			name:            "correct_table_queried",
			id:              testID,
			withExperiments: false,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				}).AddRow(
					id, "Expected Title", "Expected Description", 1, nil, time.Now(), time.Now(),
				)

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(rows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.NotNil(t, result)
				assert.Equal(t, "Expected Title", result.Title)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			tc.setupMock(mock, tc.id)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			result, err := repo.FindByID(ctx, tc.id, tc.withExperiments)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.checkResult != nil {
				tc.checkResult(t, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateProcedure(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, proc *procedure.Procedure)
		procedureFunc func() *procedure.Procedure
		ctx           func() context.Context
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "success_minimal",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "procedures"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(proc.ID, time.Now()))
				mock.ExpectCommit()
			},
			procedureFunc: func() *procedure.Procedure {
				return CreateValidProcedure()
			},
			ctx:           func() context.Context { return context.Background() },
			expectedError: nil,
		},
		{
			name: "success_full",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				mock.ExpectBegin()

				// Insert procedure
				mock.ExpectQuery(`INSERT INTO "procedures"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(proc.ID, time.Now()))

				// Insert steps
				for _, step := range proc.Steps {
					mock.ExpectQuery(`INSERT INTO "procedure_steps"`).
						WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
							AddRow(step.ID, time.Now()))
				}

				mock.ExpectCommit()
			},
			procedureFunc: func() *procedure.Procedure {
				procID := uuid.New()

				return &procedure.Procedure{
					ID:          procID,
					Title:       "Full Test Procedure",
					Description: stringutil.StringPtr("Full Test Description"),
					Version:     1,
					ParentID:    nil,
				}
			},
			ctx:           func() context.Context { return context.Background() },
			expectedError: nil,
		},
		{
			name: "duplicate_title",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "procedures"`).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			procedureFunc: func() *procedure.Procedure {
				return CreateValidProcedure()
			},
			ctx: func() context.Context { return context.Background() },
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "success_with_steps",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "procedures"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(proc.ID, time.Now()))

				// GORM batches the steps insert into a single query
				if len(proc.Steps) > 0 {
					rows := sqlmock.NewRows([]string{"id", "created_at"})
					for _, step := range proc.Steps {
						rows.AddRow(step.ID, time.Now())
					}
					mock.ExpectQuery(`INSERT INTO "procedure_steps"`).
						WillReturnRows(rows)
				}

				mock.ExpectCommit()
			},
			procedureFunc: func() *procedure.Procedure {
				return CreateProcedureWithSteps()
			},
			ctx:           func() context.Context { return context.Background() },
			expectedError: nil,
		},
		{
			name: "success_with_parent",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "procedures"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(proc.ID, time.Now()))
				mock.ExpectCommit()
			},
			procedureFunc: func() *procedure.Procedure {
				parentID := uuid.New()
				return CreateProcedureWithParent(parentID)
			},
			ctx:           func() context.Context { return context.Background() },
			expectedError: nil,
		},
		{
			name: "invalid_parent_id",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "procedures"`).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			procedureFunc: func() *procedure.Procedure {
				invalidParentID := uuid.New()
				return CreateProcedureWithParent(invalidParentID)
			},
			ctx: func() context.Context { return context.Background() },
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "context_cancelled",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				// No DB expectations since context is cancelled
			},
			procedureFunc: func() *procedure.Procedure {
				return CreateValidProcedure()
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "db_error",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "procedures"`).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			procedureFunc: func() *procedure.Procedure {
				return CreateValidProcedure()
			},
			ctx: func() context.Context { return context.Background() },
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			proc := tc.procedureFunc()
			tc.setupMock(mock, proc)

			err := repo.CreateProcedure(ctx, proc)

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

func TestRepository_UpdateProcedure(t *testing.T) {
	testID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := []struct {
		name          string
		id            uuid.UUID
		procedureFunc func() *procedure.Procedure
		setupMock     func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "success",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: stringutil.StringPtr("Updated Description"),
					Version:     1,
					UpdatedAt:   timeutil.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "procedures" SET`).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "db_error",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: stringutil.StringPtr("Updated Description"),
					Version:     1,
					UpdatedAt:   timeutil.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "procedures" SET`).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			name: "not_found",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: stringutil.StringPtr("Updated Description"),
					Version:     1,
					UpdatedAt:   timeutil.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "procedures" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures" WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedError: constant.ErrProcedureNotFound,
		},
		{
			name: "version_conflict",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: stringutil.StringPtr("Updated Description"),
					Version:     1, // Version mismatch
					UpdatedAt:   timeutil.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "procedures" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()

				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures" WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedError: constant.ErrProcedureConflict,
		},
		{
			name: "count_error",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: stringutil.StringPtr("Updated Description"),
					Version:     1,
					UpdatedAt:   timeutil.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				mock.ExpectBegin()
				// UPDATE succeeds but affects 0 rows
				mock.ExpectExec(`UPDATE "procedures" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()

				// COUNT query returns error
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures" WHERE id = \$1`).
					WithArgs(id).
					WillReturnError(assert.AnError)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				// The error should be the count query error, not wrapped
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := procedure.NewRepository(db)
			ctx := context.Background()

			proc := tc.procedureFunc()
			tc.setupMock(mock, tc.id, proc)

			err := repo.UpdateProcedure(ctx, tc.id, proc)

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

func TestRepository_CreateProcedureStep(t *testing.T) {
	tests := []struct {
		name          string
		stepFunc      func() *procedure.ProcedureStep
		ctx           func() context.Context
		setupMock     func(mock sqlmock.Sqlmock, step *procedure.ProcedureStep)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name:     "success",
			stepFunc: TestProcedureStep,
			ctx:      func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, step *procedure.ProcedureStep) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "procedure_steps"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(step.ID, time.Now()))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:     "db_error",
			stepFunc: TestProcedureStep,
			ctx:      func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, step *procedure.ProcedureStep) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "procedure_steps"`).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			name:     "context_cancelled",
			stepFunc: TestProcedureStep,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setupMock: func(mock sqlmock.Sqlmock, step *procedure.ProcedureStep) {
				// No DB expectations since context is cancelled
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, context.Canceled)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			step := tc.stepFunc()
			tc.setupMock(mock, step)

			err := repo.CreateProcedureStep(ctx, step)

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

func TestRepository_GetStepIDsByProcID(t *testing.T) {
	testProcID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := []struct {
		name          string
		procedureID   uuid.UUID
		ctx           func() context.Context
		setupMock     func(mock sqlmock.Sqlmock, procedureID uuid.UUID)
		expectedError error
		checkError    func(t *testing.T, err error)
		checkResult   func(t *testing.T, result []procedure.StepMetadata)
	}{
		{
			name:        "success_multiple",
			procedureID: testProcID,
			ctx:         func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, procedureID uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "version"}).
					AddRow(uuid.MustParse("33333333-1234-1234-1234-444433332222"), 1).
					AddRow(uuid.MustParse("33333333-1234-1234-1234-555533332222"), 2)
				mock.ExpectQuery(`SELECT "id","version" FROM "procedure_steps" WHERE procedure_id = \$1`).
					WithArgs(procedureID).
					WillReturnRows(rows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []procedure.StepMetadata) {
				assert.Len(t, result, 2)
			},
		},
		{
			name:        "success_empty",
			procedureID: testProcID,
			ctx:         func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, procedureID uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "version"})
				mock.ExpectQuery(`SELECT "id","version" FROM "procedure_steps" WHERE procedure_id = \$1`).
					WithArgs(procedureID).
					WillReturnRows(rows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []procedure.StepMetadata) {
				assert.Len(t, result, 0)
			},
		},
		{
			name:        "db_error",
			procedureID: testProcID,
			ctx:         func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, procedureID uuid.UUID) {
				mock.ExpectQuery(`SELECT "id","version" FROM "procedure_steps" WHERE procedure_id = \$1`).
					WithArgs(procedureID).
					WillReturnError(assert.AnError)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			name:        "context_cancelled",
			procedureID: testProcID,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setupMock: func(mock sqlmock.Sqlmock, procedureID uuid.UUID) {
				// No DB expectations since context is cancelled
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, context.Canceled)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			tc.setupMock(mock, tc.procedureID)

			result, err := repo.GetStepIDsByProcID(ctx, tc.procedureID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.checkResult != nil {
				tc.checkResult(t, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_UpdateProcedureStep(t *testing.T) {
	testStepID := uuid.MustParse("33333333-1234-1234-1234-444433332222")
	testProcID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := []struct {
		name          string
		stepID        uuid.UUID
		procedureID   uuid.UUID
		stepFunc      func() *procedure.ProcedureStep
		ctx           func() context.Context
		setupMock     func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID, step *procedure.ProcedureStep)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name:        "success",
			stepID:      testStepID,
			procedureID: testProcID,
			stepFunc: func() *procedure.ProcedureStep {
				s := TestProcedureStep()
				s.Title = "Updated Step Title"
				s.Description = stringutil.StringPtr("Updated Step Description")
				s.UpdatedAt = timeutil.TimePtr(time.Now())
				return s
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID, step *procedure.ProcedureStep) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "procedure_steps" SET`).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:        "not_found",
			stepID:      testStepID,
			procedureID: testProcID,
			stepFunc: func() *procedure.ProcedureStep {
				s := TestProcedureStep()
				s.UpdatedAt = timeutil.TimePtr(time.Now())
				return s
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID, step *procedure.ProcedureStep) {
				mock.ExpectBegin()
				// UPDATE affects 0 rows
				mock.ExpectExec(`UPDATE "procedure_steps" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
				// COUNT returns 0 — step not found
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedure_steps" WHERE id = \$1 AND procedure_id = \$2`).
					WithArgs(stepID, procedureID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedError: constant.ErrProcedureStepNotFound,
		},
		{
			name:        "optimistic_lock_conflict",
			stepID:      testStepID,
			procedureID: testProcID,
			stepFunc: func() *procedure.ProcedureStep {
				s := TestProcedureStep()
				s.UpdatedAt = timeutil.TimePtr(time.Now())
				return s
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID, step *procedure.ProcedureStep) {
				mock.ExpectBegin()
				// UPDATE affects 0 rows (version mismatch)
				mock.ExpectExec(`UPDATE "procedure_steps" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
				// COUNT returns 1 — step exists but version mismatch
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedure_steps" WHERE id = \$1 AND procedure_id = \$2`).
					WithArgs(stepID, procedureID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedError: constant.ErrProcedureStepConflict,
		},
		{
			name:        "db_error_on_update",
			stepID:      testStepID,
			procedureID: testProcID,
			stepFunc: func() *procedure.ProcedureStep {
				s := TestProcedureStep()
				s.UpdatedAt = timeutil.TimePtr(time.Now())
				return s
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID, step *procedure.ProcedureStep) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "procedure_steps" SET`).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			name:        "db_error_on_count",
			stepID:      testStepID,
			procedureID: testProcID,
			stepFunc: func() *procedure.ProcedureStep {
				s := TestProcedureStep()
				s.UpdatedAt = timeutil.TimePtr(time.Now())
				return s
			},
			ctx: func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID, step *procedure.ProcedureStep) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "procedure_steps" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
				// COUNT query returns error
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedure_steps" WHERE id = \$1 AND procedure_id = \$2`).
					WithArgs(stepID, procedureID).
					WillReturnError(assert.AnError)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			name:        "context_cancelled",
			stepID:      testStepID,
			procedureID: testProcID,
			stepFunc: func() *procedure.ProcedureStep {
				s := TestProcedureStep()
				s.UpdatedAt = timeutil.TimePtr(time.Now())
				return s
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID, step *procedure.ProcedureStep) {
				// No DB expectations since context is cancelled
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			step := tc.stepFunc()
			tc.setupMock(mock, tc.stepID, tc.procedureID, step)

			err := repo.UpdateProcedureStep(ctx, tc.stepID, tc.procedureID, step)

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

func TestRepository_DeleteProcedureStep(t *testing.T) {
	testStepID := uuid.MustParse("33333333-1234-1234-1234-444433332222")
	testProcID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := []struct {
		name          string
		stepID        uuid.UUID
		procedureID   uuid.UUID
		ctx           func() context.Context
		setupMock     func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name:        "success",
			stepID:      testStepID,
			procedureID: testProcID,
			ctx:         func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "procedure_steps" WHERE id = \$1 AND procedure_id = \$2`).
					WithArgs(stepID, procedureID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:        "not_found",
			stepID:      testStepID,
			procedureID: testProcID,
			ctx:         func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID) {
				mock.ExpectBegin()
				// Hard delete affects 0 rows
				mock.ExpectExec(`DELETE FROM "procedure_steps" WHERE id = \$1 AND procedure_id = \$2`).
					WithArgs(stepID, procedureID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectedError: constant.ErrProcedureStepNotFound,
		},
		{
			name:        "db_error",
			stepID:      testStepID,
			procedureID: testProcID,
			ctx:         func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "procedure_steps" WHERE id = \$1 AND procedure_id = \$2`).
					WithArgs(stepID, procedureID).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			name:        "context_cancelled",
			stepID:      testStepID,
			procedureID: testProcID,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setupMock: func(mock sqlmock.Sqlmock, stepID uuid.UUID, procedureID uuid.UUID) {
				// No DB expectations since context is cancelled
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			tc.setupMock(mock, tc.stepID, tc.procedureID)

			err := repo.DeleteProcedureStep(ctx, tc.stepID, tc.procedureID)

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

func TestRepository_DeleteProcedure(t *testing.T) {
	testID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := []struct {
		name          string
		id            uuid.UUID
		ctx           func() context.Context
		setupMock     func(mock sqlmock.Sqlmock, id uuid.UUID)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "success",
			id:   testID,
			ctx:  func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "procedures" WHERE id = \$1`).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "not_found",
			id:   testID,
			ctx:  func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "procedures" WHERE id = \$1`).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectedError: constant.ErrProcedureNotFound,
		},
		{
			name: "db_error",
			id:   testID,
			ctx:  func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "procedures" WHERE id = \$1`).
					WithArgs(id).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		{
			name: "context_cancelled",
			id:   testID,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {

			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, context.Canceled)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			tc.setupMock(mock, tc.id)

			err := repo.DeleteProcedure(ctx, tc.id)

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
