package result

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, r *ExperimentResult) error
	FindByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error)
	FindByID(ctx context.Context, id uuid.UUID) (*ExperimentResult, error)
	FindByIDAndExperimentID(ctx context.Context, id uuid.UUID, experimentID uuid.UUID) (*ExperimentResult, error)
	Update(ctx context.Context, id uuid.UUID, experimentID uuid.UUID, r *UpdateFields, currentVersion int) error
}

// UpdateFields holds pointer fields for partial updates.
// nil means "don't update", non-nil means "set to this value".
type UpdateFields struct {
	Outcome         *Outcome
	Summary         *string
	OutcomeReason   *string
	ConfidenceLevel *ConfidenceLevel
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, result *ExperimentResult) error {
	err := r.db.WithContext(ctx).Create(result).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "experiment_id") {
			return ErrExperimentResultAlreadyExists
		}
		if strings.Contains(err.Error(), "foreign key constraint") && strings.Contains(err.Error(), "experiment_id") {
			return ErrExperimentNotFound
		}
		return err
	}
	return nil
}

func (r *repository) FindByExperimentID(ctx context.Context, experimentID uuid.UUID) (*ExperimentResult, error) {
	var result ExperimentResult
	err := r.db.WithContext(ctx).
		Where("experiment_id = ?", experimentID).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExperimentResultNotFound
		}
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExperimentResultNotFound
		}
		return nil, err
	}
	return &result, nil
}

func (r *repository) FindByIDAndExperimentID(ctx context.Context, id uuid.UUID, experimentID uuid.UUID) (*ExperimentResult, error) {
	var result ExperimentResult
	err := r.db.WithContext(ctx).
		Where("id = ? AND experiment_id = ?", id, experimentID).
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExperimentResultNotFound
		}
		return nil, err
	}
	return &result, nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, experimentID uuid.UUID, fields *UpdateFields, currentVersion int) error {
	updates := map[string]interface{}{
		"version": currentVersion + 1,
	}
	if fields.Outcome != nil {
		updates["outcome"] = *fields.Outcome
	}
	if fields.Summary != nil {
		updates["summary"] = *fields.Summary
	}
	if fields.OutcomeReason != nil {
		updates["outcome_reason"] = *fields.OutcomeReason
	}
	if fields.ConfidenceLevel != nil {
		updates["confidence_level"] = *fields.ConfidenceLevel
	}

	dbResult := r.db.WithContext(ctx).
		Model(&ExperimentResult{}).
		Where("id = ? AND experiment_id = ? AND version = ?", id, experimentID, currentVersion).
		Updates(updates)

	if dbResult.Error != nil {
		return dbResult.Error
	}

	if dbResult.RowsAffected == 0 {
		return ErrOptimisticLockingConflict
	}

	return nil
}
