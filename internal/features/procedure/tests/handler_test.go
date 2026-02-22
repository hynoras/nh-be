package procedure

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"nh-be/internal/constant"
	"nh-be/internal/features/procedure"
	proceduremocks "nh-be/internal/features/procedure/mocks"
	"nh-be/internal/utils/testutil"
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

			router := testutil.SetupTestRouter(http.MethodGet, "/procedures",
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
					Return(nil, constant.ErrProcedureConflict)
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
					ID:          validID.String(),
					Title:       "Test Procedure",
					Description: "Test Procedure Description",
					CreatedAt:   testProcedure.CreatedAt,
					UpdatedAt:   *testProcedure.UpdatedAt,
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

			router := testutil.SetupTestRouter(http.MethodGet, "/procedures/:procedureId",
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
			requestBody: `{"title":"Test Procedure","description":"Test Description","steps":[{"step_order":1,"title":"Step 1","description":"Description Step 1","is_optional":false,"type":"action"}]}`,
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
			requestBody: `{"title":"Test Procedure","description":"Test Description","steps":[{"step_order":1,"title":"Step 1","description":null,"is_optional":false,"wait_time":null,"type":"action"}]}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*procedure.CreateProcedureDto")).
					Return(constant.ErrForbidCreateProcedure)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "procedure_already_exists",
			requestBody: `{"title":"Test Procedure","description":"Test Description","steps":[{"step_order":1,"title":"Step 1","description":null,"is_optional":false,"wait_time":null,"type":"action"}]}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("CreateProcedure", mock.Anything, mock.AnythingOfType("*procedure.CreateProcedureDto")).
					Return(constant.ErrProcedureAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "internal_server_error",
			requestBody: `{"title":"Test Procedure","description":"Test Description","steps":[{"step_order":1,"title":"Step 1","description":null,"is_optional":false,"wait_time":null,"type":"action"}]}`,
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

			router := testutil.SetupTestRouter(http.MethodPost, "/procedures",
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

func TestHandler_UpdateProcedure(t *testing.T) {
	validID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := []struct {
		name           string
		procedureID    string
		requestBody    string
		setupMocks     func(svc *proceduremocks.Service)
		expectedStatus int
	}{
		{
			name:        "invalid_uuid",
			procedureID: "invalid",
			requestBody: `{"title":"Updated Title","description":"Updated Description","version":1}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "UpdateProcedure")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid_body",
			procedureID: validID.String(),
			requestBody: `{"title":""}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "UpdateProcedure")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:        "update_success",
			procedureID: validID.String(),
			requestBody: `{"title":"Updated Title","description":"Updated Description","version":1}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedure", mock.Anything, validID, mock.AnythingOfType("*procedure.UpdateProcedureDto")).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "forbidden_update",
			procedureID: validID.String(),
			requestBody: `{"title":"Updated Title","description":"Updated Description","version":1}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedure", mock.Anything, validID, mock.AnythingOfType("*procedure.UpdateProcedureDto")).
					Return(constant.ErrForbidUpdateProcedure)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "procedure_not_found",
			procedureID: validID.String(),
			requestBody: `{"title":"Updated Title","description":"Updated Description","version":1}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedure", mock.Anything, validID, mock.AnythingOfType("*procedure.UpdateProcedureDto")).
					Return(constant.ErrProcedureNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "version_conflict",
			procedureID: validID.String(),
			requestBody: `{"title":"Updated Title","description":"Updated Description","version":1}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedure", mock.Anything, validID, mock.AnythingOfType("*procedure.UpdateProcedureDto")).
					Return(constant.ErrProcedureConflict)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "unexpected_error",
			procedureID: validID.String(),
			requestBody: `{"title":"Updated Title","description":"Updated Description","version":1}`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedure", mock.Anything, validID, mock.AnythingOfType("*procedure.UpdateProcedureDto")).
					Return(errors.New("unexpected database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := proceduremocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := testutil.SetupTestRouter(http.MethodPut, "/procedures/:procedureId",
				procedure.UpdateProcedureHandler(mockSvc))

			req := httptest.NewRequest(http.MethodPut,
				"/procedures/"+tc.procedureID,
				strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandler_UpdateProcedureStep(t *testing.T) {
	validID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	// Valid UUID v4 for use in step bodies
	validStepID := "550e8400-e29b-41d4-a716-446655440000"
	validStepBody := fmt.Sprintf(`[{"id":"%s","title":"Step 1","type":"action","version":1}]`, validStepID)

	tests := []struct {
		name           string
		procedureID    string
		requestBody    string
		setupMocks     func(svc *proceduremocks.Service)
		expectedStatus int
	}{
		{
			name:        "invalid_procedure_id",
			procedureID: "invalid",
			requestBody: `[]`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "UpdateProcedureStep")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid_json_body",
			procedureID: validID.String(),
			requestBody: `invalid-json`,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "UpdateProcedureStep")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "validation_error",
			procedureID: validID.String(),
			requestBody: fmt.Sprintf(`[{"id":"%s","title":"Valid Title","type":"not_valid_type"}]`, validStepID),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "UpdateProcedureStep")
			},
			expectedStatus: http.StatusBadRequest, // gin returns 400 for slice binding validation
		},
		{
			name:        "permission_denied",
			procedureID: validID.String(),
			requestBody: validStepBody,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedureStep", mock.Anything, validID, mock.AnythingOfType("[]procedure.UpdateProcedureStepInput")).
					Return(constant.ErrForbidUpdateProcedure)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "procedure_not_found",
			procedureID: validID.String(),
			requestBody: validStepBody,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedureStep", mock.Anything, validID, mock.AnythingOfType("[]procedure.UpdateProcedureStepInput")).
					Return(constant.ErrProcedureNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "procedure_step_conflict",
			procedureID: validID.String(),
			requestBody: validStepBody,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedureStep", mock.Anything, validID, mock.AnythingOfType("[]procedure.UpdateProcedureStepInput")).
					Return(constant.ErrProcedureConflict)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "procedure_step_not_found",
			procedureID: validID.String(),
			requestBody: validStepBody,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedureStep", mock.Anything, validID, mock.AnythingOfType("[]procedure.UpdateProcedureStepInput")).
					Return(constant.ErrProcedureStepNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "unexpected_error",
			procedureID: validID.String(),
			requestBody: validStepBody,
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedureStep", mock.Anything, validID, mock.AnythingOfType("[]procedure.UpdateProcedureStepInput")).
					Return(errors.New("unexpected error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "successful_update",
			procedureID: validID.String(),
			requestBody: fmt.Sprintf(`[{"id":"%s","title":"Updated Step 1","type":"action","version":1}]`, validStepID),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("UpdateProcedureStep", mock.Anything, validID, mock.AnythingOfType("[]procedure.UpdateProcedureStepInput")).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := proceduremocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := testutil.SetupTestRouter(http.MethodPut, "/procedures/:procedureId/procedure-steps",
				procedure.UpdateProcedureStepHandler(mockSvc))

			req := httptest.NewRequest(http.MethodPut,
				"/procedures/"+tc.procedureID+"/procedure-steps",
				strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandler_DeleteProcedure(t *testing.T) {
	validID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := []struct {
		name           string
		procedureID    string
		setupMocks     func(svc *proceduremocks.Service)
		expectedStatus int
	}{
		{
			name:        "success",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("DeleteProcedure", mock.Anything, validID).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid_uuid",
			procedureID: "invalid-uuid",
			setupMocks: func(svc *proceduremocks.Service) {
				svc.AssertNotCalled(t, "DeleteProcedure", mock.Anything, mock.Anything)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "service_forbidden",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("DeleteProcedure", mock.Anything, validID).
					Return(constant.ErrForbidDeleteProcedure)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:        "service_not_found",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("DeleteProcedure", mock.Anything, validID).
					Return(constant.ErrProcedureNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "service_internal_error",
			procedureID: validID.String(),
			setupMocks: func(svc *proceduremocks.Service) {
				svc.On("DeleteProcedure", mock.Anything, validID).
					Return(errors.New("generic error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := proceduremocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := testutil.SetupTestRouter(http.MethodDelete, "/procedures/:procedureId",
				procedure.DeleteProcedureHandler(mockSvc))

			req := httptest.NewRequest(http.MethodDelete,
				"/procedures/"+tc.procedureID,
				nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
