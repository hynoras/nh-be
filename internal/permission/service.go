package permission

import (
	"context"
	"errors"

	"nh-be/internal/user"
	"nh-be/utils"

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
	// Permission
	GetAllPermissions(ctx context.Context, search string) ([]PermissionResponseDto, int64, error)
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error)

	// Permission Group
	CreatePermissionGroup(ctx context.Context, dto *CreatePermissionGroupDto) error
	GetAllPermissionGroups(ctx context.Context, name string, assignedUser string, page, pageSize int) ([]PermissionGroup, int64, error)
	GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error)
	UpdatePermissionGroup(ctx context.Context, id uuid.UUID, dto *UpdatePermissionGroupDto) error
	DeletePermissionGroup(ctx context.Context, id uuid.UUID) error

	// User Permission
	AssignUserToGroup(ctx context.Context, dto *AssignUserGroupDto) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
}

type service struct {
	permissionRepo Repository
	userRepo       user.Repository
	userService    user.Service
}

func NewService(permissionRepo Repository, userRepo user.Repository) Service {
	return &service{permissionRepo: permissionRepo, userRepo: userRepo, userService: user.NewService(userRepo)}
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
	if len(permissions) == 0 {
		return nil, ErrPermissionNotFound
	}
	return permissions, nil
}

// Permission Implementations
func (s *service) GetAllPermissions(ctx context.Context, search string) ([]PermissionResponseDto, int64, error) {
	return s.permissionRepo.FindAllPermissions(ctx, search)
}

func (s *service) GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error) {
	return s.permissionRepo.FindPermissionByID(ctx, id)
}

// Permission Group Implementations
func (s *service) CreatePermissionGroup(ctx context.Context, dto *CreatePermissionGroupDto) error {
	var permissions []Permission
	var assignedUsers []user.User
	
	if len(dto.Permissions) > 0 {
		ids, parseErr := utils.ParseStringsToUUIDs(dto.Permissions)
		if parseErr != nil {
			return parseErr
		}
		var checkExistingPermissionsErr error
		permissions, checkExistingPermissionsErr = s.CheckExistingPermissions(ctx, ids)
		if checkExistingPermissionsErr != nil {
			return checkExistingPermissionsErr
		}
	}

	if len(dto.Users) > 0 {
		ids, parseErr := utils.ParseStringsToUUIDs(dto.Users)
		if parseErr != nil {
			return parseErr
		}
		var err error
		assignedUsers, err = user.NewService(s.userRepo).CheckExistingUsers(ctx, ids)
		if err != nil {
			return err
		}
	}

	return s.permissionRepo.WithTransaction(ctx, func(txRepo Repository) error {
		pg := &PermissionGroup{
			Name:          dto.Name,
			Description:   dto.Description,
			Permissions:   permissions,
			AssignedUsers: assignedUsers,
		}

		err := txRepo.CreatePermissionGroup(ctx, pg)
		if err != nil {
			return err
		} 
		
		return nil 
	})
}

func (s *service) GetAllPermissionGroups(ctx context.Context, name string, assignedUser string, page, pageSize int) ([]PermissionGroup, int64, error) {
	offset := (page - 1) * pageSize
	return s.permissionRepo.FindAllPermissionGroups(ctx, name, assignedUser, offset, pageSize)
}

func (s *service) GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error) {
	permissionGroup, err := s.permissionRepo.FindPermissionGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if permissionGroup == nil {
		return nil, ErrPermissionGroupNotFound
	}
	return permissionGroup, nil
}

func (s *service) UpdatePermissionGroup(ctx context.Context, id uuid.UUID, dto *UpdatePermissionGroupDto) error {
	var permissions []Permission
	var assignedUsers []user.User
	
	if len(dto.Permissions) == 0 || dto.Permissions == nil {
		return ErrNotNullPermissions
	} else {
		ids, parseErr := utils.ParseStringsToUUIDs(dto.Permissions)
		if parseErr != nil {
			return parseErr
		}
		var checkExistingPermissionsErr error
		permissions, checkExistingPermissionsErr = s.CheckExistingPermissions(ctx, ids)
		if checkExistingPermissionsErr != nil {
			return checkExistingPermissionsErr
		}
	}

	if len(dto.Users) > 0 {
		ids, parseErr := utils.ParseStringsToUUIDs(dto.Users)
		if parseErr != nil {
			return parseErr
		}
		var err error
		assignedUsers, err = s.userService.CheckExistingUsers(ctx, ids)
		if err != nil {
			return err
		}
	}

	pg := &PermissionGroup{
		Name:        dto.Name,
		Description: dto.Description,
		AssignedUsers: assignedUsers,
		Permissions: permissions,
	}
	return s.permissionRepo.UpdatePermissionGroup(ctx, id, pg)
}

func (s *service) DeletePermissionGroup(ctx context.Context, id uuid.UUID) error {
	return s.permissionRepo.DeletePermissionGroup(ctx, id)
}

// User Permission Implementations

func (s *service) AssignUserToGroup(ctx context.Context, dto *AssignUserGroupDto) error {
	uid, err := uuid.Parse(dto.UserID)
	if err != nil {
		return err
	}
	gid, err := uuid.Parse(dto.PermissionGroupID)
	if err != nil {
		return err
	}
	return s.permissionRepo.AssignUserToGroup(ctx, uid, gid)
}

func (s *service) RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	return s.permissionRepo.RemoveUserFromGroup(ctx, userID, groupID)
}
