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
