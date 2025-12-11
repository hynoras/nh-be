package permission

import (
	"context"
	"nh-be/utils"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Permission
	FindIDByID(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	FindIDByIDs(ctx context.Context, ids []uuid.UUID) ([]uuid.UUID, error)
	FindAllPermissions(ctx context.Context, name string) ([]Permission, int64, error)
	FindPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error)
	FindPermissionsByIDs(ctx context.Context, ids []uuid.UUID) ([]Permission, error)

	// Permission Group
	CreatePermissionGroup(ctx context.Context, pg *PermissionGroup) error
	FindAllPermissionGroups(ctx context.Context, name string, assignedUser string, offset, limit int) ([]PermissionGroup, int64, error)
	FindPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error)
	UpdatePermissionGroup(ctx context.Context, id uuid.UUID, pg *PermissionGroup) error
	DeletePermissionGroup(ctx context.Context, id uuid.UUID) error
	FindPermissionGroupsByIDs(ctx context.Context, ids []uuid.UUID) ([]PermissionGroup, error)

	// Transaction
	WithTransaction(ctx context.Context, fn func(repo Repository) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Permission Implementations
func (r *repository) FindAllPermissions(ctx context.Context, search string) ([]Permission, int64, error) {
	var permissions []Permission
	var length int64

	query := r.db.WithContext(ctx).Select("id", "name", "description").Model(&Permission{}).
	Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%").
	Count(&length)

	result := query.Find(&permissions)
	return permissions, length, result.Error
}

func (r *repository) FindIDByID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var p uuid.UUID
	if err := r.db.WithContext(ctx).Model(&Permission{}).Select("id").Where("id = ?", id).First(&p).Error; err != nil {
		return uuid.Nil, err
	}
	return p, nil
}

func (r *repository) FindIDByIDs(ctx context.Context, ids []uuid.UUID) ([]uuid.UUID, error) {
	var p []uuid.UUID
	if err := r.db.WithContext(ctx).Model(&Permission{}).Select("id").Where("id IN ?", ids).Find(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *repository) FindPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error) {
	var p Permission
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) FindPermissionsByIDs(ctx context.Context, ids []uuid.UUID) ([]Permission, error) {
	var permissions []Permission
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// Permission Group Implementations
func (r *repository) CreatePermissionGroup(ctx context.Context, pg *PermissionGroup) error {
	return r.db.WithContext(ctx).Create(pg).Error
}

func (r *repository) FindAllPermissionGroups(
	ctx context.Context,
	name string,
	assignedUser string,
	offset,
	limit int,
) ([]PermissionGroup, int64, error) {
	var groups []PermissionGroup
	var length int64

	query := r.db.WithContext(ctx).Model(&PermissionGroup{})

	if name != "" {
		query = query.Where("permission_groups.name LIKE ?", "%"+name+"%")
	}
	
	if assignedUser != "" {
		query = query.Joins("JOIN user_permissions ON user_permissions.permission_group_id = permission_groups.id").
			Joins("JOIN users ON users.id = user_permissions.user_id").
			Where("users.username LIKE ?", "%"+assignedUser+"%").
			Distinct()
	}

	query.Count(&length)

	result := query.
		Preload("Permissions").
		Preload("AssignedUsers").
		Scopes(utils.Paginate(offset, limit)).
		Find(&groups)
	
	return groups, length, result.Error
}

func (r *repository) FindPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error) {
	var pg PermissionGroup
	if err := r.db.WithContext(ctx).Preload("Permissions").Preload("AssignedUsers").First(&pg, id).Error; err != nil {
		return nil, err
	}
	return &pg, nil
}

func (r *repository) UpdatePermissionGroup(ctx context.Context, id uuid.UUID, pg *PermissionGroup) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		basicFields := map[string]interface{}{
			"name":        pg.Name,
			"description": pg.Description,
		}
		if err := tx.Model(&PermissionGroup{}).Where("id = ?", id).Updates(basicFields).Error; err != nil {
			return err
		}
		
		pg.ID = id

		if err := tx.Model(pg).Association("Permissions").Replace(pg.Permissions); err != nil {
			return err
		}
		
		return nil
	})
}

func (r *repository) DeletePermissionGroup(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var pg PermissionGroup
		if err := tx.First(&pg, id).Error; err != nil {
			return err
		}
		
		if err := tx.Model(&pg).Association("Permissions").Clear(); err != nil {
			return err
		}
		
		if err := tx.Model(&pg).Association("AssignedUsers").Clear(); err != nil {
			return err
		}
		
		return tx.Delete(&pg).Error
	})
}

func (r *repository) FindPermissionGroupsByIDs(ctx context.Context, ids []uuid.UUID) ([]PermissionGroup, error) {
	var permissionGroups []PermissionGroup
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&permissionGroups).Error; err != nil {
		return nil, err
	}
	return permissionGroups, nil
}

// Transaction Implementations
func (r *repository) WithTransaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &repository{db: tx}
		return fn(txRepo)
	})
}
