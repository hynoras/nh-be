package experiment

import (
	"context"
	"testing"

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
