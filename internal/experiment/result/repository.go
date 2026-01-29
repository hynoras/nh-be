package result

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, r *ExperimentResult) error
	FindByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error)
	Update(ctx context.Context, experimentID uuid.UUID, r *ExperimentResult) error
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

func (r *repository) Update(ctx context.Context, experimentID uuid.UUID, result *ExperimentResult) error {
	fields := map[string]interface{}{
		"updated_at": result.UpdatedAt,
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

	return r.db.WithContext(ctx).
		Model(&ExperimentResult{}).
		Where("experiment_id = ?", experimentID).
		Updates(fields).Error
}
