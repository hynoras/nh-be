package procedure

import (
	"context"
	"nh-be/utils"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, search string, offset, limit int, withExperiments bool) ([]Procedure, int64, error)
	// FindByID(ctx context.Context, id uuid.UUID) (*Procedure, error)
	// Create(ctx context.Context, procedure *Procedure) error
	// Update(ctx context.Context, procedure *Procedure) error
	// Delete(ctx context.Context, id uuid.UUID) error
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

// func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*Procedure, error) {
// 	var procedure Procedure
// 	err := r.db.WithContext(ctx).Preload("Steps").Preload("Experiments").First(&procedure, id).Error
// 	return &procedure, err
// }

// func (r *repository) Create(ctx context.Context, procedure *Procedure) error {
// 	return r.db.WithContext(ctx).Create(procedure).Error
// }

// func (r *repository) Update(ctx context.Context, procedure *Procedure) error {
// 	return r.db.WithContext(ctx).Save(procedure).Error
// }

// func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
// 	return r.db.WithContext(ctx).Delete(&Procedure{}, id).Error
// }
