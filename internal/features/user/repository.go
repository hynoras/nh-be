package user

import (
	"context"
	"errors"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/dbutil"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, u *User) (User, error)
	FindAll(ctx context.Context, search string, page, pageSize int) ([]User, int64, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]User, error)
	FindPasswordById(ctx context.Context, id uuid.UUID) (*string, error)
	FindPermissionGroupsByIDs(ctx context.Context, ids []uuid.UUID) ([]permission.PermissionGroup, error)
	Update(ctx context.Context, id uuid.UUID, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	WithTransaction(ctx context.Context, fn func(repo Repository) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, u *User) (User, error) {
	err := r.db.WithContext(ctx).Create(&u).Error
	if err != nil {
		return User{}, err
	}
	return *u, nil
}

func (r *repository) FindAll(ctx context.Context, search string, page, pageSize int) ([]User, int64, error) {
	var users []User
	var length int64

	baseQuery := r.db.WithContext(ctx).Model(&User{})
	if search != "" {
		baseQuery = baseQuery.Where("LOWER(username) LIKE ?", "%"+strings.ToLower(search)+"%")
	}
	if err := baseQuery.Count(&length).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Preload("AssignedPermissionGroups").
		Select("id", "username", "email", "created_at").
		Scopes(dbutil.Paginate(page, pageSize)).
		Find(&users).Error

	return users, length, err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&u)

	//replace gorm error with domain error
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &u, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&u)

	//replace gorm error with domain error
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}

func (r *repository) FindPasswordById(ctx context.Context, id uuid.UUID) (*string, error) {
	var u User
	result := r.db.WithContext(ctx).Where("id = ?", id).Select("password").First(&u)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return u.Password, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User
	result := r.db.WithContext(ctx).Where("id = ?", id).
		Preload("AssignedPermissionGroups").
		Preload("AssignedPermissionGroups.Permissions").
		First(&u)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &u, nil
}

func (r *repository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]User, error) {
	var u []User
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (r *repository) FindPermissionGroupsByIDs(ctx context.Context, ids []uuid.UUID) ([]permission.PermissionGroup, error) {
	var groups []permission.PermissionGroup
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}
	return groups, nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, u *User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update basic fields
		basicFields := map[string]interface{}{}
		if u.Username != "" {
			basicFields["username"] = u.Username
		}
		if u.Email != "" {
			basicFields["email"] = u.Email
		}
		if u.IsVerified {
			basicFields["is_verified"] = u.IsVerified
		}
		basicFields["updated_at"] = u.UpdatedAt

		if err := tx.Model(&User{}).Where("id = ?", id).Updates(basicFields).Error; err != nil {
			return err
		}

		// Update permission groups association if provided
		if u.AssignedPermissionGroups != nil {
			u.ID = id
			if err := tx.Model(u).Association("AssignedPermissionGroups").Replace(u.AssignedPermissionGroups); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}

func (r *repository) WithTransaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &repository{db: tx}
		return fn(txRepo)
	})
}
