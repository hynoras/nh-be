package procedure

import (
	"context"
	"nh-be/internal/procedure"
	"nh-be/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
