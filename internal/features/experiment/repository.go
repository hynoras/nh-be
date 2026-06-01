package experiment

import (
	"context"
	"errors"
	"nh-be/internal/constant"
	"nh-be/internal/utils/dbutil"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(ctx context.Context, e *Experiment) error
	FindAllExperiments(
		ctx context.Context,
		sortBy, search string,
		sortOrder constant.Order,
		experimentStatus *ExperimentStatus,
		experimentType *ExperimentType,
		createdBy *uuid.UUID,
		page, pageSize int,
	) ([]ExperimentsQueryDto, error)
	CountExperiments(ctx context.Context, createdBy *uuid.UUID) (int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Experiment, error)
	Update(ctx context.Context, id uuid.UUID, e *Experiment) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status ExperimentStatus, currentVersion int) error
	GetProcedureIDByID(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateProcedureID(ctx context.Context, id uuid.UUID, procedureID uuid.UUID, currentVersion int) error
	Delete(ctx context.Context, id uuid.UUID) error
	WithTransaction(ctx context.Context, fn func(repo Repository) error) error
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

func (r *repository) FindAllExperiments(
	ctx context.Context,
	sortBy, search string,
	sortOrder constant.Order,
	experimentStatus *ExperimentStatus,
	experimentType *ExperimentType,
	createdBy *uuid.UUID,
	page, pageSize int,
) ([]ExperimentsQueryDto, error) {
	var experiments []ExperimentsQueryDto

	query := r.db.WithContext(ctx).Model(&Experiment{}).
		Select(
			"e.identifier",
			"e.title",
			"e.objective",
			"e.status",
			"e.type",
			"creator.name as creator",
			"updater.name as updater",
			"e.created_at",
			"e.updated_at",
			"proc.name as procedure_name",
		).
		Joins("INNER JOIN users creator ON creator.id = experiments.created_by_id").
		Joins("LEFT JOIN users updater ON updater.id = experiments.updated_by_id").
		Joins("LEFT JOIN procedures proc ON proc.id = experiments.procedure_id").
		Where("LOWER(title) LIKE ?", "%"+strings.ToLower(search)+"%")

	if experimentStatus != nil {
		query = query.Where("status = ?", *experimentStatus)
	}
	if experimentType != nil {
		query = query.Where("type = ?", *experimentType)
	}
	if createdBy != nil {
		query = query.Where("created_by_id = ?", *createdBy)
	}

	allowedSortFields := map[string]string{
		"created_at": "e.created_at",
		"updated_at": "e.updated_at",
	}

	column, ok := allowedSortFields[sortBy]
	if !ok {
		column = "e.created_at"
	}

	result := query.Scopes(dbutil.Paginate(page, pageSize)).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: column},
			Desc:   sortOrder == constant.DESC,
		}).
		Scan(&experiments)

	return experiments, result.Error
}

func (r *repository) CountExperiments(ctx context.Context, createdBy *uuid.UUID) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&Experiment{})

	//NOTE: uncomment if business requirement changes
	//For now just return count filtered by createdBy
	// .Where("LOWER(title) LIKE ?", "%"+strings.ToLower(search)+"%")

	// if experimentStatus != nil {
	// 	query = query.Where("status = ?", *experimentStatus)
	// }
	// if experimentType != nil {
	// 	query = query.Where("type = ?", *experimentType)
	// }
	if createdBy != nil {
		query = query.Where("created_by_id = ?", *createdBy)
	}

	result := query.Count(&count)
	return count, result.Error
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

func (r *repository) WithTransaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&repository{db: tx})
	})
}
