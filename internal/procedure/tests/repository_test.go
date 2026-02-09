package procedure

import (
	"context"
	"nh-be/constant"
	"nh-be/internal/procedure"
	"nh-be/utils"
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
			db, mock := utils.SetupMockDB(t)
			tc.setupMock(mock)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			procedures, total, err := repo.FindAll(ctx, tc.queryParams.Search, tc.queryParams.Offset, tc.queryParams.Limit, false)

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
		withSteps       bool
		withExperiments bool
		ctx             func() context.Context
		setupMock       func(mock sqlmock.Sqlmock, id uuid.UUID)
		expectedError   error
		checkError      func(t *testing.T, err error)
		checkResult     func(t *testing.T, result *procedure.Procedure)
	}{
		{
			name:            "happy_path_no_preloads",
			id:              testID,
			withSteps:       false,
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
			withSteps:       false,
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
			withSteps:       false,
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
			withSteps:       false,
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
			withSteps:       false,
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
		{
			name:            "load_steps_only",
			id:              testID,
			withSteps:       true,
			withExperiments: false,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				procedureRows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				}).AddRow(
					proc.ID, proc.Title, proc.Description, proc.Version, proc.ParentID, proc.CreatedAt, proc.UpdatedAt,
				)

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(procedureRows)

				stepRows := sqlmock.NewRows([]string{
					"id", "procedure_id", "index", "title", "description", "step_type", "created_at", "updated_at",
				}).AddRow(
					uuid.New(), id, 1, "Step 1", "Description 1", "manual", time.Now(), time.Now(),
				)

				mock.ExpectQuery(`SELECT \* FROM "procedure_steps" WHERE "procedure_steps"\."procedure_id" = \$1`).
					WithArgs(id).
					WillReturnRows(stepRows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.NotNil(t, result)
				assert.Greater(t, len(result.Steps), 0)
				assert.Equal(t, 0, len(result.Experiments))
			},
		},
		{
			name:            "load_experiments_only",
			id:              testID,
			withSteps:       false,
			withExperiments: true,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				procedureRows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				}).AddRow(
					proc.ID, proc.Title, proc.Description, proc.Version, proc.ParentID, proc.CreatedAt, proc.UpdatedAt,
				)

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(procedureRows)

				experimentRows := sqlmock.NewRows([]string{
					"id", "procedure_id", "experiment_id", "created_at", "updated_at",
				}).AddRow(
					uuid.New(), id, uuid.New(), time.Now(), time.Now(),
				)

				mock.ExpectQuery(`SELECT \* FROM "procedure_experiment_assignments" WHERE "procedure_experiment_assignments"\."procedure_id" = \$1`).
					WithArgs(id).
					WillReturnRows(experimentRows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.NotNil(t, result)
				assert.Equal(t, 0, len(result.Steps))
				assert.Greater(t, len(result.Experiments), 0)
			},
		},
		{
			name:            "load_both_relations",
			id:              testID,
			withSteps:       true,
			withExperiments: true,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				// Main procedure query
				procedureRows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				}).AddRow(
					proc.ID, proc.Title, proc.Description, proc.Version, proc.ParentID, proc.CreatedAt, proc.UpdatedAt,
				)

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(procedureRows)

				// Experiments preload query
				experimentRows := sqlmock.NewRows([]string{
					"id", "procedure_id", "experiment_id", "created_at", "updated_at",
				}).AddRow(
					uuid.New(), id, uuid.New(), time.Now(), time.Now(),
				)

				mock.ExpectQuery(`SELECT \* FROM "procedure_experiment_assignments" WHERE "procedure_experiment_assignments"\."procedure_id" = \$1`).
					WithArgs(id).
					WillReturnRows(experimentRows)

				// Steps preload query
				stepRows := sqlmock.NewRows([]string{
					"id", "procedure_id", "index", "title", "description", "step_type", "created_at", "updated_at",
				}).AddRow(
					uuid.New(), id, 1, "Step 1", "Description 1", "manual", time.Now(), time.Now(),
				)

				mock.ExpectQuery(`SELECT \* FROM "procedure_steps" WHERE "procedure_steps"\."procedure_id" = \$1`).
					WithArgs(id).
					WillReturnRows(stepRows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.NotNil(t, result)
				assert.Greater(t, len(result.Steps), 0)
				assert.Greater(t, len(result.Experiments), 0)
			},
		},
		{
			name:            "procedure_has_no_steps",
			id:              testID,
			withSteps:       true,
			withExperiments: false,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				procedureRows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				}).AddRow(
					proc.ID, proc.Title, proc.Description, proc.Version, proc.ParentID, proc.CreatedAt, proc.UpdatedAt,
				)

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(procedureRows)

				stepRows := sqlmock.NewRows([]string{
					"id", "procedure_id", "index", "title", "description", "step_type", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT \* FROM "procedure_steps" WHERE "procedure_steps"\."procedure_id" = \$1`).
					WithArgs(id).
					WillReturnRows(stepRows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Steps)
				assert.Equal(t, 0, len(result.Steps))
			},
		},
		{
			name:            "procedure_has_no_experiments",
			id:              testID,
			withSteps:       false,
			withExperiments: true,
			ctx:             func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				procedureRows := sqlmock.NewRows([]string{
					"id", "title", "description", "version", "parent_id", "created_at", "updated_at",
				}).AddRow(
					proc.ID, proc.Title, proc.Description, proc.Version, proc.ParentID, proc.CreatedAt, proc.UpdatedAt,
				)

				mock.ExpectQuery(`SELECT \* FROM "procedures" WHERE "procedures"\."id" = \$1 ORDER BY "procedures"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(procedureRows)

				// Empty experiments preload query
				experimentRows := sqlmock.NewRows([]string{
					"id", "procedure_id", "experiment_id", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT \* FROM "procedure_experiment_assignments" WHERE "procedure_experiment_assignments"\."procedure_id" = \$1`).
					WithArgs(id).
					WillReturnRows(experimentRows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result *procedure.Procedure) {
				assert.NotNil(t, result)
				assert.NotNil(t, result.Experiments)
				assert.Equal(t, 0, len(result.Experiments))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := utils.SetupMockDB(t)
			tc.setupMock(mock, tc.id)
			repo := procedure.NewRepository(db)
			ctx := tc.ctx()

			result, err := repo.FindByID(ctx, tc.id, tc.withSteps, tc.withExperiments)

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

				// Insert experiment assignments
				for _, exp := range proc.Experiments {
					mock.ExpectQuery(`INSERT INTO "procedure_experiment_assignments"`).
						WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
							AddRow(exp.ID, time.Now()))
				}

				mock.ExpectCommit()
			},
			procedureFunc: func() *procedure.Procedure {
				procID := uuid.New()
				exp1ID := uuid.New()

				return &procedure.Procedure{
					ID:          procID,
					Title:       "Full Test Procedure",
					Description: "Full Test Description",
					Version:     1,
					ParentID:    nil,
					Steps: []procedure.ProcedureStep{
						{
							ID:          uuid.New(),
							ProcedureID: procID,
							Index:       1,
							Title:       "Step 1",
							Description: "Desc 1",
							StepType:    "manual",
						},
					},
					Experiments: []procedure.ProcedureExperimentAssignment{
						{
							ID:           uuid.New(),
							ProcedureID:  procID,
							ExperimentID: exp1ID,
						},
					},
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
			name: "duplicate_proc_exp_assignment",
			setupMock: func(mock sqlmock.Sqlmock, proc *procedure.Procedure) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "procedures"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(proc.ID, time.Now()))

				mock.ExpectQuery(`INSERT INTO "procedure_experiment_assignments"`).
					WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
			procedureFunc: func() *procedure.Procedure {
				return CreateProcedureWithExperiments()
			},
			ctx: func() context.Context { return context.Background() },
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
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
			db, mock := utils.SetupMockDB(t)
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
			name: "update_success",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: "Updated Description",
					Version:     1,
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				// GORM wraps Updates in a transaction
				mock.ExpectBegin()
				// Expect UPDATE query to succeed with 1 row affected
				// Don't use WithArgs for Updates() with map - GORM handles it internally
				mock.ExpectExec(`UPDATE "procedures" SET`).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "update_db_error",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: "Updated Description",
					Version:     1,
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				mock.ExpectBegin()
				// Expect UPDATE query to return error
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
			name: "update_not_found",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: "Updated Description",
					Version:     1,
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				mock.ExpectBegin()
				// UPDATE succeeds but affects 0 rows
				mock.ExpectExec(`UPDATE "procedures" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()

				// COUNT query returns 0 (procedure not found)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures" WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedError: constant.ErrProcedureNotFound,
		},
		{
			name: "update_version_conflict",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: "Updated Description",
					Version:     1, // Version mismatch
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
			},
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID, proc *procedure.Procedure) {
				mock.ExpectBegin()
				// UPDATE succeeds but affects 0 rows (version mismatch)
				mock.ExpectExec(`UPDATE "procedures" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()

				// COUNT query returns 1 (procedure exists but version mismatch)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "procedures" WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedError: constant.ErrOptimisticLockingConflict,
		},
		{
			name: "update_count_error",
			id:   testID,
			procedureFunc: func() *procedure.Procedure {
				return &procedure.Procedure{
					Title:       "Updated Title",
					Description: "Updated Description",
					Version:     1,
					UpdatedAt:   utils.TimePtr(time.Now()),
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
			db, mock := utils.SetupMockDB(t)
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
