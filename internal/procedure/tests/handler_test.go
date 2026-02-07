package procedure

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"nh-be/constant"
	"nh-be/internal/procedure"
	proceduremocks "nh-be/internal/procedure/mocks"
	"nh-be/utils"
)

func TestHandler_GetAllProcedures(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		setupMocks     func(svc *proceduremocks.Service)
		expectedStatus int
	}{
		{
			name:        "success",
			queryParams: "?page=1&pageSize=10",
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("GetAllProcedures", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return([]procedure.ProcedureListResponseDto{}, int64(0), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid_pagination",
			queryParams: "?page=invalid&pageSize=10",
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "GetAllProcedures")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "forbidden_no_permission",
			queryParams: "?page=1&pageSize=10",
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("GetAllProcedures", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, int64(0), constant.ErrForbidViewProcedure)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "internal_server_error",
			queryParams: "?page=1&pageSize=10",
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("GetAllProcedures", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, int64(0), errors.New("unexpected database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := proceduremocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := utils.SetupTestRouter(http.MethodGet, "/procedures",
				procedure.GetAllProceduresHandler(mockSvc))

			req := httptest.NewRequest(http.MethodGet,
				"/procedures"+tc.queryParams,
				nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandler_GetProcedureByID(t *testing.T) {
	validID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	testProcedure := TestProcedureDetail()

	tests := []struct {
		name           string
		procedureID    string
		setupMocks     func(svc *proceduremocks.Service)
		expectedStatus int
	}{
		{
			name:        "invalid_uuid",
			procedureID: "abc",
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "GetProcedureByID")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "procedure_not_found",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("GetProcedureByID", mock.Anything, validID).
					Return(nil, constant.ErrProcedureNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "forbid_view_procedure",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("GetProcedureByID", mock.Anything, validID).
					Return(nil, constant.ErrForbidViewProcedure)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "optimistic_locking_conflict",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("GetProcedureByID", mock.Anything, validID).
					Return(nil, constant.ErrOptimisticLockingConflict)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "internal_error",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("GetProcedureByID", mock.Anything, validID).
					Return(nil, errors.New("db down"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "success",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				testResponseDto := &procedure.ProcedureResponseDto{
					ID:                validID.String(),
					Title:             "Test Procedure",
					Description:       "Test Procedure Description",
					UsedByExperiments: []procedure.UsedByExperiment{},
					Steps:             []procedure.Steps{},
					CreatedAt:         testProcedure.CreatedAt,
					UpdatedAt:         *testProcedure.UpdatedAt,
				}
				svc.On("GetProcedureByID", mock.Anything, validID).
					Return(testResponseDto, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := proceduremocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := utils.SetupTestRouter(http.MethodGet, "/procedures/:procedureId",
				procedure.GetProcedureByIDHandler(mockSvc))

			req := httptest.NewRequest(http.MethodGet,
				"/procedures/"+tc.procedureID,
				nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandler_CreateProcedure(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		setupMocks     func(svc *proceduremocks.Service)
		expectedStatus int
	}{
		{
			name:        "success_created",
			requestBody: `{"title":"Test Procedure","description":"Test Description","steps":[{"step_order":1,"title":"Step 1","description":"Desc 1","type":"action"}],"assigned_experiments":[]}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*procedure.CreateProcedureDto")).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "invalid_request_body",
			requestBody: `{"title":""}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "CreateProcedure")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:        "forbidden_create",
			requestBody: `{"title":"Test Procedure","description":"Test Description","steps":[{"step_order":1,"title":"Step 1","description":"Desc 1","type":"action"}],"assigned_experiments":[]}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*procedure.CreateProcedureDto")).
					Return(constant.ErrForbidCreateProcedure)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "procedure_already_exists",
			requestBody: `{"title":"Test Procedure","description":"Test Description","steps":[{"step_order":1,"title":"Step 1","description":"Desc 1","type":"action"}],"assigned_experiments":[]}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*procedure.CreateProcedureDto")).
					Return(constant.ErrProcedureAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "internal_server_error",
			requestBody: `{"title":"Test Procedure","description":"Test Description","steps":[{"step_order":1,"title":"Step 1","description":"Desc 1","type":"action"}],"assigned_experiments":[]}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*procedure.CreateProcedureDto")).
					Return(errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := proceduremocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := utils.SetupTestRouter(http.MethodPost, "/procedures",
				procedure.CreateProcedureHandler(mockSvc))

			req := httptest.NewRequest(http.MethodPost,
				"/procedures",
				strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
