package integration

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"nh-be/internal/constant"
	"nh-be/internal/features/experiment"
	"nh-be/internal/features/experiment/result"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/internal/utils/timeutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestContext struct {
	DB *gorm.DB
	Tx *gorm.DB
}

func SetupTestDB(ctx context.Context) (*TestContext, error) {
	// Try multiple paths for .env file based on test location
	envPaths := []string{
		"../../.env",    // from tests/integration/
		"../../../.env", // from tests/integration/experiment/
		".env",          // from project experiment
	}
	loaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			loaded = true
			break
		}
	}
	if !loaded {
		log.Printf("Warning: Could not load .env file from any path")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

	if host == "" || port == "" || user == "" || dbname == "" {
		return nil, fmt.Errorf("missing required database environment variables (DB_HOST, DB_PORT, DB_USERNAME, DB_NAME)")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		user, pass, host, port, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(
		&permission.Permission{},
		&permission.PermissionGroup{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	return &TestContext{
		DB: db,
		Tx: tx,
	}, nil
}

func (tc *TestContext) Teardown(ctx context.Context) {
	if tc.Tx != nil {
		tc.Tx.Rollback()
	}
}

func CreateTestUserWithoutPermission(ctx context.Context, tx *gorm.DB) (*user.User, error) {
	testUser := &user.User{
		ID:       uuid.New(),
		Username: fmt.Sprintf("noperm_%s", uuid.New().String()[:8]),
		Email:    fmt.Sprintf("noperm_%s@example.com", uuid.New().String()[:8]),
		Password: "hashedpassword",
	}
	if err := tx.Create(testUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return testUser, nil
}

func CreateTestUser(ctx context.Context, tx *gorm.DB) (*user.User, error) {
	perm := &permission.Permission{}
	err := tx.Where("code_name = ?", constant.ManageExperiment).First(perm).Error
	if err != nil {
		perm = &permission.Permission{
			ID:       uuid.New(),
			Name:     "Manage Experiment Test",
			CodeName: constant.ManageExperiment,
		}
		if err := tx.Create(perm).Error; err != nil {
			return nil, fmt.Errorf("failed to create permission: %w", err)
		}
	}

	permGroup := &permission.PermissionGroup{
		ID:          uuid.New(),
		Name:        fmt.Sprintf("Test Group %s", uuid.New().String()[:8]),
		Permissions: []permission.Permission{*perm},
	}
	if err := tx.Create(permGroup).Error; err != nil {
		return nil, fmt.Errorf("failed to create permission group: %w", err)
	}

	testUser := &user.User{
		ID:                       uuid.New(),
		Username:                 fmt.Sprintf("testuser_%s", uuid.New().String()[:8]),
		Email:                    fmt.Sprintf("test_%s@example.com", uuid.New().String()[:8]),
		Password:                 "hashedpassword",
		AssignedPermissionGroups: []permission.PermissionGroup{*permGroup},
	}
	if err := tx.Create(testUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return testUser, nil
}

func CreateTestExperiment(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (*experiment.Experiment, error) {
	exp := &experiment.Experiment{
		ID:          uuid.New(),
		Title:       "Test Experiment",
		Objective:   "Test objective for integration testing",
		Status:      experiment.ExperimentDraft,
		Type:        experiment.ExperimentExploratoryType,
		Version:     1,
		CreatedByID: userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   timeutil.TimePtr(time.Now()),
	}
	if err := tx.Create(exp).Error; err != nil {
		return nil, fmt.Errorf("failed to create experiment: %w", err)
	}
	return exp, nil
}

func SetupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")

	protected := api.Group("")
	protected.Use(contextAuthMiddleware())

	experimentsGroup := protected.Group("/experiments")

	resultRepo := result.NewRepository(db)
	permissionRepo := permission.NewRepository(db)
	permissionService := permission.NewService(permissionRepo, permission.NewNoOpPermissionCache())
	resultService := result.NewService(resultRepo, permissionService)

	experimentsGroup.POST("/:experimentId/result", result.CreateResultHandler(resultService))
	experimentsGroup.GET("/:experimentId/result", result.GetResultByExperimentIDHandler(resultService))
	experimentsGroup.PUT("/:experimentId/result/:resultId", result.UpdateResultHandler(resultService))

	return r
}

func contextAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID := c.Request.Context().Value(constant.CtxUserId); userID != nil {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
			"message": "Unauthorized",
		})
	}
}

func MakeAuthenticatedRequest(router *gin.Engine, method, path, body string, userID uuid.UUID) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	ctx := context.WithValue(req.Context(), constant.CtxUserId, userID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
