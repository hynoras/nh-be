package procedure

import (
	"context"
	"errors"
	"nh-be/constant"
	"nh-be/utils"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, search string, offset, limit int, withExperiments bool) ([]Procedure, int64, error)
	FindByID(ctx context.Context, id uuid.UUID, withSteps, withExperiments bool) (*Procedure, error)
	CreateProcedure(ctx context.Context, procedure *Procedure) error
	// Update(ctx context.Context, procedure *Procedure) error
	// Delete(ctx context.Context, id uuid.UUID) error
	WithTransaction(ctx context.Context, fn func(repo Repository) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(ctx context.Context, search string, offset, limit int, withExperiments bool) ([]Procedure, int64, error) {
	var procedures []Procedure
	var length int64
	query := r.db.WithContext(ctx).Model(&Procedure{})

	if withExperiments {
		query = query.Preload("Experiments")
	}

	if search != "" {
		query = query.Where("LOWER(procedures.title) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	err := query.Count(&length).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Scopes(utils.Paginate(offset, limit)).Find(&procedures).Error
	return procedures, length, err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID, withSteps, withExperiments bool) (*Procedure, error) {
	var procedure Procedure
	query := r.db.WithContext(ctx).Model(&Procedure{})

	if withSteps {
		query = query.Preload("Steps")
	}

	if withExperiments {
		query = query.Preload("Experiments")
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

// func (r *repository) Update(ctx context.Context, procedure *Procedure) error {
// 	return r.db.WithContext(ctx).Save(procedure).Error
// }

// func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
// 	return r.db.WithContext(ctx).Delete(&Procedure{}, id).Error
// }

func (r *repository) WithTransaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &repository{db: tx}
		return fn(txRepo)
	})
}
