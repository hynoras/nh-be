package result

import (
	"context"
	"errors"
	"testing"

	"nh-be/constant"
	"nh-be/internal/features/experiment/result"
	resultmocks "nh-be/internal/features/experiment/result/mocks"
	permissionmocks "nh-be/internal/features/permission/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetResultByExperimentID(t *testing.T) {
	experimentID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name           string
		ctx            context.Context
		setupMocks     func(repo *resultmocks.Repository, permSvc *permissionmocks.Service)
		expectedResult *result.ExperimentResult
		expectedError  error
		checkError     func(t *testing.T, err error)
	}{
		{
			name: "ok_view_permission",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindByExperimentID", mock.Anything, experimentID).
					Return(CreateTestResult(experimentID), nil)
			},
			expectedResult: CreateTestResult(experimentID),
			expectedError:  nil,
		},
		{
			name: "ok_manage_permission",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("FindByExperimentID", mock.Anything, experimentID).
					Return(CreateTestResult(experimentID), nil)
			},
			expectedResult: CreateTestResult(experimentID),
			expectedError:  nil,
		},
		{
			name: "forbidden_no_permission",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
				// Repo should NOT be called
			},
			expectedResult: nil,
			expectedError:  result.ErrForbidViewExperimentResult,
		},
		{
			name: "permission_service_error",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, errors.New("permission service unavailable"))
				// Repo should NOT be called
			},
			expectedResult: nil,
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "permission service unavailable")
			},
		},
		{
			name: "no_user_in_context",
			ctx:  context.Background(), // No user in context
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				// Neither permission service nor repo should be called
			},
			expectedResult: nil,
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name: "result_not_found",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindByExperimentID", mock.Anything, experimentID).
					Return(nil, result.ErrExperimentResultNotFound)
			},
			expectedResult: nil,
			expectedError:  result.ErrExperimentResultNotFound,
		},
		{
			name: "repo_error",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindByExperimentID", mock.Anything, experimentID).
					Return(nil, errors.New("database connection error"))
			},
			expectedResult: nil,
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database connection error")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh mocks for each test case
			mockRepo := resultmocks.NewRepository(t)
			mockPermSvc := permissionmocks.NewService(t)

			// Setup mock expectations
			tc.setupMocks(mockRepo, mockPermSvc)

			// Create service with mocks
			svc := result.NewService(mockRepo, mockPermSvc)

			// Execute
			res, err := svc.GetResultByExperimentID(tc.ctx, experimentID)

			// Assert
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, res)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tc.expectedResult.ExperimentID, res.ExperimentID)
			}

			// Mock expectations are automatically asserted on cleanup via NewRepository/NewService
		})
	}
}

func TestService_CreateResult(t *testing.T) {
	experimentID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		ctx           context.Context
		setupMocks    func(repo *resultmocks.Repository, permSvc *permissionmocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "ok_create",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*result.ExperimentResult")).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "forbidden_no_permission",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
				// Repo should NOT be called
			},
			expectedError: result.ErrForbidCreateExperimentResult,
		},
		{
			name: "permission_service_error",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, errors.New("permission service unavailable"))
				// Repo should NOT be called
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "permission service unavailable")
			},
		},
		{
			name: "no_user_in_context",
			ctx:  context.Background(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				// Neither permission service nor repo should be called
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name: "duplicate_result",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*result.ExperimentResult")).
					Return(result.ErrExperimentResultAlreadyExists)
			},
			expectedError: result.ErrExperimentResultAlreadyExists,
		},
		{
			name: "repo_error",
			ctx:  ContextWithUser(userID),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*result.ExperimentResult")).
					Return(errors.New("database connection error"))
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database connection error")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := resultmocks.NewRepository(t)
			mockPermSvc := permissionmocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := result.NewService(mockRepo, mockPermSvc)
			dto := CreateTestDto()

			err := svc.CreateResult(tc.ctx, experimentID, dto)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_UpdateResult(t *testing.T) {
	experimentID := uuid.New()
	resultID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		ctx           context.Context
		dto           *result.UpdateResultDto
		setupMocks    func(repo *resultmocks.Repository, permSvc *permissionmocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "ok_partial_update",
			ctx:  ContextWithUser(userID),
			dto:  CreateTestUpdateDto(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("FindByIDAndExperimentID", mock.Anything, resultID, experimentID).
					Return(CreateTestResult(experimentID), nil)
				repo.On("Update", mock.Anything, resultID, experimentID, mock.AnythingOfType("*result.UpdateFields"), 1).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "forbidden_no_permission",
			ctx:  ContextWithUser(userID),
			dto:  CreateTestUpdateDto(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
				// Repo should NOT be called
			},
			expectedError: result.ErrForbidUpdateExperimentResult,
		},
		{
			name: "permission_service_error",
			ctx:  ContextWithUser(userID),
			dto:  CreateTestUpdateDto(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, errors.New("permission service unavailable"))
				// Repo should NOT be called
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "permission service unavailable")
			},
		},
		{
			name: "no_user_in_context",
			ctx:  context.Background(),
			dto:  CreateTestUpdateDto(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				// Neither permission service nor repo should be called
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name: "result_not_found",
			ctx:  ContextWithUser(userID),
			dto:  CreateTestUpdateDto(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("FindByIDAndExperimentID", mock.Anything, resultID, experimentID).
					Return(nil, result.ErrExperimentResultNotFound)
				// Update should NOT be called
			},
			expectedError: result.ErrExperimentResultNotFound,
		},
		{
			name: "optimistic_lock_conflict",
			ctx:  ContextWithUser(userID),
			dto:  CreateTestUpdateDto(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("FindByIDAndExperimentID", mock.Anything, resultID, experimentID).
					Return(CreateTestResult(experimentID), nil)
				repo.On("Update", mock.Anything, resultID, experimentID, mock.AnythingOfType("*result.UpdateFields"), 1).
					Return(result.ErrOptimisticLockingConflict)
			},
			expectedError: result.ErrOptimisticLockingConflict,
		},
		{
			name: "repo_update_error",
			ctx:  ContextWithUser(userID),
			dto:  CreateTestUpdateDto(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("FindByIDAndExperimentID", mock.Anything, resultID, experimentID).
					Return(CreateTestResult(experimentID), nil)
				repo.On("Update", mock.Anything, resultID, experimentID, mock.AnythingOfType("*result.UpdateFields"), 1).
					Return(errors.New("database connection error"))
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database connection error")
			},
		},
		{
			name: "empty_update_fields",
			ctx:  ContextWithUser(userID),
			dto:  CreateEmptyUpdateDto(),
			setupMocks: func(repo *resultmocks.Repository, permSvc *permissionmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("FindByIDAndExperimentID", mock.Anything, resultID, experimentID).
					Return(CreateTestResult(experimentID), nil)
				// Update is still called with empty fields (only version incremented)
				repo.On("Update", mock.Anything, resultID, experimentID, mock.AnythingOfType("*result.UpdateFields"), 1).
					Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := resultmocks.NewRepository(t)
			mockPermSvc := permissionmocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := result.NewService(mockRepo, mockPermSvc)

			err := svc.UpdateResult(tc.ctx, resultID, experimentID, tc.dto)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
