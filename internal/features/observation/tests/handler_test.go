package observation

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"nh-be/internal/constant"
	"nh-be/internal/features/observation"
	observationmocks "nh-be/internal/features/observation/mocks"
	"nh-be/internal/utils/testutil"
)

func TestHandler_GetAllObservations(t *testing.T) {
	validExpId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	validProcStepId := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name            string
		experimentID    string
		procedureStepID string
		queryParams     string
		setupMocks      func(svc *observationmocks.Service)
		expectedStatus  int
	}{
		{
			name:            "success",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=1&pageSize=10&sortBy=created_at&sortOrder=DESC",
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("GetAllObservations", mock.Anything, validExpId, validProcStepId, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(TestObservationsResponseDtoList(3), int64(3), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:            "invalid_experiment_id",
			experimentID:    "invalid-id",
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=1&pageSize=10",
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "GetAllObservations")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "invalid_procedure_step_id",
			experimentID:    validExpId.String(),
			procedureStepID: "invalid-id",
			queryParams:     "?page=1&pageSize=10",
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "GetAllObservations")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "invalid_page",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=abc&pageSize=10",
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "GetAllObservations")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "invalid_page_size",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=1&pageSize=xyz",
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "GetAllObservations")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "invalid_sort_order",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?sortBy=created_at&sortOrder=WRONG",
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "GetAllObservations")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "permission_denied",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=1&pageSize=10",
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("GetAllObservations", mock.Anything, validExpId, validProcStepId, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, int64(0), constant.ErrForbidViewObservation)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:            "service_error",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=1&pageSize=10",
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("GetAllObservations", mock.Anything, validExpId, validProcStepId, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, int64(0), errors.New("unexpected database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:            "success_empty_result",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=1&pageSize=10",
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("GetAllObservations", mock.Anything, validExpId, validProcStepId, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return([]observation.ObservationsResponseDto{}, int64(0), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:            "success_partial_sort_params_missing_sortBy",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=1&pageSize=10&sortOrder=ASC",
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("GetAllObservations", mock.Anything, validExpId, validProcStepId, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(TestObservationsResponseDtoList(2), int64(2), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:            "success_partial_sort_params_missing_sort_order",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			queryParams:     "?page=1&pageSize=10&sortBy=created_at",
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("GetAllObservations", mock.Anything, validExpId, validProcStepId, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(TestObservationsResponseDtoList(2), int64(2), nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := observationmocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := testutil.SetupTestRouter(
				http.MethodGet,
				"/observations/:experimentId/:procedureStepId",
				observation.GetAllObservationsHandler(mockSvc),
			)

			req := httptest.NewRequest(
				http.MethodGet,
				"/observations/"+tc.experimentID+"/"+tc.procedureStepID+tc.queryParams,
				nil,
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandler_CreateObservation(t *testing.T) {
	validExpId := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	validProcStepId := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	validBody := `{"title":"Test Observation","notes":"Test Notes","observed_at":"2026-02-03T19:15:10Z"}`

	tests := []struct {
		name            string
		experimentID    string
		procedureStepID string
		requestBody     string
		setupMocks      func(svc *observationmocks.Service)
		expectedStatus  int
	}{
		{
			name:            "success",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     validBody,
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("CreateObservation", mock.Anything, validExpId, validProcStepId, mock.Anything).
					Return(TestCreatedObservationResponseDto(), nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:            "invalid_experiment_uuid",
			experimentID:    "not-a-uuid",
			procedureStepID: validProcStepId.String(),
			requestBody:     validBody,
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "CreateObservation")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "invalid_procedure_step_uuid",
			experimentID:    validExpId.String(),
			procedureStepID: "not-a-uuid",
			requestBody:     validBody,
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "CreateObservation")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "invalid_json_body",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     `{invalid json`,
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "CreateObservation")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "invalid_previous_observation_id",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     `{"title":"Test","observed_at":"2026-02-03T19:15:10Z","previous_observation_id":"not-a-uuid"}`,
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "CreateObservation")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:            "missing_required_fields",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     `{}`,
			setupMocks: func(svc *observationmocks.Service) {
				svc.AssertNotCalled(t, "CreateObservation")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:            "experiment_not_found",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     validBody,
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("CreateObservation", mock.Anything, validExpId, validProcStepId, mock.Anything).
					Return(observation.CreatedObservationResponseDto{}, constant.ErrExperimentNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "procedure_step_not_found",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     validBody,
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("CreateObservation", mock.Anything, validExpId, validProcStepId, mock.Anything).
					Return(observation.CreatedObservationResponseDto{}, constant.ErrProcedureStepNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "unauthorized",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     validBody,
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("CreateObservation", mock.Anything, validExpId, validProcStepId, mock.Anything).
					Return(observation.CreatedObservationResponseDto{}, constant.ErrForbidCreateObservation)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:            "internal_error",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     validBody,
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("CreateObservation", mock.Anything, validExpId, validProcStepId, mock.Anything).
					Return(observation.CreatedObservationResponseDto{}, errors.New("unexpected database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:            "context_cancellation",
			experimentID:    validExpId.String(),
			procedureStepID: validProcStepId.String(),
			requestBody:     validBody,
			setupMocks: func(svc *observationmocks.Service) {
				svc.On("CreateObservation", mock.Anything, validExpId, validProcStepId, mock.Anything).
					Return(observation.CreatedObservationResponseDto{}, context.Canceled)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := observationmocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := testutil.SetupTestRouter(
				http.MethodPost,
				"/observations/:experimentId/:procedureStepId",
				observation.CreateObservationHandler(mockSvc),
			)

			req := httptest.NewRequest(
				http.MethodPost,
				"/observations/"+tc.experimentID+"/"+tc.procedureStepID,
				strings.NewReader(tc.requestBody),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
