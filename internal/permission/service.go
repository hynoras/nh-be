package permission

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrPermissionNotFound = errors.New("permission not found")
	ErrPermissionGroupNotFound = errors.New("permission group not found")
	ErrNotNullPermissions = errors.New("permissions can not be null")
)

type Service interface {
	CheckExistingPermission(ctx context.Context, permissionId uuid.UUID) (*Permission, error)
	CheckExistingPermissions(ctx context.Context, permissionIds []uuid.UUID) ([]Permission, error)
	CheckExistingPermissionGroups(ctx context.Context, permissionGroupIds []uuid.UUID) ([]PermissionGroup, error)
	// Permission
	GetAllPermissions(ctx context.Context, search string) ([]Permission, int64, error)
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error)

	// Permission Group
	CreatePermissionGroup(ctx context.Context, permissionGroup *PermissionGroupInput) error
	GetAllPermissionGroups(ctx context.Context, search string, page, pageSize int) ([]PermissionGroup, int64, error)
	GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error)
	UpdatePermissionGroup(ctx context.Context, id uuid.UUID, permissionGroup *PermissionGroupInput) error
	DeletePermissionGroup(ctx context.Context, id uuid.UUID) error
}

type service struct {
	permissionRepo Repository
}

func NewService(permissionRepo Repository) Service {
	return &service{permissionRepo: permissionRepo}
}

func (s *service) CheckExistingPermission(ctx context.Context, permissionId uuid.UUID) (*Permission, error) {
	permission, err := s.permissionRepo.FindPermissionByID(ctx, permissionId)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrPermissionNotFound
	}
	return permission, nil
}

func (s *service) CheckExistingPermissions(ctx context.Context, permissionIds []uuid.UUID) ([]Permission, error) {
	permissions, err := s.permissionRepo.FindPermissionsByIDs(ctx, permissionIds)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *service) CheckExistingPermissionGroups(ctx context.Context, permissionGroupIds []uuid.UUID) ([]PermissionGroup, error) {
	permissionGroups, err := s.permissionRepo.FindPermissionGroupsByIDs(ctx, permissionGroupIds)
	if err != nil {
		return nil, err
	}
	if len(permissionGroups) == 0 {
		return nil, ErrPermissionGroupNotFound
	}
	return permissionGroups, nil
}

// Permission Implementations
func (s *service) GetAllPermissions(ctx context.Context, search string) ([]Permission, int64, error) {
	return s.permissionRepo.FindAllPermissions(ctx, search)
}

func (s *service) GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error) {
	return s.permissionRepo.FindPermissionByID(ctx, id)
}

// Permission Group Implementations
func (s *service) CreatePermissionGroup(ctx context.Context, permissionGroup *PermissionGroupInput) error {
	var permissions []Permission
	
	var checkExistingPermissionsErr error
	permissions, checkExistingPermissionsErr = s.CheckExistingPermissions(ctx, permissionGroup.Permissions)
	if checkExistingPermissionsErr != nil {
		return checkExistingPermissionsErr
	}

	return s.permissionRepo.WithTransaction(ctx, func(txRepo Repository) error {
		pg := &PermissionGroup{
			Name:          permissionGroup.Name,
			Description:   permissionGroup.Description,
			Permissions:   permissions,
		}

		err := txRepo.CreatePermissionGroup(ctx, pg)

		if err != nil {
			return err
		} 
		
		return nil 
	})
}

func (s *service) GetAllPermissionGroups(ctx context.Context, search string, page, pageSize int) ([]PermissionGroup, int64, error) {
	offset := (page - 1) * pageSize
	return s.permissionRepo.FindAllPermissionGroups(ctx, search, offset, pageSize)
}

func (s *service) GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error) {
	permissionGroup, err := s.permissionRepo.FindPermissionGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return permissionGroup, nil
}

func (s *service) UpdatePermissionGroup(ctx context.Context, id uuid.UUID, permissionGroup *PermissionGroupInput) error {
	var permissions []Permission

	_, err := s.GetPermissionGroupByID(ctx, id)
	if err != nil {
		return err
	}
	
	var checkExistingPermissionsErr error
	permissions, checkExistingPermissionsErr = s.CheckExistingPermissions(ctx, permissionGroup.Permissions)
	if checkExistingPermissionsErr != nil {
		return checkExistingPermissionsErr
	}
	
	pg := &PermissionGroup{
		Name:        permissionGroup.Name,
		Description: permissionGroup.Description,
		Permissions: permissions,
	}
	return s.permissionRepo.UpdatePermissionGroup(ctx, id, pg)
}

func (s *service) DeletePermissionGroup(ctx context.Context, id uuid.UUID) error {
	_, err := s.GetPermissionGroupByID(ctx, id)
	if err != nil {
		return err
	}
	
	return s.permissionRepo.DeletePermissionGroup(ctx, id)
}
