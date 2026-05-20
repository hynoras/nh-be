package procedure

import (
	"context"
	"errors"
	"testing"
	"time"

	"nh-be/internal/constant"
	"nh-be/internal/features/permission/mocks"
	"nh-be/internal/features/procedure"
	procmocks "nh-be/internal/features/procedure/mocks"
	"nh-be/internal/utils/stringutil"

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
				repo.On("FindAll", mock.Anything, "", 0, 10).
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
				repo.On("FindAll", mock.Anything, "", 0, 10).
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
			expectedError:  procedure.ErrForbidViewProcedure,
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
				repo.On("FindAll", mock.Anything, "", 0, 10).
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
				repo.On("FindAll", mock.Anything, "test", 0, 10).
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
				repo.On("FindAll", mock.Anything, "", 20, 50).
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

func TestService_GetProcedureByID(t *testing.T) {
	userID := uuid.New()
	procedureID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	permissionError := errors.New("permission denied")
	genericRepoError := errors.New("database error")

	tests := []struct {
		name          string
		ctx           context.Context
		procedureID   uuid.UUID
		setupMocks    func(repo *procmocks.Repository, permSvc *mocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
		checkResult   func(t *testing.T, result *procedure.ProcedureResponseDto)
	}{
		{
			name:        "mapping_success",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				testProc := TestProcedureDetailWithRelations()
				repo.On("FindByID", mock.Anything, procedureID, true).
					Return(&testProc, nil)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			checkResult: func(t *testing.T, result *procedure.ProcedureResponseDto) {
				assert.NotNil(t, result)
				assert.Equal(t, "12345678-1234-1234-1234-123456789012", result.ID)
				assert.Equal(t, "Test Procedure", result.Title)
				assert.Equal(t, "Test Procedure Description", result.Description)

				assert.Equal(t, []procedure.StepsResponseDto{
					{
						ID:          "33333333-1234-1234-1234-444433332222",
						Index:       1,
						Title:       "Test Step",
						Description: stringutil.StringPtr("Test Step Description"),
						IsOptional:  false,
						StepType:    "wait",
						CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
						UpdatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
					},
					{
						ID:          "33333333-1234-1234-1234-555533332222",
						Index:       2,
						Title:       "Test Step 2",
						Description: stringutil.StringPtr("Test Step Description 2"),
						IsOptional:  true,
						StepType:    "decision",
						CreatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
						UpdatedAt:   time.Date(2026, 2, 3, 19, 15, 10, 0, time.UTC),
					},
				}, result.Steps)
			},
		},
		{
			name:        "repository_called_with_correct_flags",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				testProc := TestProcedureDetailWithRelations()
				repo.On("FindByID", mock.Anything, procedureID, true).
					Return(&testProc, nil)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			checkResult: func(t *testing.T, result *procedure.ProcedureResponseDto) {
				assert.NotNil(t, result)
			},
		},
		{
			name:        "permission_denied",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, permissionError)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, permissionError)
			},
			checkResult: func(t *testing.T, result *procedure.ProcedureResponseDto) {
				assert.Nil(t, result)
			},
		},
		{
			name:        "procedure_found",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				testProc := TestProcedureDetailWithRelations()
				repo.On("FindByID", mock.Anything, procedureID, true).
					Return(&testProc, nil)
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			checkResult: func(t *testing.T, result *procedure.ProcedureResponseDto) {
				assert.NotNil(t, result)
				assert.Equal(t, procedureID.String(), result.ID)
			},
		},
		{
			name:        "procedure_not_found",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindByID", mock.Anything, procedureID, true).
					Return(nil, procedure.ErrProcedureNotFound)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, procedure.ErrProcedureNotFound)
			},
			checkResult: func(t *testing.T, result *procedure.ProcedureResponseDto) {
				assert.Nil(t, result)
			},
		},
		{
			name:        "repository_error",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("FindByID", mock.Anything, procedureID, true).
					Return(nil, genericRepoError)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, genericRepoError)
			},
			checkResult: func(t *testing.T, result *procedure.ProcedureResponseDto) {
				assert.Nil(t, result)
			},
		},
		{
			name:        "context_cancelled_before_permission",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, context.Canceled)
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, context.Canceled)
			},
			checkResult: func(t *testing.T, result *procedure.ProcedureResponseDto) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := procmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := procedure.NewService(mockRepo, mockPermSvc)

			result, err := svc.GetProcedureByID(tc.ctx, tc.procedureID)

			tc.checkError(t, err)
			tc.checkResult(t, result)
		})
	}
}

func TestService_CreateProcedure(t *testing.T) {
	userID := uuid.New()
	permissionError := errors.New("permission service unavailable")
	repoError := errors.New("database error")

	tests := []struct {
		name          string
		ctx           context.Context
		dto           *procedure.CreateProcedureDto
		setupMocks    func(repo *procmocks.Repository, permSvc *mocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "success_create",
			ctx:  ContextWithUser(userID),
			dto:  TestCreateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*procedure.Procedure")).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "forbidden_create",
			ctx:  ContextWithUser(userID),
			dto:  TestCreateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				// Repo should NOT be called
			},
			expectedError: procedure.ErrForbidCreateProcedure,
		},
		{
			name: "missing_user_in_context",
			ctx:  context.Background(), // No user in context
			dto:  TestCreateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				// Permission service and repo should NOT be called
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name: "permission_service_error",
			ctx:  ContextWithUser(userID),
			dto:  TestCreateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, permissionError)
				// Repo should NOT be called
			},
			expectedError: permissionError,
		},
		{
			name: "repo_error_propagated",
			ctx:  ContextWithUser(userID),
			dto:  TestCreateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*procedure.Procedure")).
					Return(repoError)
			},
			expectedError: repoError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := procmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := procedure.NewService(mockRepo, mockPermSvc)

			err := svc.CreateProcedure(tc.ctx, tc.dto)

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

func TestService_UpdateProcedure(t *testing.T) {
	userID := uuid.New()
	procedureID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	dbError := errors.New("database error")

	tests := []struct {
		name          string
		ctx           context.Context
		procedureID   uuid.UUID
		dto           *procedure.UpdateProcedureDto
		setupMocks    func(repo *procmocks.Repository, permSvc *mocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name:        "permission_denied",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			dto:         TestUpdateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				// Repo should NOT be called
			},
			expectedError: procedure.ErrForbidUpdateProcedure,
		},
		{
			name:        "update_success",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			dto:         TestUpdateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("UpdateProcedure", mock.Anything, procedureID, mock.AnythingOfType("*procedure.Procedure")).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "repo_returns_not_found",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			dto:         TestUpdateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("UpdateProcedure", mock.Anything, procedureID, mock.AnythingOfType("*procedure.Procedure")).
					Return(procedure.ErrProcedureNotFound)
			},
			expectedError: procedure.ErrProcedureNotFound,
		},
		{
			name:        "repo_returns_conflict",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			dto:         TestUpdateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("UpdateProcedure", mock.Anything, procedureID, mock.AnythingOfType("*procedure.Procedure")).
					Return(procedure.ErrProcedureConflict)
			},
			expectedError: procedure.ErrProcedureConflict,
		},
		{
			name:        "repo_returns_db_error",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			dto:         TestUpdateProcedureDto(),
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("UpdateProcedure", mock.Anything, procedureID, mock.AnythingOfType("*procedure.Procedure")).
					Return(dbError)
			},
			expectedError: dbError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := procmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := procedure.NewService(mockRepo, mockPermSvc)

			err := svc.UpdateProcedure(tc.ctx, tc.procedureID, tc.dto)

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

func TestService_UpdateProcedureStep(t *testing.T) {
	userID := uuid.New()
	procedureID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := []struct {
		name          string
		ctx           context.Context
		procedureID   uuid.UUID
		stepInput     []procedure.UpdateProcedureStepInput
		setupMocks    func(repo *procmocks.Repository, permSvc *mocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name:        "permission_denied_returns_error",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput:   []procedure.UpdateProcedureStepInput{TestUpdateStepInput(ExistingStepID1)},
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				// Repo should NOT be called
			},
			expectedError: procedure.ErrForbidUpdateProcedure,
		},
		{
			name:        "transaction_error_is_propagated",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput:   []procedure.UpdateProcedureStepInput{TestUpdateStepInput(ExistingStepID1)},
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(procedure.Repository) error")).
					Return(errors.New("transaction failed"))
			},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "transaction failed")
			},
		},
		{
			name:        "create_new_step_when_id_is_nil",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput:   []procedure.UpdateProcedureStepInput{TestNewStepInput()},
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("GetStepIDsByProcID", mock.Anything, procedureID).
					Return([]procedure.StepMetadata{}, nil)
				repo.On("CreateProcedureStep", mock.Anything, mock.AnythingOfType("*procedure.ProcedureStep")).
					Return(nil)
				repo.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(procedure.Repository) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(procedure.Repository) error)
						_ = fn(repo)
					}).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "update_existing_step_success",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput:   []procedure.UpdateProcedureStepInput{TestUpdateStepInput(ExistingStepID1)},
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("GetStepIDsByProcID", mock.Anything, procedureID).
					Return(TestExistingStepMetadata(), nil)
				repo.On("UpdateProcedureStep", mock.Anything, ExistingStepID1, procedureID, mock.AnythingOfType("*procedure.ProcedureStep")).
					Return(nil)
				repo.On("DeleteProcedureStep", mock.Anything, ExistingStepID2, procedureID).
					Return(nil)
				repo.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(procedure.Repository) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(procedure.Repository) error)
						_ = fn(repo)
					}).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "update_non_existing_step_returns_not_found",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput:   []procedure.UpdateProcedureStepInput{TestUpdateStepInput(ExistingStepID3)}, // ID not in existing
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("GetStepIDsByProcID", mock.Anything, procedureID).
					Return(TestExistingStepMetadata(), nil)
				repo.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(procedure.Repository) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(procedure.Repository) error)
						_ = fn(repo)
					}).
					Return(procedure.ErrProcedureStepNotFound)
			},
			expectedError: procedure.ErrProcedureStepNotFound,
		},
		{
			name:        "delete_missing_step_when_not_in_input",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput:   []procedure.UpdateProcedureStepInput{}, // empty input, all existing should be deleted
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("GetStepIDsByProcID", mock.Anything, procedureID).
					Return([]procedure.StepMetadata{{ID: ExistingStepID1, Version: 1}}, nil)
				repo.On("DeleteProcedureStep", mock.Anything, ExistingStepID1, procedureID).
					Return(nil)
				repo.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(procedure.Repository) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(procedure.Repository) error)
						_ = fn(repo)
					}).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "version_conflict_is_propagated",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput:   []procedure.UpdateProcedureStepInput{TestUpdateStepInput(ExistingStepID1)},
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("GetStepIDsByProcID", mock.Anything, procedureID).
					Return(TestExistingStepMetadata(), nil)
				repo.On("UpdateProcedureStep", mock.Anything, ExistingStepID1, procedureID, mock.AnythingOfType("*procedure.ProcedureStep")).
					Return(procedure.ErrProcedureConflict)
				repo.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(procedure.Repository) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(procedure.Repository) error)
						_ = fn(repo)
					}).
					Return(procedure.ErrProcedureConflict)
			},
			expectedError: procedure.ErrProcedureConflict,
		},
		{
			name:        "mixed_create_update_delete_in_single_request",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput: []procedure.UpdateProcedureStepInput{
				TestNewStepInput(),                   // create (uuid.Nil)
				TestUpdateStepInput(ExistingStepID1), // update (existing)
				// ExistingStepID2 not in input → delete
			},
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("GetStepIDsByProcID", mock.Anything, procedureID).
					Return(TestExistingStepMetadata(), nil)
				repo.On("CreateProcedureStep", mock.Anything, mock.AnythingOfType("*procedure.ProcedureStep")).
					Return(nil)
				repo.On("UpdateProcedureStep", mock.Anything, ExistingStepID1, procedureID, mock.AnythingOfType("*procedure.ProcedureStep")).
					Return(nil)
				repo.On("DeleteProcedureStep", mock.Anything, ExistingStepID2, procedureID).
					Return(nil)
				repo.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(procedure.Repository) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(procedure.Repository) error)
						_ = fn(repo)
					}).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "empty_input_deletes_all_existing_steps",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			stepInput:   []procedure.UpdateProcedureStepInput{},
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("GetStepIDsByProcID", mock.Anything, procedureID).
					Return(TestExistingStepMetadata(), nil)
				repo.On("DeleteProcedureStep", mock.Anything, ExistingStepID1, procedureID).
					Return(nil)
				repo.On("DeleteProcedureStep", mock.Anything, ExistingStepID2, procedureID).
					Return(nil)
				repo.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(procedure.Repository) error")).
					Run(func(args mock.Arguments) {
						fn := args.Get(1).(func(procedure.Repository) error)
						_ = fn(repo)
					}).
					Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := procmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := procedure.NewService(mockRepo, mockPermSvc)

			err := svc.UpdateProcedureStep(tc.ctx, tc.procedureID, tc.stepInput)

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

func TestService_DeleteProcedure(t *testing.T) {
	userID := uuid.New()
	procedureID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	dbError := errors.New("database error")
	permissionError := errors.New("permission service unavailable")

	tests := []struct {
		name          string
		ctx           context.Context
		procedureID   uuid.UUID
		setupMocks    func(repo *procmocks.Repository, permSvc *mocks.Service)
		expectedError error
		checkError    func(t *testing.T, err error)
	}{
		{
			name:        "success",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("DeleteProcedure", mock.Anything, procedureID).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "permission_denied",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.AssertNotCalled(t, "DeleteProcedure", mock.Anything, procedureID)
			},
			expectedError: procedure.ErrForbidDeleteProcedure,
		},
		{
			name:        "permission_service_error",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, permissionError)
				repo.AssertNotCalled(t, "DeleteProcedure", mock.Anything, procedureID)
			},
			expectedError: permissionError,
		},
		{
			name:        "repository_not_found",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("DeleteProcedure", mock.Anything, procedureID).
					Return(procedure.ErrProcedureNotFound)
			},
			expectedError: procedure.ErrProcedureNotFound,
		},
		{
			name:        "repository_db_error",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("DeleteProcedure", mock.Anything, procedureID).
					Return(dbError)
			},
			expectedError: dbError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := procmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := procedure.NewService(mockRepo, mockPermSvc)

			err := svc.DeleteProcedure(tc.ctx, tc.procedureID)

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

func TestService_GetProcedureSteps(t *testing.T) {
	userID := uuid.New()
	procedureID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	repoError := errors.New("repository error")

	tests := []struct {
		name           string
		ctx            context.Context
		procedureID    uuid.UUID
		offset         int
		limit          int
		setupMocks     func(repo *procmocks.Repository, permSvc *mocks.Service)
		expectedResult []procedure.StepsResponseDto
		expectedLength int64
		expectedError  error
		checkError     func(t *testing.T, err error)
	}{
		{
			name:        "permission_denied",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			offset:      0,
			limit:       10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
			},
			expectedError: procedure.ErrForbidViewProcedure,
		},
		{
			name:        "repo_returns_error",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			offset:      0,
			limit:       10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetProcStepsByProcID", mock.Anything, procedureID, 0, 10).
					Return(nil, int64(0), repoError)
			},
			expectedError: repoError,
		},
		{
			name:        "success_no_steps",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			offset:      0,
			limit:       10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetProcStepsByProcID", mock.Anything, procedureID, 0, 10).
					Return([]procedure.ProcedureStep{}, int64(0), nil)
			},
			expectedResult: []procedure.StepsResponseDto{},
			expectedLength: 0,
			expectedError:  nil,
		},
		{
			name:        "success_with_steps",
			ctx:         ContextWithUser(userID),
			procedureID: procedureID,
			offset:      0,
			limit:       10,
			setupMocks: func(repo *procmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				steps := TestStepsForProcedure(procedureID, 2)
				repo.On("GetProcStepsByProcID", mock.Anything, procedureID, 0, 10).
					Return(steps, int64(2), nil)
			},
			expectedResult: procedure.MapStepsToDto(TestStepsForProcedure(procedureID, 2)),
			expectedLength: 2,
			expectedError:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := procmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := procedure.NewService(mockRepo, mockPermSvc)

			res, length, err := svc.GetProcedureSteps(tc.ctx, tc.procedureID, tc.offset, tc.limit)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, res)
				assert.Equal(t, tc.expectedLength, length)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, len(tc.expectedResult), len(res))
				assert.Equal(t, tc.expectedLength, length)
			}
		})
	}
}
