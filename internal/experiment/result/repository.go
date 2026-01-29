package result

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, r *ExperimentResult) error
	FindByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error)
	FindByID(ctx context.Context, id uuid.UUID) (*ExperimentResult, error)
	Update(ctx context.Context, id uuid.UUID, experimentID uuid.UUID, r *ExperimentResult, currentVersion int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, result *ExperimentResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *repository) FindByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error) {
	var result ExperimentResult
	err := r.db.WithContext(ctx).
		Where("experiment_id = ?", experimentID).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*ExperimentResult, error) {
	var result ExperimentResult
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, experimentID uuid.UUID, result *ExperimentResult, currentVersion int) error {
	fields := map[string]interface{}{
		"updated_at": result.UpdatedAt,
		"version":    currentVersion + 1,
	}
	if result.Outcome != "" {
		fields["outcome"] = result.Outcome
	}
	if result.Summary != "" {
		fields["summary"] = result.Summary
	}
	if result.OutcomeReason != "" {
		fields["outcome_reason"] = result.OutcomeReason
	}
	if result.ConfidenceLevel != "" {
		fields["confidence_level"] = result.ConfidenceLevel
	}

	dbResult := r.db.WithContext(ctx).
		Model(&ExperimentResult{}).
		Where("id = ? AND experiment_id = ? AND version = ?", id, experimentID, currentVersion).
		Updates(fields)

	if dbResult.Error != nil {
		return dbResult.Error
	}

	if dbResult.RowsAffected == 0 {
		return ErrOptimisticLockingConflict
	}

	return nil
}
