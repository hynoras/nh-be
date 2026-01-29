package result

import (
	"context"
	"errors"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	return gormDB, mock
}

func createValidExperimentResult() *ExperimentResult {
	return &ExperimentResult{
		ID:              uuid.New(),
		ExperimentID:    uuid.New(),
		Outcome:         OutcomeSuccess,
		Summary:         "Test summary for experiment result",
		OutcomeReason:   "Test outcome reason for the experiment",
		ConfidenceLevel: ConfidenceHigh,
		Version:         1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func TestRepository_Create_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	result := createValidExperimentResult()

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "experiment_results" ("id","experiment_id","outcome","summary","outcome_reason","confidence_level","version","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).
		WithArgs(
			result.ID,
			result.ExperimentID,
			result.Outcome,
			result.Summary,
			result.OutcomeReason,
			result.ConfidenceLevel,
			result.Version,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.Create(ctx, result)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRepository_Create_InvalidUUIDFormat(t *testing.T) {
	invalidUUIDs := []string{
		"not-a-uuid",
		"12345",
		"",
		"xyz-abc-def-ghi",
	}

	for _, invalidID := range invalidUUIDs {
		t.Run("invalid_uuid: "+invalidID, func(t *testing.T) {
			_, err := uuid.Parse(invalidID)
			if err == nil {
				t.Errorf("expected UUID parse error for '%s', got nil", invalidID)
			}
		})
	}
}

func TestRepository_Create_FKConstraintViolation(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	result := createValidExperimentResult()
	result.ExperimentID = uuid.New()

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "experiment_results" ("id","experiment_id","outcome","summary","outcome_reason","confidence_level","version","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).
		WithArgs(
			result.ID,
			result.ExperimentID,
			result.Outcome,
			result.Summary,
			result.OutcomeReason,
			result.ConfidenceLevel,
			result.Version,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("pq: insert or update on table \"experiment_results\" violates foreign key constraint \"experiment_results_experiment_id_fkey\""))

	mock.ExpectRollback()

	err := repo.Create(ctx, result)

	if err == nil {
		t.Error("expected FK constraint error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRepository_Create_DuplicateConcurrent(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	experimentID := uuid.New()

	result1 := createValidExperimentResult()
	result1.ExperimentID = experimentID

	result2 := createValidExperimentResult()
	result2.ExperimentID = experimentID

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "experiment_results" ("id","experiment_id","outcome","summary","outcome_reason","confidence_level","version","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).
		WithArgs(
			result1.ID,
			result1.ExperimentID,
			result1.Outcome,
			result1.Summary,
			result1.OutcomeReason,
			result1.ConfidenceLevel,
			result1.Version,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "experiment_results" ("id","experiment_id","outcome","summary","outcome_reason","confidence_level","version","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).
		WithArgs(
			result2.ID,
			result2.ExperimentID,
			result2.Outcome,
			result2.Summary,
			result2.OutcomeReason,
			result2.ConfidenceLevel,
			result2.Version,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"experiment_results_experiment_id_key\""))
	mock.ExpectRollback()

	var wg sync.WaitGroup
	var err1, err2 error

	wg.Add(2)
	go func() {
		defer wg.Done()
		err1 = repo.Create(ctx, result1)
	}()

	time.Sleep(10 * time.Millisecond)

	go func() {
		defer wg.Done()
		err2 = repo.Create(ctx, result2)
	}()

	wg.Wait()

	if err1 != nil {
		t.Errorf("expected first insert to succeed, got: %v", err1)
	}
	if err2 == nil {
		t.Error("expected second insert to fail with unique violation, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// ==================== UPDATE TESTS ====================

func TestRepository_Update_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	resultID := uuid.New()
	experimentID := uuid.New()
	currentVersion := 1

	result := &ExperimentResult{
		Outcome:         OutcomeFailure,
		Summary:         "Updated summary",
		OutcomeReason:   "Updated reason",
		ConfidenceLevel: ConfidenceMedium,
		UpdatedAt:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "experiment_results" SET`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, resultID, experimentID, result, currentVersion)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRepository_Update_InvalidExperimentIDFormat(t *testing.T) {
	invalidUUIDs := []string{
		"not-a-uuid",
		"12345",
		"",
		"xyz-abc-def-ghi",
	}

	for _, invalidID := range invalidUUIDs {
		t.Run("invalid_experiment_id: "+invalidID, func(t *testing.T) {
			_, err := uuid.Parse(invalidID)
			if err == nil {
				t.Errorf("expected UUID parse error for '%s', got nil", invalidID)
			}
		})
	}
}

func TestRepository_Update_InvalidResultIDFormat(t *testing.T) {
	invalidUUIDs := []string{
		"not-a-uuid",
		"12345",
		"",
		"xyz-abc-def-ghi",
	}

	for _, invalidID := range invalidUUIDs {
		t.Run("invalid_result_id: "+invalidID, func(t *testing.T) {
			_, err := uuid.Parse(invalidID)
			if err == nil {
				t.Errorf("expected UUID parse error for '%s', got nil", invalidID)
			}
		})
	}
}

func TestRepository_Update_NoExperimentFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	resultID := uuid.New()
	experimentID := uuid.New() // Non-existent experiment
	currentVersion := 1

	result := &ExperimentResult{
		Outcome:         OutcomeFailure,
		Summary:         "Updated summary",
		OutcomeReason:   "Updated reason",
		ConfidenceLevel: ConfidenceMedium,
		UpdatedAt:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "experiment_results" SET`)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
	mock.ExpectCommit()

	err := repo.Update(ctx, resultID, experimentID, result, currentVersion)

	if err != ErrOptimisticLockingConflict {
		t.Errorf("expected ErrOptimisticLockingConflict, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRepository_Update_NoResultFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	resultID := uuid.New() // Non-existent result
	experimentID := uuid.New()
	currentVersion := 1

	result := &ExperimentResult{
		Outcome:         OutcomeFailure,
		Summary:         "Updated summary",
		OutcomeReason:   "Updated reason",
		ConfidenceLevel: ConfidenceMedium,
		UpdatedAt:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "experiment_results" SET`)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
	mock.ExpectCommit()

	err := repo.Update(ctx, resultID, experimentID, result, currentVersion)

	if err != ErrOptimisticLockingConflict {
		t.Errorf("expected ErrOptimisticLockingConflict, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRepository_Update_OptimisticLockingConflict(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	resultID := uuid.New()
	experimentID := uuid.New()
	staleVersion := 1 // Version is outdated (DB has version 2)

	result := &ExperimentResult{
		Outcome:         OutcomeFailure,
		Summary:         "Updated summary",
		OutcomeReason:   "Updated reason",
		ConfidenceLevel: ConfidenceMedium,
		UpdatedAt:       time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "experiment_results" SET`)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected due to version mismatch
	mock.ExpectCommit()

	err := repo.Update(ctx, resultID, experimentID, result, staleVersion)

	if err != ErrOptimisticLockingConflict {
		t.Errorf("expected ErrOptimisticLockingConflict, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
