package permission

import (
	"context"
	"nh-be/constant"
	"nh-be/utils"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	// Permission
	GetAllPermissions(ctx context.Context, search string) ([]Permission, int64, error)
	GetPermissionsByIDs(ctx context.Context, permissionIds []uuid.UUID) ([]Permission, error)
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error)
	GetUserPermissionCodeNames(ctx context.Context, userId uuid.UUID) ([]string, error)

	// Permission Group
	CreatePermissionGroup(ctx context.Context, permissionGroup *PermissionGroupInput) error
	GetAllPermissionGroups(ctx context.Context, search string, permissionIds []uuid.UUID, page, pageSize int) ([]PermissionGroup, int64, error)
	GetPermissionGroupsByIDs(ctx context.Context, permissionGroupIds []uuid.UUID) ([]PermissionGroup, error)
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

func (s *service) GetPermissionsByIDs(ctx context.Context, permissionIds []uuid.UUID) ([]Permission, error) {
	permissions, err := s.permissionRepo.FindPermissionsByIDs(ctx, permissionIds)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *service) GetPermissionGroupsByIDs(ctx context.Context, permissionGroupIds []uuid.UUID) ([]PermissionGroup, error) {
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
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	userPerm, err := s.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	if !slices.Contains(userPerm, constant.ViewPermissionGroup) && !slices.Contains(userPerm, constant.ManagePermissionGroup) {
		return nil, 0, ErrForbidViewPermissions
	}

	return s.permissionRepo.FindAllPermissions(ctx, search)
}

func (s *service) GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error) {
	return s.permissionRepo.FindPermissionByID(ctx, id)
}

func (s *service) GetUserPermissionCodeNames(ctx context.Context, userId uuid.UUID) ([]string, error) {
	return s.permissionRepo.FindUserPermissionCodeNames(ctx, userId)
}

// Permission Group Implementations
func (s *service) CreatePermissionGroup(ctx context.Context, permissionGroup *PermissionGroupInput) error {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManagePermissionGroup) {
		return ErrForbidCreatePermissionGroup
	}

	var permissions []Permission

	var checkExistingPermissionsErr error
	permissions, checkExistingPermissionsErr = s.GetPermissionsByIDs(ctx, permissionGroup.Permissions)
	if checkExistingPermissionsErr != nil {
		return checkExistingPermissionsErr
	}

	// Check for duplicate name
	existingGroup, err := s.permissionRepo.FindPermissionGroupByName(ctx, permissionGroup.Name)
	if err == nil && existingGroup != nil {
		return ErrPermissionGroupNameAlreadyExists
	}

	return s.permissionRepo.WithTransaction(ctx, func(txRepo Repository) error {
		pg := &PermissionGroup{
			Name:        permissionGroup.Name,
			Description: permissionGroup.Description,
			Permissions: permissions,
			CreatedAt:   time.Now(),
		}

		err := txRepo.CreatePermissionGroup(ctx, pg)

		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetAllPermissionGroups(ctx context.Context, search string, permissionIds []uuid.UUID, page, pageSize int) ([]PermissionGroup, int64, error) {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	userPerm, err := s.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	if !slices.Contains(userPerm, constant.ViewPermissionGroup) && !slices.Contains(userPerm, constant.ManagePermissionGroup) {
		return nil, 0, ErrForbidViewPermissionGroup
	}

	offset := (page - 1) * pageSize
	return s.permissionRepo.FindAllPermissionGroups(ctx, search, permissionIds, offset, pageSize)
}

func (s *service) GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error) {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userPerm, err := s.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(userPerm, constant.ViewPermissionGroup) && !slices.Contains(userPerm, constant.ManagePermissionGroup) {
		return nil, ErrForbidViewPermissionGroup
	}

	permissionGroup, err := s.permissionRepo.FindPermissionGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return permissionGroup, nil
}

func (s *service) UpdatePermissionGroup(ctx context.Context, id uuid.UUID, permissionGroup *PermissionGroupInput) error {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManagePermissionGroup) {
		return ErrForbidUpdatePermissionGroup
	}

	var permissions []Permission

	_, err = s.permissionRepo.FindPermissionGroupByID(ctx, id)
	if err != nil {
		return err
	}

	var checkExistingPermissionsErr error
	permissions, checkExistingPermissionsErr = s.GetPermissionsByIDs(ctx, permissionGroup.Permissions)
	if checkExistingPermissionsErr != nil {
		return checkExistingPermissionsErr
	}

	// Check for duplicate name (excluding current group)
	duplicateGroup, err := s.permissionRepo.FindPermissionGroupByName(ctx, permissionGroup.Name)
	if err == nil && duplicateGroup != nil && duplicateGroup.ID != id {
		return ErrPermissionGroupNameAlreadyExists
	}

	pg := &PermissionGroup{
		Name:        permissionGroup.Name,
		Description: permissionGroup.Description,
		Permissions: permissions,
		UpdatedAt:   time.Now(),
	}
	return s.permissionRepo.UpdatePermissionGroup(ctx, id, pg)
}

func (s *service) DeletePermissionGroup(ctx context.Context, id uuid.UUID) error {
	userId, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManagePermissionGroup) {
		return ErrForbidDeletePermissionGroup
	}

	permissionGroup, err := s.permissionRepo.FindPermissionGroupByID(ctx, id)
	if err != nil {
		return err
	}

	if permissionGroup.Name == "Super Admin" {
		return ErrCannotDeleteSuperAdmin
	}

	return s.permissionRepo.DeletePermissionGroup(ctx, id)
}
