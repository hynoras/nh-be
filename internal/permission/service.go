package permission

import (
	"context"
	"errors"
	"log"

	"nh-be/internal/user"

	"github.com/google/uuid"
)

type Service interface {
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
}

func NewService(permissionRepo Repository, userRepo user.Repository) Service {
	return &service{permissionRepo: permissionRepo, userRepo: userRepo}
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
	if len(dto.Permissions) > 0 {
		var ids []uuid.UUID
		for _, idStr := range dto.Permissions {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}
		var err error
		permissions, err = s.permissionRepo.FindPermissionsByIDs(ctx, ids)
		if err != nil {
			return err
		}
		if len(permissions) == 0 {
			return errors.New("permissions not found")
		}
	}

	var assignedUsers []user.User
	if len(dto.Users) > 0 {
		var ids []uuid.UUID
		for _, idStr := range dto.Users {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}
		var err error
		assignedUsers, err = s.userRepo.FindByIDs(ctx, ids)
		if err != nil {
			return err
		}
		if len(assignedUsers) == 0 {
			return errors.New("assigned users not found")
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
			log.Println("Failed to create permission group:", err)
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
	return s.permissionRepo.FindPermissionGroupByID(ctx, id)
}

func (s *service) UpdatePermissionGroup(ctx context.Context, id uuid.UUID, dto *UpdatePermissionGroupDto) error {
	var permissions []Permission
	if dto.PermissionIDs != nil { // Check if PermissionIDs is provided to update permissions
		var ids []uuid.UUID
		for _, idStr := range dto.PermissionIDs {
			uid, err := uuid.Parse(idStr)
			if err != nil {
				return err
			}
			ids = append(ids, uid)
		}
		var err error
		permissions, err = s.permissionRepo.FindPermissionsByIDs(ctx, ids)
		if err != nil {
			return err
		}
	}

	pg := &PermissionGroup{
		Name:        dto.Name,
		Description: dto.Description,
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
