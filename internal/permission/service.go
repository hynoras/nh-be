package permission

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type Service interface {
	// Permission
	GetAllPermissions(ctx context.Context, search string) ([]Permission, int64, error)
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error)

	// Permission Group
	CreatePermissionGroup(ctx context.Context, dto *CreatePermissionGroupDto) error
	GetAllPermissionGroups(ctx context.Context, name string, assignedUser uuid.UUID, page, pageSize int) ([]PermissionGroup, int64, error)
	GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error)
	UpdatePermissionGroup(ctx context.Context, id uuid.UUID, dto *UpdatePermissionGroupDto) error
	DeletePermissionGroup(ctx context.Context, id uuid.UUID) error

	// User Permission
	AssignUserToGroup(ctx context.Context, dto *AssignUserGroupDto) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Permission Implementations

func (s *service) GetAllPermissions(ctx context.Context, search string) ([]Permission, int64, error) {
	log.Println("search in service", search)
	return s.repo.FindAllPermissions(ctx, search)
}

func (s *service) GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error) {
	return s.repo.FindPermissionByID(ctx, id)
}

// Permission Group Implementations

func (s *service) CreatePermissionGroup(ctx context.Context, dto *CreatePermissionGroupDto) error {
	var permissions []Permission
	if len(dto.PermissionIDs) > 0 {
		var ids []uuid.UUID
		for _, idStr := range dto.PermissionIDs {
			id, err := uuid.Parse(idStr)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}
		var err error
		permissions, err = s.repo.FindPermissionsByIDs(ctx, ids)
		if err != nil {
			return err
		}
	}

	pg := &PermissionGroup{
		Name:        dto.Name,
		Description: dto.Description,
		Permissions: permissions,
	}
	return s.repo.CreatePermissionGroup(ctx, pg)
}

func (s *service) GetAllPermissionGroups(ctx context.Context, name string, assignedUser uuid.UUID, page, pageSize int) ([]PermissionGroup, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.FindAllPermissionGroups(ctx, name, assignedUser, offset, pageSize)
}

func (s *service) GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (*PermissionGroup, error) {
	return s.repo.FindPermissionGroupByID(ctx, id)
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
		permissions, err = s.repo.FindPermissionsByIDs(ctx, ids)
		if err != nil {
			return err
		}
	}

	pg := &PermissionGroup{
		Name:        dto.Name,
		Description: dto.Description,
		Permissions: permissions,
	}
	return s.repo.UpdatePermissionGroup(ctx, id, pg)
}

func (s *service) DeletePermissionGroup(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePermissionGroup(ctx, id)
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
	return s.repo.AssignUserToGroup(ctx, uid, gid)
}

func (s *service) RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	return s.repo.RemoveUserFromGroup(ctx, userID, groupID)
}
