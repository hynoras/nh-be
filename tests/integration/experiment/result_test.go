package experiment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"nh-be/internal/features/experiment/result"
	setup "nh-be/tests/integration"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetResultByExperimentID(t *testing.T) {
	tests := []struct {
		name           string
		setupTestData  func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func())
		expectedStatus int
		checkResponse  func(t *testing.T, resp map[string]interface{})
		customPath     string
		useNoPermUser  bool
	}{
		{
			name: "happy_path",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				res := &result.ExperimentResult{
					ID:              uuid.New(),
					ExperimentID:    exp.ID,
					Outcome:         result.OutcomeSuccess,
					Summary:         "This is a test summary for the experiment result",
					OutcomeReason:   "Reason for logic",
					ConfidenceLevel: result.ConfidenceHigh,
					Version:         1,
				}
				db.Create(res)
				return &exp.ID, func() {
					db.Delete(res)
					db.Delete(exp)
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "Experiment result fetched successfully", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "success", data["outcome"])
				assert.Equal(t, "high", data["confidence_level"])
				assert.Equal(t, "This is a test summary for the experiment result", data["summary"])
				assert.Equal(t, "Reason for logic", data["outcome_reason"])
				assert.EqualValues(t, 1, data["version"])
			},
		},
		{
			name: "result_does_not_exist",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				return &exp.ID, func() {
					db.Delete(exp)
				}
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Experiment result not found", resp["message"])
			},
		},
		{
			name: "experiment_does_not_exist",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				id := uuid.New()
				return &id, func() {}
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				// Based on current implementation, error is ErrExperimentResultNotFound when First fails
				assert.Equal(t, "Experiment result not found", resp["message"])
			},
		},
		{
			name: "invalid_experiment_uuid",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				return nil, func() {}
			},
			customPath:     "/api/v1/experiments/invalid-uuid/result",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Invalid ID format", resp["message"])
			},
		},
		{
			name: "permission_denied",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				return &exp.ID, func() {
					db.Delete(exp)
				}
			},
			useNoPermUser:  true,
			expectedStatus: http.StatusForbidden,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Authorization failed", resp["message"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			testCtx, err := setup.SetupTestDB(ctx)
			require.NoError(t, err, "Failed to setup test database")
			defer testCtx.Teardown(ctx)

			testUser, err := setup.CreateTestUser(ctx, testCtx.DB)
			require.NoError(t, err, "Failed to create test user")

			requestUserID := testUser.ID
			if tc.useNoPermUser {
				noPermUser, err := setup.CreateTestUserWithoutPermission(ctx, testCtx.DB)
				require.NoError(t, err)
				requestUserID = noPermUser.ID
			}

			var experimentID *uuid.UUID
			var cleanup func()
			if tc.setupTestData != nil {
				experimentID, cleanup = tc.setupTestData(ctx, testCtx.DB, testUser.ID)
				defer cleanup()
			}

			router := setup.SetupTestRouter(testCtx.DB)
			path := tc.customPath
			if path == "" && experimentID != nil {
				path = fmt.Sprintf("/api/v1/experiments/%s/result", experimentID.String())
			}

			resp := setup.MakeAuthenticatedRequest(router, http.MethodGet, path, "", requestUserID)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.checkResponse != nil {
				var response map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &response)
				require.NoError(t, err)
				tc.checkResponse(t, response)
			}
		})
	}
}

func TestCreateResult(t *testing.T) {
	tests := []struct {
		name           string
		setupTestData  func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func())
		requestBody    string
		expectedStatus int
		checkResponse  func(t *testing.T, resp map[string]interface{})
		checkDB        func(t *testing.T, db *gorm.DB, experimentID uuid.UUID)
		customPath     string
		useNoPermUser  bool
	}{
		{
			name: "happy_path",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				return &exp.ID, func() {
					db.Where("experiment_id = ?", exp.ID).Delete(&result.ExperimentResult{})
					db.Delete(exp)
				}
			},
			requestBody: `{
				"outcome": "success",
				"summary": "This is a test summary for the experiment result",
				"outcome_reason": "The experiment achieved its objectives successfully",
				"confidence_level": "high"
			}`,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "Experiment result created successfully", resp["message"])
			},
			checkDB: func(t *testing.T, db *gorm.DB, experimentID uuid.UUID) {
				var createdResult result.ExperimentResult
				err := db.Where("experiment_id = ?", experimentID).First(&createdResult).Error
				require.NoError(t, err)
				assert.Equal(t, experimentID, createdResult.ExperimentID)
				assert.Equal(t, result.OutcomeSuccess, createdResult.Outcome)
				assert.Equal(t, result.ConfidenceHigh, createdResult.ConfidenceLevel)
				assert.Equal(t, 1, createdResult.Version)
			},
		},
		{
			name: "experiment_does_not_exist",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				nonExistentID := uuid.New()
				return &nonExistentID, func() {}
			},
			requestBody: `{
				"outcome": "success",
				"summary": "This is a test summary for the experiment result",
				"outcome_reason": "The experiment achieved its objectives successfully",
				"confidence_level": "high"
			}`,
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Experiment not found", resp["message"])
			},
			checkDB: nil,
		},
		{
			name: "experiment_result_already_exists",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				existingResult := &result.ExperimentResult{
					ID:              uuid.New(),
					ExperimentID:    exp.ID,
					Outcome:         result.OutcomeSuccess,
					Summary:         "Existing result",
					OutcomeReason:   "Existing reason",
					ConfidenceLevel: result.ConfidenceHigh,
					Version:         1,
				}
				db.Create(existingResult)
				return &exp.ID, func() {
					db.Where("experiment_id = ?", exp.ID).Delete(&result.ExperimentResult{})
					db.Delete(exp)
				}
			},
			requestBody: `{
				"outcome": "success",
				"summary": "This is a test summary for the experiment result",
				"outcome_reason": "The experiment achieved its objectives successfully",
				"confidence_level": "high"
			}`,
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Experiment result already exists", resp["message"])
			},
			checkDB: func(t *testing.T, db *gorm.DB, experimentID uuid.UUID) {
				var count int64
				db.Model(&result.ExperimentResult{}).Where("experiment_id = ?", experimentID).Count(&count)
				assert.Equal(t, int64(1), count)
			},
		},
		{
			name: "invalid_enum_values",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				return &exp.ID, func() {
					db.Delete(exp)
				}
			},
			requestBody: `{
				"outcome": "invalid_outcome",
				"summary": "This is a valid summary text longer than 10 chars",
				"outcome_reason": "This is a valid reason text longer than 10 chars",
				"confidence_level": "invalid_level"
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Validation failed", resp["message"])
			},
		},
		{
			name: "missing_required_fields",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				return &exp.ID, func() {
					db.Delete(exp)
				}
			},
			requestBody:    `{}`,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Validation failed", resp["message"])
			},
		},
		{
			name: "user_not_allowed_to_create",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				return &exp.ID, func() {
					db.Delete(exp)
				}
			},
			requestBody: `{
				"outcome": "success",
				"summary": "This is a valid summary text longer than 10 chars",
				"outcome_reason": "This is a valid reason text longer than 10 chars",
				"confidence_level": "high"
			}`,
			expectedStatus: http.StatusForbidden,
			useNoPermUser:  true,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Authorization failed", resp["message"])
			},
		},
		{
			name: "invalid_experiment_uuid",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*uuid.UUID, func()) {
				return nil, func() {}
			},
			customPath:     "/api/v1/experiments/invalid-uuid/result",
			requestBody:    `{}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Invalid ID format", resp["message"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			testCtx, err := setup.SetupTestDB(ctx)
			require.NoError(t, err, "Failed to setup test database")
			defer testCtx.Teardown(ctx)

			testUser, err := setup.CreateTestUser(ctx, testCtx.DB)
			require.NoError(t, err, "Failed to create test user")

			requestUserID := testUser.ID
			if tc.useNoPermUser {
				noPermUser, err := setup.CreateTestUserWithoutPermission(ctx, testCtx.DB)
				require.NoError(t, err)
				requestUserID = noPermUser.ID
			}

			var experimentID *uuid.UUID
			var cleanup func()
			if tc.setupTestData != nil {
				experimentID, cleanup = tc.setupTestData(ctx, testCtx.DB, testUser.ID)
				defer cleanup()
			}

			router := setup.SetupTestRouter(testCtx.DB)
			path := tc.customPath
			if path == "" && experimentID != nil {
				path = fmt.Sprintf("/api/v1/experiments/%s/result", experimentID.String())
			}

			resp := setup.MakeAuthenticatedRequest(router, http.MethodPost, path, tc.requestBody, requestUserID)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.checkResponse != nil {
				var response map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &response)
				require.NoError(t, err)
				tc.checkResponse(t, response)
			}

			if tc.checkDB != nil && experimentID != nil {
				tc.checkDB(t, testCtx.DB, *experimentID)
			}
		})
	}
}

func TestCreateResult_Idempotency(t *testing.T) {
	tests := []struct {
		name            string
		numRequests     int
		expected201     int
		expected409     int
		expectedDBCount int64
	}{
		{
			name:            "concurrent_requests",
			numRequests:     2,
			expected201:     1,
			expected409:     1,
			expectedDBCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			testCtx, err := setup.SetupTestDB(ctx)
			require.NoError(t, err, "Failed to setup test database")

			testUser, err := setup.CreateTestUser(ctx, testCtx.DB)
			require.NoError(t, err, "Failed to create test user")

			exp, err := setup.CreateTestExperiment(ctx, testCtx.DB, testUser.ID)
			require.NoError(t, err, "Failed to create test experiment")

			defer func() {
				testCtx.DB.Where("experiment_id = ?", exp.ID).Delete(&result.ExperimentResult{})
				testCtx.DB.Delete(exp)
				testCtx.DB.Delete(testUser)
				testCtx.Teardown(ctx)
			}()

			router := setup.SetupTestRouter(testCtx.DB)

			requestBody := `{
				"outcome": "success",
				"summary": "This is a test summary for the experiment result",
				"outcome_reason": "The experiment achieved its objectives successfully",
				"confidence_level": "high"
			}`
			path := fmt.Sprintf("/api/v1/experiments/%s/result", exp.ID.String())

			var wg sync.WaitGroup
			responseCodes := make([]int, tc.numRequests)

			for i := 0; i < tc.numRequests; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					resp := setup.MakeAuthenticatedRequest(router, http.MethodPost, path, requestBody, testUser.ID)
					responseCodes[idx] = resp.Code
				}(i)
			}

			wg.Wait()

			count201, count409 := 0, 0
			for _, code := range responseCodes {
				switch code {
				case http.StatusCreated:
					count201++
				case http.StatusConflict:
					count409++
				}
			}

			assert.Equal(t, tc.expected201, count201, "Expected %d 201 responses", tc.expected201)
			assert.Equal(t, tc.expected409, count409, "Expected %d 409 responses", tc.expected409)

			var count int64
			err = testCtx.DB.Model(&result.ExperimentResult{}).Where("experiment_id = ?", exp.ID).Count(&count).Error
			require.NoError(t, err)
			assert.Equal(t, tc.expectedDBCount, count, "Expected %d result(s) in database", tc.expectedDBCount)
		})
	}
}

func TestUpdateResult(t *testing.T) {
	tests := []struct {
		name           string
		setupTestData  func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (resultID *uuid.UUID, experimentID *uuid.UUID, cleanup func())
		requestBody    string
		expectedStatus int
		checkResponse  func(t *testing.T, resp map[string]interface{})
		checkDB        func(t *testing.T, db *gorm.DB, resultID uuid.UUID, experimentID uuid.UUID)
		customPath     func(resultID *uuid.UUID, experimentID *uuid.UUID) string
		useNoPermUser  bool
	}{
		{
			name: "happy_path",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (resultID *uuid.UUID, experimentID *uuid.UUID, cleanup func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				res := &result.ExperimentResult{
					ID:              uuid.New(),
					ExperimentID:    exp.ID,
					Outcome:         result.OutcomeSuccess,
					Summary:         "Initial summary that is long enough",
					OutcomeReason:   "Initial reason that is long enough",
					ConfidenceLevel: result.ConfidenceHigh,
					Version:         1,
				}
				db.Create(res)
				return &res.ID, &exp.ID, func() {
					db.Delete(res)
					db.Delete(exp)
				}
			},
			requestBody: `{
				"version": 1,
				"outcome": "failure",
				"summary": "Updated summary that is also long enough",
				"outcome_reason": "Updated reason that is also long enough",
				"confidence_level": "low"
			}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "Experiment result updated successfully", resp["message"])
			},
			checkDB: func(t *testing.T, db *gorm.DB, resultID uuid.UUID, experimentID uuid.UUID) {
				var updatedResult result.ExperimentResult
				err := db.Where("id = ? AND experiment_id = ?", resultID, experimentID).First(&updatedResult).Error
				require.NoError(t, err)
				assert.Equal(t, result.OutcomeFailure, updatedResult.Outcome)
				assert.Equal(t, "Updated summary that is also long enough", updatedResult.Summary)
				assert.Equal(t, result.ConfidenceLow, updatedResult.ConfidenceLevel)
				assert.Equal(t, 2, updatedResult.Version)
			},
		},
		{
			name: "permission_denied",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (resultID *uuid.UUID, experimentID *uuid.UUID, cleanup func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				res := &result.ExperimentResult{
					ID:              uuid.New(),
					ExperimentID:    exp.ID,
					Outcome:         result.OutcomeSuccess,
					Summary:         "Initial summary that is long enough",
					OutcomeReason:   "Initial reason that is long enough",
					ConfidenceLevel: result.ConfidenceHigh,
					Version:         1,
				}
				db.Create(res)
				return &res.ID, &exp.ID, func() {
					db.Delete(res)
					db.Delete(exp)
				}
			},
			useNoPermUser: true,
			requestBody: `{
				"version": 1,
				"outcome": "failure"
			}`,
			expectedStatus: http.StatusForbidden,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Authorization failed", resp["message"])
			},
		},
		{
			name: "invalid_body_validation_failure",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (resultID *uuid.UUID, experimentID *uuid.UUID, cleanup func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				res := &result.ExperimentResult{
					ID:              uuid.New(),
					ExperimentID:    exp.ID,
					Outcome:         result.OutcomeSuccess,
					Summary:         "Initial summary that is long enough",
					OutcomeReason:   "Initial reason that is long enough",
					ConfidenceLevel: result.ConfidenceHigh,
					Version:         1,
				}
				db.Create(res)
				return &res.ID, &exp.ID, func() {
					db.Delete(res)
					db.Delete(exp)
				}
			},
			requestBody: `{
				"version": 1,
				"outcome": "invalid_outcome",
				"summary": "short"
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Validation failed", resp["message"])
			},
		},
		{
			name: "result_not_found",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (resultID *uuid.UUID, experimentID *uuid.UUID, cleanup func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				nonExistentResultID := uuid.New()
				return &nonExistentResultID, &exp.ID, func() {
					db.Delete(exp)
				}
			},
			requestBody: `{
				"version": 1,
				"outcome": "success"
			}`,
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Experiment result not found", resp["message"])
			},
		},
		{
			name: "experiment_not_found",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (resultID *uuid.UUID, experimentID *uuid.UUID, cleanup func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				res := &result.ExperimentResult{
					ID:              uuid.New(),
					ExperimentID:    exp.ID,
					Outcome:         result.OutcomeSuccess,
					Summary:         "Initial summary that is long enough",
					OutcomeReason:   "Initial reason that is long enough",
					ConfidenceLevel: result.ConfidenceHigh,
					Version:         1,
				}
				db.Create(res)
				nonExistentExperimentID := uuid.New()
				return &res.ID, &nonExistentExperimentID, func() {
					db.Delete(res)
					db.Delete(exp)
				}
			},
			requestBody: `{
				"version": 1,
				"outcome": "success"
			}`,
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Experiment result not found", resp["message"])
			},
		},
		{
			name: "invalid_experiment_uuid",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (resultID *uuid.UUID, experimentID *uuid.UUID, cleanup func()) {
				return nil, nil, func() {}
			},
			customPath: func(resultID *uuid.UUID, experimentID *uuid.UUID) string {
				return fmt.Sprintf("/api/v1/experiments/invalid-uuid/result/%s", uuid.New().String())
			},
			requestBody:    `{"version": 1}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Invalid ID format", resp["message"])
			},
		},
		{
			name: "invalid_result_uuid",
			setupTestData: func(ctx context.Context, db *gorm.DB, userID uuid.UUID) (resultID *uuid.UUID, experimentID *uuid.UUID, cleanup func()) {
				exp, _ := setup.CreateTestExperiment(ctx, db, userID)
				return nil, &exp.ID, func() {
					db.Delete(exp)
				}
			},
			customPath: func(resultID *uuid.UUID, experimentID *uuid.UUID) string {
				return fmt.Sprintf("/api/v1/experiments/%s/result/invalid-uuid", experimentID.String())
			},
			requestBody:    `{"version": 1}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				assert.Equal(t, "Invalid ID format", resp["message"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			testCtx, err := setup.SetupTestDB(ctx)
			require.NoError(t, err, "Failed to setup test database")
			defer testCtx.Teardown(ctx)

			testUser, err := setup.CreateTestUser(ctx, testCtx.DB)
			require.NoError(t, err, "Failed to create test user")

			requestUserID := testUser.ID
			if tc.useNoPermUser {
				noPermUser, err := setup.CreateTestUserWithoutPermission(ctx, testCtx.DB)
				require.NoError(t, err)
				requestUserID = noPermUser.ID
			}

			var resultID, experimentID *uuid.UUID
			var cleanup func()
			if tc.setupTestData != nil {
				resultID, experimentID, cleanup = tc.setupTestData(ctx, testCtx.DB, testUser.ID)
				defer cleanup()
			}

			router := setup.SetupTestRouter(testCtx.DB)
			path := ""
			if tc.customPath != nil {
				path = tc.customPath(resultID, experimentID)
			} else if resultID != nil && experimentID != nil {
				path = fmt.Sprintf("/api/v1/experiments/%s/result/%s", experimentID.String(), resultID.String())
			}

			resp := setup.MakeAuthenticatedRequest(router, http.MethodPut, path, tc.requestBody, requestUserID)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.checkResponse != nil {
				var response map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &response)
				require.NoError(t, err)
				tc.checkResponse(t, response)
			}

			if tc.checkDB != nil && resultID != nil && experimentID != nil {
				tc.checkDB(t, testCtx.DB, *resultID, *experimentID)
			}
		})
	}
}

func TestUpdateResult_OptimisticLocking(t *testing.T) {
	ctx := context.Background()

	testCtx, err := setup.SetupTestDB(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer testCtx.Teardown(ctx)

	testUser, err := setup.CreateTestUser(ctx, testCtx.DB)
	require.NoError(t, err, "Failed to create test user")

	exp, _ := setup.CreateTestExperiment(ctx, testCtx.DB, testUser.ID)
	res := &result.ExperimentResult{
		ID:              uuid.New(),
		ExperimentID:    exp.ID,
		Outcome:         result.OutcomeSuccess,
		Summary:         "Initial summary that is long enough",
		OutcomeReason:   "Initial reason that is long enough",
		ConfidenceLevel: result.ConfidenceHigh,
		Version:         1,
	}
	testCtx.DB.Create(res)

	defer func() {
		testCtx.DB.Delete(res)
		testCtx.DB.Delete(exp)
	}()

	router := setup.SetupTestRouter(testCtx.DB)
	path := fmt.Sprintf("/api/v1/experiments/%s/result/%s", exp.ID.String(), res.ID.String())

	requestBody1 := `{
		"version": 1,
		"summary": "First caller updated summary long enough"
	}`
	requestBody2 := `{
		"version": 1,
		"summary": "Second caller updated summary long enough"
	}`

	var wg sync.WaitGroup
	wg.Add(2)

	codes := make([]int, 2)
	go func() {
		defer wg.Done()
		resp := setup.MakeAuthenticatedRequest(router, http.MethodPut, path, requestBody1, testUser.ID)
		codes[0] = resp.Code
	}()

	go func() {
		defer wg.Done()
		resp := setup.MakeAuthenticatedRequest(router, http.MethodPut, path, requestBody2, testUser.ID)
		codes[1] = resp.Code
	}()

	wg.Wait()

	// One should succeed (200 OK), one should fail (409 Conflict)
	successCount := 0
	conflictCount := 0
	for _, code := range codes {
		switch code {
		case http.StatusOK:
			successCount++
		case http.StatusConflict:
			conflictCount++
		}
	}

	assert.Equal(t, 1, successCount, "Expected exactly one success")
	assert.Equal(t, 1, conflictCount, "Expected exactly one conflict")

	// Final version in DB should be 2
	var finalRes result.ExperimentResult
	testCtx.DB.Where("id = ?", res.ID).First(&finalRes)
	assert.Equal(t, 2, finalRes.Version)
}
