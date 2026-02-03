package procedure

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
