package procedure

import (
	"context"
	"errors"
	"nh-be/internal/constant"
	"nh-be/internal/utils/dbutil"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, search string, offset, limit int) ([]Procedure, int64, error)
	FindByID(ctx context.Context, id uuid.UUID, withSteps bool) (*Procedure, error)
	CreateProcedure(ctx context.Context, procedure *Procedure) error
	UpdateProcedure(ctx context.Context, id uuid.UUID, procedure *Procedure) error
	DeleteProcedure(ctx context.Context, id uuid.UUID) error

	GetStepIDsByProcID(ctx context.Context, procedureId uuid.UUID) ([]StepMetadata, error)
	CreateProcedureStep(ctx context.Context, step *ProcedureStep) error
	UpdateProcedureStep(ctx context.Context, stepId uuid.UUID, procedureId uuid.UUID, step *ProcedureStep) error
	DeleteProcedureStep(ctx context.Context, stepId uuid.UUID, procedureId uuid.UUID) error
	WithTransaction(ctx context.Context, fn func(repo Repository) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(ctx context.Context, search string, offset, limit int) ([]Procedure, int64, error) {
	var procedures []Procedure
	var length int64
	query := r.db.WithContext(ctx).Model(&Procedure{})

	if search != "" {
		query = query.Where("LOWER(procedures.title) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	err := query.Count(&length).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Scopes(dbutil.Paginate(offset, limit)).Find(&procedures).Error
	return procedures, length, err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID, withSteps bool) (*Procedure, error) {
	var procedure Procedure
	query := r.db.WithContext(ctx).Model(&Procedure{})

	if withSteps {
		query = query.Preload("Steps")
	}

	err := query.First(&procedure, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrProcedureNotFound
		}
		return nil, err
	}

	return &procedure, nil
}

func (r *repository) CreateProcedure(ctx context.Context, procedure *Procedure) error {
	return r.db.WithContext(ctx).Create(procedure).Error
}

func (r *repository) UpdateProcedure(ctx context.Context, id uuid.UUID, procedure *Procedure) error {
	fields := map[string]interface{}{
		"title":       procedure.Title,
		"description": procedure.Description,
		"updated_at":  procedure.UpdatedAt,
		"version":     gorm.Expr("version + 1"),
	}

	dbResult := r.db.WithContext(ctx).Model(&Procedure{}).Where("id = ? AND version = ?", id, procedure.Version).Updates(fields)
	if dbResult.Error != nil {
		return dbResult.Error
	}

	if dbResult.RowsAffected == 0 {
		var count int64
		err := r.db.Model(&Procedure{}).Where("id = ?", id).Count(&count).Error
		if err != nil {
			return err
		}

		if count == 0 {
			return constant.ErrProcedureNotFound
		}
		return constant.ErrProcedureConflict
	}
	return nil
}

func (r *repository) DeleteProcedure(ctx context.Context, id uuid.UUID) error {
	dbResult := r.db.WithContext(ctx).Model(&Procedure{}).Where("id = ?", id).Delete(&Procedure{})
	if dbResult.Error != nil {
		return dbResult.Error
	}
	if dbResult.RowsAffected == 0 {
		return constant.ErrProcedureNotFound
	}
	return nil
}

func (r *repository) CreateProcedureStep(ctx context.Context, step *ProcedureStep) error {
	err := r.db.WithContext(ctx).Create(step).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetStepIDsByProcID(ctx context.Context, procedureId uuid.UUID) ([]StepMetadata, error) {
	var steps []StepMetadata
	err := r.db.WithContext(ctx).
		Model(&ProcedureStep{}).
		Where("procedure_id = ?", procedureId).
		Select("id", "version").
		Scan(&steps).Error

	if err != nil {
		return nil, err
	}

	return steps, nil
}

func (r *repository) UpdateProcedureStep(ctx context.Context, stepId uuid.UUID, procedureId uuid.UUID, step *ProcedureStep) error {
	fields := map[string]interface{}{
		"title":       step.Title,
		"index":       step.Index,
		"description": step.Description,
		"updated_at":  step.UpdatedAt,
		"step_type":   step.StepType,
		"version":     gorm.Expr("version + 1"),
	}

	dbResult :=
		r.db.WithContext(ctx).
			Model(&ProcedureStep{}).
			Where("id = ? AND procedure_id = ? AND version = ?", stepId, procedureId, step.Version).
			Updates(fields)
	if dbResult.Error != nil {
		return dbResult.Error
	}

	if dbResult.RowsAffected == 0 {
		var count int64
		err := r.db.WithContext(ctx).
			Model(&ProcedureStep{}).
			Where("id = ? AND procedure_id = ?", stepId, procedureId).
			Count(&count).Error
		if err != nil {
			return err
		}

		if count == 0 {
			return constant.ErrProcedureStepNotFound
		}
		return constant.ErrProcedureStepConflict
	}

	return nil
}

func (r *repository) DeleteProcedureStep(ctx context.Context, stepId uuid.UUID, procedureId uuid.UUID) error {
	dbResult := r.db.WithContext(ctx).
		Model(&ProcedureStep{}).
		Where("id = ? AND procedure_id = ?", stepId, procedureId).
		Delete(&ProcedureStep{})

	if dbResult.Error != nil {
		return dbResult.Error
	}

	if dbResult.RowsAffected == 0 {
		return constant.ErrProcedureStepNotFound
	}

	return nil
}

func (r *repository) WithTransaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &repository{db: tx}
		return fn(txRepo)
	})
}
