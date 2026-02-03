package procedure

import (
	"context"
	"errors"
	"testing"

	"nh-be/constant"
	"nh-be/internal/permission/mocks"
	"nh-be/internal/procedure"
	procmocks "nh-be/internal/procedure/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetAllProcedures(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		ctx            context.Context
		search         string
		offset         int
		limit          int
		setupMocks     func(repo *procmocks.Repository, permSvc *mocks.Service)
		expectedResult []procedure.ProcedureListResponseDto
		expectedLength int64
		expectedError  error
		checkError     func(t *testing.T, err error)
	}{
		{
			name:   "ok_view_permission",
			ctx:    ContextWithUser(userID),
			search: "",
			offset: 0,
			limit:  10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindAll", mock.Anything, "", 0, 10, true).
					Return(TestProcedureList(), int64(2), nil)
			},
			expectedResult: procedure.MapProceduresToDto(TestProcedureList()),
			expectedLength: 2,
			expectedError:  nil,
		},
		{
			name:   "ok_manage_permission",
			ctx:    ContextWithUser(userID),
			search: "",
			offset: 0,
			limit:  10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("FindAll", mock.Anything, "", 0, 10, true).
					Return(TestProcedureList(), int64(2), nil)
			},
			expectedResult: procedure.MapProceduresToDto(TestProcedureList()),
			expectedLength: 2,
			expectedError:  nil,
		},
		{
			name:   "forbidden_no_permission",
			ctx:    ContextWithUser(userID),
			search: "",
			offset: 0,
			limit:  10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
				// Repo should NOT be called
			},
			expectedResult: nil,
			expectedLength: 0,
			expectedError:  constant.ErrForbidViewProcedure,
		},
		{
			name:   "permission_service_error",
			ctx:    ContextWithUser(userID),
			search: "",
			offset: 0,
			limit:  10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, errors.New("permission service unavailable"))
				// Repo should NOT be called
			},
			expectedResult: nil,
			expectedLength: 0,
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "permission service unavailable")
			},
		},
		{
			name:   "no_user_in_context",
			ctx:    context.Background(), // No user in context
			search: "",
			offset: 0,
			limit:  10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {

			},
			expectedResult: nil,
			expectedLength: 0,
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name:   "empty_result",
			ctx:    ContextWithUser(userID),
			search: "",
			offset: 0,
			limit:  10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindAll", mock.Anything, "", 0, 10, true).
					Return([]procedure.Procedure{}, int64(0), nil)
			},
			expectedResult: []procedure.ProcedureListResponseDto{},
			expectedLength: 0,
			expectedError:  nil,
		},
		{
			name:   "repo_error",
			ctx:    ContextWithUser(userID),
			search: "test",
			offset: 0,
			limit:  10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindAll", mock.Anything, "test", 0, 10, true).
					Return(nil, int64(0), errors.New("database connection error"))
			},
			expectedResult: nil,
			expectedLength: 0,
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database connection error")
			},
		},
		{
			name:   "pagination_passthrough",
			ctx:    ContextWithUser(userID),
			search: "",
			offset: 20,
			limit:  50,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindAll", mock.Anything, "", 20, 50, true).
					Return([]procedure.Procedure{}, int64(0), nil)
			},
			expectedResult: []procedure.ProcedureListResponseDto{},
			expectedLength: 0,
			expectedError:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := procmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := procedure.NewService(mockRepo, mockPermSvc)

			res, length, err := svc.GetAllProcedures(tc.ctx, tc.search, tc.offset, tc.limit)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, res)
				assert.Equal(t, tc.expectedLength, length)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
				assert.Nil(t, res)
				assert.Equal(t, tc.expectedLength, length)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, len(tc.expectedResult), len(res))
				assert.Equal(t, tc.expectedLength, length)
			}

		})
	}
}
