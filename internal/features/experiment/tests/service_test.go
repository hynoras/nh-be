package experiment

import (
	"context"
	"errors"
	"testing"

	"nh-be/internal/constant"
	"nh-be/internal/features/experiment"
	"nh-be/internal/features/procedure"
	"nh-be/internal/features/experiment/mocks"
	permmocks "nh-be/internal/features/permission/mocks"
	procmocks "nh-be/internal/features/procedure/mocks"
	"nh-be/internal/utils/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_AssignProcedureToExperiment(t *testing.T) {
	userID := uuid.New()
	experimentID := uuid.New()
	procedureID := uuid.New()
	existingProcedureID := uuid.New()
	version := 1

	tests := []struct {
		name          string
		ctx           context.Context
		experimentID  uuid.UUID
		procedureID   uuid.UUID
		version       int
		setupMocks    func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name:         "success",
			ctx:          testutil.ContextWithUser(userID),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				procSvc.On("GetProcedureByID", mock.Anything, procedureID).
					Return(nil, nil)
				repo.On("GetProcedureIDByID", mock.Anything, experimentID).
					Return(existingProcedureID, nil)
				repo.On("UpdateProcedureID", mock.Anything, experimentID, procedureID, version).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:         "permission_denied",
			ctx:          testutil.ContextWithUser(userID),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
			},
			expectedError: experiment.ErrForbidUpdateExperiment,
		},
		{
			name:         "get_user_id_error",
			ctx:          context.Background(),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks:   func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name:         "get_permission_error",
			ctx:          testutil.ContextWithUser(userID),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, errors.New("permission error"))
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "permission error")
			},
		},
		{
			name:         "procedure_not_found",
			ctx:          testutil.ContextWithUser(userID),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				procSvc.On("GetProcedureByID", mock.Anything, procedureID).
					Return(nil, procedure.ErrProcedureNotFound)
			},
			expectedError: procedure.ErrProcedureNotFound,
		},
		{
			name:         "experiment_not_found",
			ctx:          testutil.ContextWithUser(userID),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				procSvc.On("GetProcedureByID", mock.Anything, procedureID).
					Return(nil, nil)
				repo.On("GetProcedureIDByID", mock.Anything, experimentID).
					Return(uuid.Nil, experiment.ErrExperimentNotFound)
			},
			expectedError: experiment.ErrExperimentNotFound,
		},
		{
			name:         "duplicate_assignment",
			ctx:          testutil.ContextWithUser(userID),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				procSvc.On("GetProcedureByID", mock.Anything, procedureID).
					Return(nil, nil)
				repo.On("GetProcedureIDByID", mock.Anything, experimentID).
					Return(procedureID, nil)
			},
			expectedError: experiment.ErrDuplicateProcedureAssignment,
		},
		{
			name:         "version_conflict",
			ctx:          testutil.ContextWithUser(userID),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				procSvc.On("GetProcedureByID", mock.Anything, procedureID).
					Return(nil, nil)
				repo.On("GetProcedureIDByID", mock.Anything, experimentID).
					Return(existingProcedureID, nil)
				repo.On("UpdateProcedureID", mock.Anything, experimentID, procedureID, version).
					Return(experiment.ErrExperimentConflict)
			},
			expectedError: experiment.ErrExperimentConflict,
		},
		{
			name:         "update_db_error",
			ctx:          testutil.ContextWithUser(userID),
			experimentID: experimentID,
			procedureID:  procedureID,
			version:      version,
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service, procSvc *procmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				procSvc.On("GetProcedureByID", mock.Anything, procedureID).
					Return(nil, nil)
				repo.On("GetProcedureIDByID", mock.Anything, experimentID).
					Return(existingProcedureID, nil)
				repo.On("UpdateProcedureID", mock.Anything, experimentID, procedureID, version).
					Return(errors.New("db error"))
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "db error")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockPermSvc := permmocks.NewService(t)
			mockProcSvc := procmocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc, mockProcSvc)

			svc := experiment.NewService(mockRepo, mockPermSvc, mockProcSvc)

			err := svc.AssignProcedureToExperiment(tc.ctx, tc.experimentID, tc.procedureID, tc.version)

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

func TestService_GetAllExperiments(t *testing.T) {
	userID := uuid.New()

	queryDtos := TestExperimentsQueryDto()
	expectedDtos := TestExperimentsResponseDto()
	expectedCount := int64(len(queryDtos))

	tests := []struct {
		name          string
		ctx           context.Context
		setupMocks    func(repo *mocks.Repository, permSvc *permmocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
		checkResult   func(t *testing.T, result []experiment.ExperimentsResponseDto, count int64)
	}{
		{
			name: "permission_denied",
			ctx:  testutil.ContextWithUser(userID),
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
			},
			expectedError: experiment.ErrForbidViewExperiments,
		},
		{
			name: "get_user_id_from_context_returns_error",
			ctx:  context.Background(), // no user ID in context — RequirePermission short-circuits here
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				// no mock needed: RequirePermission reads the user ID from ctx before calling permSvc
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name: "count_experiments_returns_error",
			ctx:  testutil.ContextWithUser(userID),
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("CountExperiments", mock.Anything, &userID).
					Return(int64(0), errors.New("db count error"))
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "db count error")
			},
		},
		{
			name: "find_all_experiments_returns_error",
			ctx:  testutil.ContextWithUser(userID),
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("CountExperiments", mock.Anything, &userID).
					Return(expectedCount, nil)
				repo.On("FindAllExperiments", mock.Anything,
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("constant.Order"),
					mock.Anything, mock.Anything, &userID,
					mock.AnythingOfType("int"),
					mock.AnythingOfType("int"),
				).Return(nil, errors.New("db find error"))
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "db find error")
			},
		},
		{
			name: "success",
			ctx:  testutil.ContextWithUser(userID),
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("CountExperiments", mock.Anything, &userID).
					Return(expectedCount, nil)
				repo.On("FindAllExperiments", mock.Anything,
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("constant.Order"),
					mock.Anything, mock.Anything, &userID,
					mock.AnythingOfType("int"),
					mock.AnythingOfType("int"),
				).Return(queryDtos, nil)
			},
			checkResult: func(t *testing.T, result []experiment.ExperimentsResponseDto, count int64) {
				assert.Equal(t, expectedDtos, result)
				assert.Equal(t, expectedCount, count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockPermSvc := permmocks.NewService(t)
			mockProcSvc := procmocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := experiment.NewService(mockRepo, mockPermSvc, mockProcSvc)

			result, count, err := svc.GetAllExperiments(tc.ctx, "created_at", "", constant.DESC, nil, nil, 1, 10)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, result)
				assert.Equal(t, int64(0), count)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
				assert.Nil(t, result)
				assert.Equal(t, int64(0), count)
			} else {
				assert.NoError(t, err)
			}

			if tc.checkResult != nil {
				tc.checkResult(t, result, count)
			}
		})
	}
}

func TestService_GetExperimentDetail(t *testing.T) {
	userID := uuid.New()
	experimentID := "EXP-0001"
	exp := TestExperiment()
	expectedDto := TestExperimentDetailResponseDto()

	tests := []struct {
		name          string
		ctx           context.Context
		setupMocks    func(repo *mocks.Repository, permSvc *permmocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
		checkResult   func(t *testing.T, result *experiment.ExperimentResponseDto)
	}{
		{
			name: "success",
			ctx:  testutil.ContextWithUser(userID),
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindByIdentifierAndCreatedBy", mock.Anything, experimentID, userID).
					Return(&exp, nil)
			},
			checkResult: func(t *testing.T, result *experiment.ExperimentResponseDto) {
				assert.NotNil(t, result)
				assert.Equal(t, &expectedDto, result)
			},
		},
		{
			name: "permission_denied",
			ctx:  testutil.ContextWithUser(userID),
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
			},
			expectedError: experiment.ErrForbidViewExperiment,
		},
		{
			name: "get_user_id_failed",
			ctx:  context.Background(),
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				// No mock needed; RequirePermission short-circuits.
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name: "repository_error",
			ctx:  testutil.ContextWithUser(userID),
			setupMocks: func(repo *mocks.Repository, permSvc *permmocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindByIdentifierAndCreatedBy", mock.Anything, experimentID, userID).
					Return(nil, errors.New("db find error"))
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "db find error")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockPermSvc := permmocks.NewService(t)
			mockProcSvc := procmocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := experiment.NewService(mockRepo, mockPermSvc, mockProcSvc)

			result, err := svc.GetExperimentDetail(tc.ctx, experimentID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, result)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
			}

			if tc.checkResult != nil {
				tc.checkResult(t, result)
			}
		})
	}
}
