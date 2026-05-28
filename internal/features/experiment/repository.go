package experiment

import (
	"context"
	"errors"
	"nh-be/internal/utils/dbutil"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, e *Experiment) error
	FindAll(ctx context.Context, search string, page, pageSize int) ([]Experiment, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Experiment, error)
	Update(ctx context.Context, id uuid.UUID, e *Experiment) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status ExperimentStatus, currentVersion int) error
	GetProcedureIDByID(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateProcedureID(ctx context.Context, id uuid.UUID, procedureID uuid.UUID, currentVersion int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, e *Experiment) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *repository) FindAll(ctx context.Context, search string, page, pageSize int) ([]Experiment, int64, error) {
	var experiments []Experiment
	var length int64

	query := r.db.WithContext(ctx).Model(&Experiment{}).
		Preload("CreatedBy").
		Where("LOWER(title) LIKE ?", "%"+strings.ToLower(search)+"%")

	if err := query.Count(&length).Error; err != nil {
		return nil, 0, err
	}

	result := query.Scopes(dbutil.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&experiments)

	return experiments, length, result.Error
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*Experiment, error) {
	var e Experiment
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("CreatedBy").
		First(&e)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrExperimentNotFound
	}
	return &e, nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, e *Experiment) error {
	fields := map[string]interface{}{
		"updated_at": e.UpdatedAt,
	}
	if e.Title != "" {
		fields["title"] = e.Title
	}
	if e.Type != "" {
		fields["type"] = e.Type
	}
	if e.Objective != "" {
		fields["objective"] = e.Objective
	}

	return r.db.WithContext(ctx).Model(&Experiment{}).Where("id = ?", id).Updates(fields).Error
}

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status ExperimentStatus, currentVersion int) error {
	result := r.db.WithContext(ctx).
		Model(&Experiment{}).
		Where("id = ? AND version = ?", id, currentVersion).
		Updates(map[string]interface{}{
			"status":  status,
			"version": currentVersion + 1,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrExperimentConflict
	}

	return nil
}

func (r *repository) GetProcedureIDByID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var e Experiment
	result := r.db.WithContext(ctx).
		Model(&Experiment{}).
		Select("procedure_id").
		Where("id = ?", id).
		First(&e)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return uuid.Nil, result.Error
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return uuid.Nil, ErrExperimentNotFound
	}
	if e.ProcedureID == nil {
		return uuid.Nil, nil
	}
	return *e.ProcedureID, nil
}

func (r *repository) UpdateProcedureID(ctx context.Context, id uuid.UUID, procedureID uuid.UUID, currentVersion int) error {
	result := r.db.WithContext(ctx).
		Model(&Experiment{}).
		Where("id = ? AND version = ?", id, currentVersion).
		Updates(map[string]interface{}{
			"procedure_id": procedureID,
			"version":      currentVersion + 1,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrExperimentConflict
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Experiment{}).Error
}
