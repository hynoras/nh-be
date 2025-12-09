package permission

import (
	"context"
	"nh-be/utils"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	// Permission
	FindAllPermissions(ctx context.Context, name string) ([]Permission, int64, error)
	FindPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error)
	FindPermissionsByIDs(ctx context.Context, ids []uuid.UUID) ([]Permission, error)

	// Permission Group
	CreatePermissionGroup(ctx context.Context, pg *PermissionGroup) error
	FindAllPermissionGroups(ctx context.Context, name string, assignedUser uuid.UUID, offset, limit int) ([]PermissionGroup, int64, error)
	FindPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error)
	UpdatePermissionGroup(ctx context.Context, id uuid.UUID, pg *PermissionGroup) error
	DeletePermissionGroup(ctx context.Context, id uuid.UUID) error

	// User Permission
	AssignUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
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

	query := r.db.WithContext(ctx).Model(&Permission{}).Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	query.Count(&length)

	result := query.Find(&permissions).Error
	return permissions, length, result
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

func (r *repository) FindAllPermissionGroups(ctx context.Context, name string, assignedUser uuid.UUID, offset, limit int) ([]PermissionGroup, int64, error) {
	var groups []PermissionGroup
	var length int64
	query := r.db.WithContext(ctx).Model(&PermissionGroup{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if assignedUser != uuid.Nil {
		query = query.Joins("JOIN user_permissions ON user_permissions.permission_group_id = permission_groups.id").Where("user_permissions.user_id = ?", assignedUser)
	}
	query.Find(&groups).Count(&length)
	result := query.Preload("Permissions").Scopes(utils.Paginate(offset, limit)).Find(&groups)
	return groups, length, result.Error
}

func (r *repository) FindPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error) {
	var pg PermissionGroup
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&pg, id).Error; err != nil {
		return nil, err
	}
	return &pg, nil
}

func (r *repository) UpdatePermissionGroup(ctx context.Context, id uuid.UUID, pg *PermissionGroup) error {
    // For many2many updates, we might need to handle association replacement carefully.
    // However, basic Updates only touches fields. To replace permissions, we usually use Association replace.
    // Here we just update basic fields. Association update is better handled in Service by interacting with Associations.
    // But for simplicity in repo, let's assume pg has ID set.
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Updates(pg).Error
}

func (r *repository) DeletePermissionGroup(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Select(clause.Associations).Delete(&PermissionGroup{}).Error
}

// User Permission Implementations
func (r *repository) AssignUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	up := UserPermission{UserID: userID, PermissionGroupID: groupID}
	return r.db.WithContext(ctx).Create(&up).Error
}

func (r *repository) RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND permission_group_id = ?", userID, groupID).Delete(&UserPermission{}).Error
}
