package observation

import (
	"context"
	"fmt"
	"nh-be/internal/constant"
	exp "nh-be/internal/features/experiment"
	proc "nh-be/internal/features/procedure"
	"nh-be/internal/utils/dbutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllObsByExpIDAndProcID(
		ctx context.Context,
		expId, procId uuid.UUID,
		offset, limit int,
		sortBy *string,
		sortOrder *constant.Order,
	) ([]ObservationMetadata, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllObsByExpIDAndProcID(
	ctx context.Context,
	expId, procId uuid.UUID,
	offset, limit int,
	sortBy *string,
	sortOrder *constant.Order,
) ([]ObservationMetadata, int64, error) {
	err := r.db.WithContext(ctx).
		Model(&exp.Experiment{}).
		Select("id").
		Where("id = ?", expId).
		Take(&exp.Experiment{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []ObservationMetadata{}, 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return []ObservationMetadata{}, 0, constant.ErrExperimentNotFound
	}

	err = r.db.WithContext(ctx).
		Model(&proc.ProcedureStep{}).
		Select("id").
		Where("id = ?", procId).
		Take(&proc.ProcedureStep{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []ObservationMetadata{}, 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return []ObservationMetadata{}, 0, constant.ErrProcedureNotFound
	}

	baseQuery := r.db.WithContext(ctx).
		Model(&Observation{}).
		Where("procedure_step_id = ? AND experiment_id = ?", procId, expId)

	var length int64
	countErr := baseQuery.Count(&length).Error
	if countErr != nil {
		return []ObservationMetadata{}, 0, countErr
	}

	allowedColumns := map[string]bool{"created_at": true}
	if sortBy != nil && sortOrder != nil {
		if allowedColumns[*sortBy] {
			baseQuery = baseQuery.Order(fmt.Sprintf("%s %s", *sortBy, string(*sortOrder)))
		}
	}

	var observation []ObservationMetadata
	findErr := baseQuery.
		Scopes(dbutil.Paginate(offset, limit)).
		Select("id, observed_at, title, notes, previous_observation_id, created_by, created_at").
		Find(&observation).Error
	if findErr != nil {
		return []ObservationMetadata{}, 0, findErr
	}

	return observation, length, nil
}
