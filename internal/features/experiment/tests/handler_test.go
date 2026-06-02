package experiment

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"nh-be/internal/constant"
	"nh-be/internal/features/experiment"
	experimentmocks "nh-be/internal/features/experiment/mocks"
	"nh-be/internal/utils/testutil"
)

func TestHandler_AssignProcedureToExperiment(t *testing.T) {
	validExpID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	validProcID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	tests := []struct {
		name           string
		experimentID   string
		procedureID    string
		requestBody    string
		setupMocks     func(svc *experimentmocks.Service)
		expectedStatus int
	}{
		{
			name:         "invalid_experiment_id",
			experimentID: "not-a-uuid",
			procedureID:  validProcID.String(),
			requestBody:  `{"version":1}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.AssertNotCalled(t, "AssignProcedureToExperiment")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "invalid_procedure_id",
			experimentID: validExpID.String(),
			procedureID:  "not-a-uuid",
			requestBody:  `{"version":1}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.AssertNotCalled(t, "AssignProcedureToExperiment")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "invalid_request_body",
			experimentID: validExpID.String(),
			procedureID:  validProcID.String(),
			requestBody:  `{}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.AssertNotCalled(t, "AssignProcedureToExperiment")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:         "success",
			experimentID: validExpID.String(),
			procedureID:  validProcID.String(),
			requestBody:  `{"version":1}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("AssignProcedureToExperiment", mock.Anything, validExpID, validProcID, 1).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "forbidden",
			experimentID: validExpID.String(),
			procedureID:  validProcID.String(),
			requestBody:  `{"version":1}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("AssignProcedureToExperiment", mock.Anything, validExpID, validProcID, 1).
					Return(experiment.ErrForbidUpdateExperiment)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "duplicate",
			experimentID: validExpID.String(),
			procedureID:  validProcID.String(),
			requestBody:  `{"version":1}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("AssignProcedureToExperiment", mock.Anything, validExpID, validProcID, 1).
					Return(experiment.ErrDuplicateProcedureAssignment)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:         "conflict",
			experimentID: validExpID.String(),
			procedureID:  validProcID.String(),
			requestBody:  `{"version":1}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("AssignProcedureToExperiment", mock.Anything, validExpID, validProcID, 1).
					Return(experiment.ErrExperimentConflict)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:         "experiment_not_found",
			experimentID: validExpID.String(),
			procedureID:  validProcID.String(),
			requestBody:  `{"version":1}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("AssignProcedureToExperiment", mock.Anything, validExpID, validProcID, 1).
					Return(experiment.ErrExperimentNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:         "internal_error",
			experimentID: validExpID.String(),
			procedureID:  validProcID.String(),
			requestBody:  `{"version":1}`,
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("AssignProcedureToExperiment", mock.Anything, validExpID, validProcID, 1).
					Return(errors.New("unexpected database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := experimentmocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := testutil.SetupTestRouter(
				http.MethodPut,
				"/experiments/:experimentId/procedures/:procedureId",
				experiment.AssignProcedureToExperimentHandler(mockSvc),
			)

			req := httptest.NewRequest(
				http.MethodPut,
				"/experiments/"+tc.experimentID+"/procedures/"+tc.procedureID,
				strings.NewReader(tc.requestBody),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandler_GetAllExperiments(t *testing.T) {
	queryDtos := TestExperimentsQueryDto()
	expectedResponseDtos := TestExperimentsResponseDto()
	expectedCount := int64(len(queryDtos))

	tests := []struct {
		name           string
		queryParams    string
		setupMocks     func(svc *experimentmocks.Service)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name:        "success_default_params",
			queryParams: "",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("GetAllExperiments",
					mock.Anything,
					"updated_at",
					"",
					constant.DESC,
					(*experiment.ExperimentStatus)(nil),
					(*experiment.ExperimentType)(nil),
					1,
					10,
				).Return(expectedResponseDtos, expectedCount, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "success_with_all_filters",
			queryParams: "?search=test&status=draft&type=exploratory&sortBy=title&sortOrder=ASC&page=2&pageSize=5",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("GetAllExperiments",
					mock.Anything,
					"title",
					"test",
					constant.ASC,
					mock.Anything, // *ExperimentStatus pointer created locally in handler
					mock.Anything, // *ExperimentType pointer created locally in handler
					2,
					5,
				).Return(expectedResponseDtos, expectedCount, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid_page",
			queryParams: "?page=invalid",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.AssertNotCalled(t, "GetAllExperiments")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid_page_size",
			queryParams: "?pageSize=invalid",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.AssertNotCalled(t, "GetAllExperiments")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid_sort_order",
			queryParams: "?sortOrder=invalid",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.AssertNotCalled(t, "GetAllExperiments")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid_status",
			queryParams: "?status=invalid",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.AssertNotCalled(t, "GetAllExperiments")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid_type",
			queryParams: "?type=invalid",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.AssertNotCalled(t, "GetAllExperiments")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "service_returns_forbidden",
			queryParams: "",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("GetAllExperiments", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, int64(0), experiment.ErrForbidViewExperiments)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "service_returns_internal_error",
			queryParams: "",
			setupMocks: func(svc *experimentmocks.Service) {
				svc.On("GetAllExperiments", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, int64(0), errors.New("unexpected db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := experimentmocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := testutil.SetupTestRouter(
				http.MethodGet,
				"/experiments",
				experiment.GetAllExperimentsHandler(mockSvc),
			)

			req := httptest.NewRequest(http.MethodGet, "/experiments"+tc.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
