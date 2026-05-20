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
