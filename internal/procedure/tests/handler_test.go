package procedure

import (
	"errors"
	"net/http"
	"net/http/httptest"
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
