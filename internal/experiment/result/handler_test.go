package result_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	resHandler "nh-be/internal/experiment/result"
	resultmocks "nh-be/internal/experiment/result/mocks"
)

// setupTestRouter creates a test router with the given handler at the specified route
func setupTestRouter(method, path string, handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(method, path, handler)
	return router
}

func TestHandler_GetResultByExperimentID(t *testing.T) {
	tests := []struct {
		name           string
		experimentID   string
		setupMocks     func(svc *resultmocks.Service)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:         "success",
			experimentID: "00000000-0000-0000-0000-000000000000",
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("GetResultByExperimentID", mock.Anything, mock.Anything).
					Return(createTestResult(uuid.MustParse("00000000-0000-0000-0000-000000000000")), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "invalid_uuid",
			experimentID: "not-a-uuid",
			setupMocks: func(svc *resultmocks.Service) {
				svc.AssertNotCalled(t, "GetResultByExperimentID")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "experiment_not_found",
			experimentID: "00000000-0000-0000-0000-000000000000",
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("GetResultByExperimentID", mock.Anything, mock.Anything).
					Return(nil, resHandler.ErrExperimentNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:         "forbidden_no_permission",
			experimentID: "00000000-0000-0000-0000-000000000000",
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("GetResultByExperimentID", mock.Anything, mock.Anything).
					Return(nil, resHandler.ErrForbidViewExperimentResult)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "internal_server_error",
			experimentID: "00000000-0000-0000-0000-000000000000",
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("GetResultByExperimentID", mock.Anything, mock.Anything).
					Return(nil, errors.New("unexpected database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := resultmocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := setupTestRouter(http.MethodGet, "/experiments/:experimentId/result",
				resHandler.GetResultByExperimentIDHandler(mockSvc))

			req := httptest.NewRequest(http.MethodGet,
				"/experiments/"+tc.experimentID+"/result",
				nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandler_CreateResult(t *testing.T) {
	tests := []struct {
		name           string
		experimentID   string
		requestBody    interface{}
		setupMocks     func(svc *resultmocks.Service)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:         "success",
			experimentID: "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("CreateResult", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:         "invalid_uuid",
			experimentID: "not-a-uuid",
			requestBody:  createTestDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.AssertNotCalled(t, "CreateResult")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "invalid_format",
			experimentID: uuid.New().String(),
			requestBody:  createInvalidTestDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.AssertNotCalled(t, "CreateResult")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:         "experiment_not_found",
			experimentID: "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("CreateResult", mock.Anything, mock.Anything, mock.Anything).
					Return(resHandler.ErrExperimentNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:         "experiment_result_already_exists",
			experimentID: "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("CreateResult", mock.Anything, mock.Anything, mock.Anything).
					Return(resHandler.ErrExperimentResultAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:         "forbidden_no_permission",
			experimentID: "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("CreateResult", mock.Anything, mock.Anything, mock.Anything).
					Return(resHandler.ErrForbidCreateExperimentResult)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "internal_server_error",
			experimentID: "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("CreateResult", mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("unexpected database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := resultmocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := setupTestRouter(http.MethodPost, "/experiments/:experimentId/result",
				resHandler.CreateResultHandler(mockSvc))

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost,
				"/experiments/"+tc.experimentID+"/result",
				bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandler_UpdateResult(t *testing.T) {
	tests := []struct {
		name           string
		experimentID   string
		resultID       string
		requestBody    interface{}
		setupMocks     func(svc *resultmocks.Service)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:         "success",
			experimentID: "00000000-0000-0000-0000-000000000000",
			resultID:     "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestUpdateDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("UpdateResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "invalid_uuid",
			experimentID: "not-a-uuid",
			resultID:     "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestUpdateDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.AssertNotCalled(t, "UpdateResult")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "invalid_format",
			experimentID: "00000000-0000-0000-0000-000000000000",
			resultID:     "00000000-0000-0000-0000-000000000000",
			requestBody:  createInvalidTestDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.AssertNotCalled(t, "UpdateResult")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:         "experiment_not_found",
			experimentID: "00000000-0000-0000-0000-000000000000",
			resultID:     "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestUpdateDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("UpdateResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(resHandler.ErrExperimentNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:         "result_not_found",
			experimentID: "00000000-0000-0000-0000-000000000000",
			resultID:     "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestUpdateDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("UpdateResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(resHandler.ErrExperimentResultNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:         "forbidden_no_permission",
			experimentID: "00000000-0000-0000-0000-000000000000",
			resultID:     "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestUpdateDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("UpdateResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(resHandler.ErrForbidUpdateExperimentResult)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "optimistic_locking_conflict",
			experimentID: "00000000-0000-0000-0000-000000000000",
			resultID:     "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestUpdateDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("UpdateResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(resHandler.ErrOptimisticLockingConflict)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:         "internal_server_error",
			experimentID: "00000000-0000-0000-0000-000000000000",
			resultID:     "00000000-0000-0000-0000-000000000000",
			requestBody:  createTestUpdateDto(),
			setupMocks: func(svc *resultmocks.Service) {
				svc.On("UpdateResult", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("unexpected database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := resultmocks.NewService(t)
			tc.setupMocks(mockSvc)

			router := setupTestRouter(http.MethodPut, "/experiments/:experimentId/result/:resultId",
				resHandler.UpdateResultHandler(mockSvc))

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPut,
				"/experiments/"+tc.experimentID+"/result/"+tc.resultID,
				bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
