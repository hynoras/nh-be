package experiment

import (
	"context"
	"testing"

	"nh-be/internal/constant"
	"nh-be/internal/features/experiment"
	"nh-be/internal/utils/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetProcedureIDByID(t *testing.T) {
	exp := TestExperiment()
	procedureID := *exp.ProcedureID

	tests := []struct {
		name          string
		id            uuid.UUID
		ctx           func() context.Context
		setupMock     func(mock sqlmock.Sqlmock, id uuid.UUID)
		expectedError error
		checkError    func(t *testing.T, err error)
		checkResult   func(t *testing.T, result uuid.UUID)
	}{
		{
			name: "success",
			id:   exp.ID,
			ctx:  func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{"procedure_id"}).
					AddRow(procedureID.String())

				mock.ExpectQuery(`SELECT "procedure_id" FROM "experiments" WHERE id = \$1 ORDER BY "experiments"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(rows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result uuid.UUID) {
				assert.Equal(t, procedureID, result)
			},
		},
		{
			name: "success_null",
			id:   exp.ID,
			ctx:  func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{"procedure_id"}).
					AddRow(nil)

				mock.ExpectQuery(`SELECT "procedure_id" FROM "experiments" WHERE id = \$1 ORDER BY "experiments"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(rows)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result uuid.UUID) {
				assert.Equal(t, uuid.Nil, result)
			},
		},
		{
			name: "not_found",
			id:   uuid.New(),
			ctx:  func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT "procedure_id" FROM "experiments" WHERE id = \$1 ORDER BY "experiments"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"procedure_id"}))
			},
			expectedError: experiment.ErrExperimentNotFound,
			checkResult: func(t *testing.T, result uuid.UUID) {
				assert.Equal(t, uuid.Nil, result)
			},
		},
		{
			name: "db_error",
			id:   exp.ID,
			ctx:  func() context.Context { return context.Background() },
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT "procedure_id" FROM "experiments" WHERE id = \$1 ORDER BY "experiments"\."id" LIMIT \$2`).
					WithArgs(id, 1).
					WillReturnError(assert.AnError)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.NotErrorIs(t, err, experiment.ErrExperimentNotFound)
			},
			checkResult: func(t *testing.T, result uuid.UUID) {
				assert.Equal(t, uuid.Nil, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			tc.setupMock(mock, tc.id)
			repo := experiment.NewRepository(db)
			ctx := tc.ctx()

			result, err := repo.GetProcedureIDByID(ctx, tc.id)

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

func TestRepository_UpdateProcedureID(t *testing.T) {
	testID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	procedureID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	currentVersion := 1

	tests := []struct {
		name          string
		id            uuid.UUID
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "success",
			id:   testID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiments" SET`).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "conflict",
			id:   testID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiments" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectedError: experiment.ErrExperimentConflict,
		},
		{
			name: "not_found",
			id:   uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiments" SET`).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectedError: experiment.ErrExperimentConflict,
		},
		{
			name: "db_error",
			id:   testID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "experiments" SET`).
					WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			tc.setupMock(mock)
			repo := experiment.NewRepository(db)
			ctx := context.Background()

			err := repo.UpdateProcedureID(ctx, tc.id, procedureID, currentVersion)

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

func TestRepository_FindAllExperiments(t *testing.T) {
	sample := TestExperimentsQueryDto()

	// columns as returned by the Scan into ExperimentsQueryDto
	cols := []string{
		"identifier", "title", "objective", "status", "type",
		"creator", "updater", "created_at", "updated_at", "procedure_name",
	}

	// reusable helper to build a rows result from the sample data
	makeRows := func() *sqlmock.Rows {
		rows := sqlmock.NewRows(cols)
		for _, e := range sample {
			rows.AddRow(
				e.Identifier, e.Title, e.Objective,
				string(e.Status), string(e.Type),
				e.Creator, e.Updater,
				e.CreatedAt, e.UpdatedAt,
				e.ProcedureName,
			)
		}
		return rows
	}

	// broad regex that matches the SELECT … FROM "experiments" … LIMIT … OFFSET pattern
	queryRegex := `SELECT .+ FROM "experiments" .+`

	status := experiment.ExperimentDraft
	expType := experiment.ExperimentExploratoryType
	creatorID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	tests := []struct {
		name             string
		sortBy           string
		search           string
		sortOrder        constant.Order
		experimentStatus *experiment.ExperimentStatus
		experimentType   *experiment.ExperimentType
		createdBy        *uuid.UUID
		page             int
		pageSize         int
		setupMock        func(mock sqlmock.Sqlmock)
		expectedError    error
		checkError       func(t *testing.T, err error)
		checkResult      func(t *testing.T, result []experiment.ExperimentsQueryDto)
	}{
		{
			name:             "success",
			sortBy:           "created_at",
			search:           "",
			sortOrder:        constant.DESC,
			experimentStatus: nil,
			experimentType:   nil,
			createdBy:        nil,
			page:             1,
			pageSize:         10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// args: LIKE pattern, LIMIT (OFFSET omitted by GORM when 0)
				mock.ExpectQuery(queryRegex).
					WithArgs("%%", 10).
					WillReturnRows(makeRows())
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []experiment.ExperimentsQueryDto) {
				assert.Equal(t, sample, result)
			},
		},
		{
			name:             "success_status_filter",
			sortBy:           "created_at",
			search:           "",
			sortOrder:        constant.ASC,
			experimentStatus: &status,
			experimentType:   nil,
			createdBy:        nil,
			page:             1,
			pageSize:         10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// args: LIKE pattern, status, LIMIT (OFFSET omitted when 0)
				mock.ExpectQuery(queryRegex).
					WithArgs("%%", string(status), 10).
					WillReturnRows(makeRows())
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []experiment.ExperimentsQueryDto) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:             "success_type_filter",
			sortBy:           "created_at",
			search:           "",
			sortOrder:        constant.ASC,
			experimentStatus: nil,
			experimentType:   &expType,
			createdBy:        nil,
			page:             1,
			pageSize:         10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// args: LIKE pattern, type, LIMIT (OFFSET omitted when 0)
				mock.ExpectQuery(queryRegex).
					WithArgs("%%", string(expType), 10).
					WillReturnRows(makeRows())
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []experiment.ExperimentsQueryDto) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:             "success_status_and_type_filter",
			sortBy:           "created_at",
			search:           "",
			sortOrder:        constant.ASC,
			experimentStatus: &status,
			experimentType:   &expType,
			createdBy:        nil,
			page:             1,
			pageSize:         10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// args: LIKE pattern, status, type, LIMIT (OFFSET omitted when 0)
				mock.ExpectQuery(queryRegex).
					WithArgs("%%", string(status), string(expType), 10).
					WillReturnRows(makeRows())
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []experiment.ExperimentsQueryDto) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:             "success_search_keyword",
			sortBy:           "created_at",
			search:           "test",
			sortOrder:        constant.ASC,
			experimentStatus: nil,
			experimentType:   nil,
			createdBy:        nil,
			page:             1,
			pageSize:         10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// args: LIKE pattern with keyword, LIMIT (OFFSET omitted when 0)
				mock.ExpectQuery(queryRegex).
					WithArgs("%test%", 10).
					WillReturnRows(makeRows())
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []experiment.ExperimentsQueryDto) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:             "success_created_by_filter",
			sortBy:           "created_at",
			search:           "",
			sortOrder:        constant.ASC,
			experimentStatus: nil,
			experimentType:   nil,
			createdBy:        &creatorID,
			page:             1,
			pageSize:         10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// args: LIKE pattern, created_by, LIMIT (OFFSET omitted when 0)
				mock.ExpectQuery(queryRegex).
					WithArgs("%%", creatorID.String(), 10).
					WillReturnRows(makeRows())
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []experiment.ExperimentsQueryDto) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:             "success_pagination",
			sortBy:           "created_at",
			search:           "",
			sortOrder:        constant.ASC,
			experimentStatus: nil,
			experimentType:   nil,
			createdBy:        nil,
			page:             3,
			pageSize:         5,
			setupMock: func(mock sqlmock.Sqlmock) {
				// page=3, pageSize=5 → offset=(3-1)*5=10, limit=5
				mock.ExpectQuery(queryRegex).
					WithArgs("%%", 5, 10).
					WillReturnRows(makeRows())
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result []experiment.ExperimentsQueryDto) {
				assert.NotEmpty(t, result)
			},
		},
		{
			name:             "database_error",
			sortBy:           "created_at",
			search:           "",
			sortOrder:        constant.ASC,
			experimentStatus: nil,
			experimentType:   nil,
			createdBy:        nil,
			page:             1,
			pageSize:         10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// args: LIKE pattern, LIMIT (OFFSET omitted when 0)
				mock.ExpectQuery(queryRegex).
					WithArgs("%%", 10).
					WillReturnError(assert.AnError)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			tc.setupMock(mock)
			repo := experiment.NewRepository(db)
			ctx := context.Background()

			result, err := repo.FindAllExperiments(
				ctx,
				tc.sortBy,
				tc.search,
				tc.sortOrder,
				tc.experimentStatus,
				tc.experimentType,
				tc.createdBy,
				tc.page,
				tc.pageSize,
			)

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

func TestRepository_FindByIdentifierAndCreatedBy(t *testing.T) {
	exp := TestExperiment()

	tests := []struct {
		name          string
		identifier    string
		createdBy     uuid.UUID
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
		checkError    func(t *testing.T, err error)
		checkResult   func(t *testing.T, result *experiment.Experiment)
	}{
		{
			name:       "success",
			identifier: exp.Identifier,
			createdBy:  exp.CreatedByID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "identifier", "title", "objective", "status", "type", "version",
					"created_by_id", "updated_by_id", "procedure_id", "created_at", "updated_at",
				}).AddRow(
					exp.ID, exp.Identifier, exp.Title, exp.Objective, exp.Status, exp.Type, exp.Version,
					exp.CreatedByID, exp.UpdatedByID, exp.ProcedureID, exp.CreatedAt, exp.UpdatedAt,
				)

				mock.ExpectQuery(`^SELECT .+ FROM "experiments" WHERE identifier = \$1 AND created_by_id = \$2 ORDER BY "experiments"\."id" LIMIT \$3$`).
					WithArgs(exp.Identifier, exp.CreatedByID, 1).
					WillReturnRows(rows)

				mock.ExpectQuery(`^SELECT .+ FROM "users" WHERE "users"\."id" = \$1$`).
					WithArgs(exp.CreatedByID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(exp.CreatedByID))
			},
			expectedError: nil,
			checkResult: func(t *testing.T, result *experiment.Experiment) {
				assert.NotNil(t, result)
				assert.Equal(t, exp.ID, result.ID)
				assert.Equal(t, exp.Identifier, result.Identifier)
			},
		},
		{
			name:       "not_found",
			identifier: "non-existent",
			createdBy:  uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`^SELECT .+ FROM "experiments" WHERE identifier = \$1 AND created_by_id = \$2 ORDER BY "experiments"\."id" LIMIT \$3$`).
					WithArgs("non-existent", sqlmock.AnyArg(), 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"})) // empty rows
			},
			expectedError: experiment.ErrExperimentNotFound,
			checkResult: func(t *testing.T, result *experiment.Experiment) {
				assert.Nil(t, result)
			},
		},
		{
			name:       "database_error",
			identifier: exp.Identifier,
			createdBy:  exp.CreatedByID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`^SELECT .+ FROM "experiments" WHERE identifier = \$1 AND created_by_id = \$2 ORDER BY "experiments"\."id" LIMIT \$3$`).
					WithArgs(exp.Identifier, exp.CreatedByID, 1).
					WillReturnError(assert.AnError)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
			checkResult: func(t *testing.T, result *experiment.Experiment) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := testutil.SetupMockDB(t)
			tc.setupMock(mock)
			repo := experiment.NewRepository(db)
			ctx := context.Background()

			result, err := repo.FindByIdentifierAndCreatedBy(ctx, tc.identifier, tc.createdBy)

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
