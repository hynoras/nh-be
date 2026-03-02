package observation

import (
	"context"
	"errors"
	"testing"

	"nh-be/internal/constant"
	"nh-be/internal/features/observation"
	obsmocks "nh-be/internal/features/observation/mocks"
	"nh-be/internal/features/permission/mocks"
	"nh-be/internal/utils/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetAllObservations(t *testing.T) {
	userID := uuid.New()
	expID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	procID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name           string
		ctx            context.Context
		expId          uuid.UUID
		procId         uuid.UUID
		offset         int
		limit          int
		sortBy         *string
		sortOrder      *constant.Order
		setupMocks     func(repo *obsmocks.Repository, permSvc *mocks.Service)
		expectedResult []observation.ObservationsResponseDto
		expectedLength int64
		expectedError  error
		checkError     func(t *testing.T, err error)
	}{
		{
			name:   "forbidden_view_no_permission",
			ctx:    testutil.ContextWithUser(userID),
			expId:  expID,
			procId: procID,
			offset: 0,
			limit:  10,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
				repo.AssertNotCalled(t, "GetAllObsByExpIDAndProcID",
					mock.Anything, expID, procID, 0, 10, (*string)(nil), (*constant.Order)(nil))
			},
			expectedResult: nil,
			expectedLength: 0,
			expectedError:  constant.ErrForbidViewObservation,
		},
		{
			name:   "forbidden_view_permission_error",
			ctx:    testutil.ContextWithUser(userID),
			expId:  expID,
			procId: procID,
			offset: 0,
			limit:  10,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return(nil, errors.New("permission service unavailable"))
				repo.AssertNotCalled(t, "GetAllObsByExpIDAndProcID",
					mock.Anything, expID, procID, 0, 10, (*string)(nil), (*constant.Order)(nil))
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
			ctx:    context.Background(),
			expId:  expID,
			procId: procID,
			offset: 0,
			limit:  10,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.AssertNotCalled(t, "GetUserPermissionCodeNames", mock.Anything, userID)
				repo.AssertNotCalled(t, "GetAllObsByExpIDAndProcID",
					mock.Anything, expID, procID, 0, 10, (*string)(nil), (*constant.Order)(nil))
			},
			expectedResult: nil,
			expectedLength: 0,
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name:   "repo_error_propagated",
			ctx:    testutil.ContextWithUser(userID),
			expId:  expID,
			procId: procID,
			offset: 0,
			limit:  10,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetAllObsByExpIDAndProcID", mock.Anything, expID, procID, 0, 10, (*string)(nil), (*constant.Order)(nil)).
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
			name:   "empty_result_no_observations",
			ctx:    testutil.ContextWithUser(userID),
			expId:  expID,
			procId: procID,
			offset: 0,
			limit:  10,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetAllObsByExpIDAndProcID", mock.Anything, expID, procID, 0, 10, (*string)(nil), (*constant.Order)(nil)).
					Return([]observation.ObservationMetadata{}, int64(0), nil)
			},
			expectedResult: []observation.ObservationsResponseDto{},
			expectedLength: 0,
			expectedError:  nil,
		},
		{
			name:   "success_single_page_no_sort",
			ctx:    testutil.ContextWithUser(userID),
			expId:  expID,
			procId: procID,
			offset: 0,
			limit:  10,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetAllObsByExpIDAndProcID", mock.Anything, expID, procID, 0, 10, (*string)(nil), (*constant.Order)(nil)).
					Return(TestObservationMetadataList(3), int64(3), nil)
			},
			expectedResult: observation.MapObservationsMetadataToDto(TestObservationMetadataList(3)),
			expectedLength: 3,
			expectedError:  nil,
		},
		{
			name:   "success_with_manage_permission",
			ctx:    testutil.ContextWithUser(userID),
			expId:  expID,
			procId: procID,
			offset: 0,
			limit:  10,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("GetAllObsByExpIDAndProcID", mock.Anything, expID, procID, 0, 10, (*string)(nil), (*constant.Order)(nil)).
					Return(TestObservationMetadataList(2), int64(2), nil)
			},
			expectedResult: observation.MapObservationsMetadataToDto(TestObservationMetadataList(2)),
			expectedLength: 2,
			expectedError:  nil,
		},
		{
			name:   "success_with_pagination",
			ctx:    testutil.ContextWithUser(userID),
			expId:  expID,
			procId: procID,
			offset: 2,
			limit:  1,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetAllObsByExpIDAndProcID", mock.Anything, expID, procID, 2, 1, (*string)(nil), (*constant.Order)(nil)).
					Return(TestObservationMetadataList(1), int64(5), nil)
			},
			expectedResult: observation.MapObservationsMetadataToDto(TestObservationMetadataList(1)),
			expectedLength: 5,
			expectedError:  nil,
		},
		{
			name:      "success_sort_by_created_at_asc",
			ctx:       testutil.ContextWithUser(userID),
			expId:     expID,
			procId:    procID,
			offset:    0,
			limit:     10,
			sortBy:    func() *string { s := "created_at"; return &s }(),
			sortOrder: func() *constant.Order { o := constant.ASC; return &o }(),
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				sortBy := "created_at"
				sortOrder := constant.ASC
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetAllObsByExpIDAndProcID", mock.Anything, expID, procID, 0, 10, &sortBy, &sortOrder).
					Return(TestObservationMetadataList(2), int64(2), nil)
			},
			expectedResult: observation.MapObservationsMetadataToDto(TestObservationMetadataList(2)),
			expectedLength: 2,
			expectedError:  nil,
		},
		{
			name:      "success_sort_by_created_at_desc",
			ctx:       testutil.ContextWithUser(userID),
			expId:     expID,
			procId:    procID,
			offset:    0,
			limit:     10,
			sortBy:    func() *string { s := "created_at"; return &s }(),
			sortOrder: func() *constant.Order { o := constant.DESC; return &o }(),
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				sortBy := "created_at"
				sortOrder := constant.DESC
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetAllObsByExpIDAndProcID", mock.Anything, expID, procID, 0, 10, &sortBy, &sortOrder).
					Return(TestObservationMetadataList(2), int64(2), nil)
			},
			expectedResult: observation.MapObservationsMetadataToDto(TestObservationMetadataList(2)),
			expectedLength: 2,
			expectedError:  nil,
		},
		{
			name:      "invalid_sort_column_ignored",
			ctx:       testutil.ContextWithUser(userID),
			expId:     expID,
			procId:    procID,
			offset:    0,
			limit:     10,
			sortBy:    func() *string { s := "title"; return &s }(),
			sortOrder: func() *constant.Order { o := constant.ASC; return &o }(),
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				sortBy := "title"
				sortOrder := constant.ASC
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ViewExperiment}, nil)
				repo.On("GetAllObsByExpIDAndProcID", mock.Anything, expID, procID, 0, 10, &sortBy, &sortOrder).
					Return(TestObservationMetadataList(1), int64(1), nil)
			},
			expectedResult: observation.MapObservationsMetadataToDto(TestObservationMetadataList(1)),
			expectedLength: 1,
			expectedError:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := obsmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := observation.NewService(mockRepo, mockPermSvc)

			res, length, err := svc.GetAllObservations(tc.ctx, tc.expId, tc.procId, tc.offset, tc.limit, tc.sortBy, tc.sortOrder)

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

func TestService_CreateObservation(t *testing.T) {
	userID := uuid.New()
	expID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	procStepID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	input := observation.CreateObservationInput{
		ObservedAt: TestObservation().ObservedAt,
		Title:      TestObservation().Title,
		Notes:      TestObservation().Notes,
	}

	tests := []struct {
		name           string
		ctx            context.Context
		expId          uuid.UUID
		procStepId     uuid.UUID
		input          observation.CreateObservationInput
		setupMocks     func(repo *obsmocks.Repository, permSvc *mocks.Service)
		expectedResult observation.CreatedObservationResponseDto
		expectedError  error
		checkError     func(t *testing.T, err error)
	}{
		{
			name:       "success",
			ctx:        testutil.ContextWithUser(userID),
			expId:      expID,
			procStepId: procStepID,
			input:      input,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				createdObs := TestObservation()
				repo.On("CreateObservation", mock.Anything, expID, procStepID, mock.AnythingOfType("observation.Observation")).
					Return(createdObs, nil)
			},
			expectedResult: observation.MapObsToCreatedObsResponseDto(TestObservation()),
			expectedError:  nil,
		},
		{
			name:       "permission_denied",
			ctx:        testutil.ContextWithUser(userID),
			expId:      expID,
			procStepId: procStepID,
			input:      input,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{}, nil)
				repo.AssertNotCalled(t, "CreateObservation",
					mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
			expectedResult: observation.CreatedObservationResponseDto{},
			expectedError:  constant.ErrForbidCreateObservation,
		},
		{
			name:       "missing_user_id_in_context",
			ctx:        context.Background(),
			expId:      expID,
			procStepId: procStepID,
			input:      input,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.AssertNotCalled(t, "GetUserPermissionCodeNames", mock.Anything, mock.Anything)
				repo.AssertNotCalled(t, "CreateObservation",
					mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
			expectedResult: observation.CreatedObservationResponseDto{},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user ID not found in context")
			},
		},
		{
			name:       "experiment_not_found",
			ctx:        testutil.ContextWithUser(userID),
			expId:      expID,
			procStepId: procStepID,
			input:      input,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("CreateObservation", mock.Anything, expID, procStepID, mock.AnythingOfType("observation.Observation")).
					Return(observation.Observation{}, constant.ErrExperimentNotFound)
			},
			expectedResult: observation.CreatedObservationResponseDto{},
			expectedError:  constant.ErrExperimentNotFound,
		},
		{
			name:       "procedure_not_found",
			ctx:        testutil.ContextWithUser(userID),
			expId:      expID,
			procStepId: procStepID,
			input:      input,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("CreateObservation", mock.Anything, expID, procStepID, mock.AnythingOfType("observation.Observation")).
					Return(observation.Observation{}, constant.ErrProcedureNotFound)
			},
			expectedResult: observation.CreatedObservationResponseDto{},
			expectedError:  constant.ErrProcedureNotFound,
		},
		{
			name:       "unexpected_error",
			ctx:        testutil.ContextWithUser(userID),
			expId:      expID,
			procStepId: procStepID,
			input:      input,
			setupMocks: func(repo *obsmocks.Repository, permSvc *mocks.Service) {
				permSvc.On("GetUserPermissionCodeNames", mock.Anything, userID).
					Return([]string{constant.ManageExperiment}, nil)
				repo.On("CreateObservation", mock.Anything, expID, procStepID, mock.AnythingOfType("observation.Observation")).
					Return(observation.Observation{}, errors.New("database connection error"))
			},
			expectedResult: observation.CreatedObservationResponseDto{},
			checkError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database connection error")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := obsmocks.NewRepository(t)
			mockPermSvc := mocks.NewService(t)

			tc.setupMocks(mockRepo, mockPermSvc)

			svc := observation.NewService(mockRepo, mockPermSvc)

			res, err := svc.CreateObservation(tc.ctx, tc.expId, tc.procStepId, tc.input)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Equal(t, tc.expectedResult, res)
			} else if tc.checkError != nil {
				tc.checkError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			}
		})
	}
}
