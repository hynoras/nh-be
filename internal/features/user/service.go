package user

import (
	"context"
	"errors"
	"nh-be/internal/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/authutil"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CheckExistingUser(ctx context.Context, userId uuid.UUID) (*User, error)
	CheckExistingUsers(ctx context.Context, userIds []uuid.UUID) ([]User, error)
	GetAllUsers(ctx context.Context, search string, page, pageSize int) ([]UserResponseDto, int64, error)
	GetUserById(ctx context.Context, id uuid.UUID, isMe bool) (interface{}, error)
	CreateUser(ctx context.Context, userInput *UserInput) error
	UpdateUser(ctx context.Context, id uuid.UUID, userInput *UserInput) error
	DeleteUsers(ctx context.Context, ids []uuid.UUID) error
}

type service struct {
	userRepo          Repository
	permissionRepo    permission.Repository
	permissionService permission.Service
}

func NewService(
	userRepo Repository,
	permissionRepo permission.Repository,
	permissionService permission.Service,
) Service {
	return &service{
		userRepo:          userRepo,
		permissionRepo:    permissionRepo,
		permissionService: permissionService,
	}
}

func (s *service) CheckExistingUser(ctx context.Context, userId uuid.UUID) (*User, error) {
	assignedUser, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	if assignedUser == nil {
		return nil, ErrAssignedUserNotFound
	}
	return assignedUser, nil
}

func (s *service) CheckExistingUsers(ctx context.Context, userIds []uuid.UUID) ([]User, error) {
	users, err := s.userRepo.FindByIDs(ctx, userIds)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUsersNotFound
	}
	return users, nil

}

func (s *service) CreateUser(ctx context.Context, userInput *UserInput) error {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidCreateUser, constant.ManageUser); err != nil {
		return err
	}

	// Check for duplicate username
	existingUser, err := s.userRepo.FindByUsername(ctx, userInput.Username)
	if err != nil && !errors.Is(err, constant.ErrUserNotFound) {
		return err
	}
	if existingUser != nil {
		return ErrDuplicateUsername
	}

	// Check for duplicate email
	existingUser, err = s.userRepo.FindByEmail(ctx, userInput.Email)
	if err != nil && !errors.Is(err, constant.ErrUserNotFound) {
		return err
	}
	if existingUser != nil {
		return ErrDuplicateEmail
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	// Parse and validate permission groups exist if provided
	var permissionGroups []permission.PermissionGroup
	if len(userInput.Permissions) > 0 {
		permissionGroups, err = s.permissionService.GetPermissionGroupsByIDs(ctx, userInput.Permissions)
		if err != nil {
			return err
		}
	}

	return s.userRepo.WithTransaction(ctx, func(txRepo Repository) error {
		userToCreate := &User{
			Username:                 userInput.Username,
			Email:                    userInput.Email,
			Password:                 string(hashedPassword),
			AssignedPermissionGroups: permissionGroups,
			CreatedAt:                time.Now(),
		}
		_, err := txRepo.Create(ctx, userToCreate)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) GetAllUsers(ctx context.Context, search string, page, pageSize int) ([]UserResponseDto, int64, error) {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidViewUsers, constant.ViewUser, constant.ManageUser); err != nil {
		return nil, 0, err
	}

	users, length, err := s.userRepo.FindAll(ctx, search, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	mappedUser := MapUsersToDto(users)
	return mappedUser, length, nil
}

func (s *service) GetUserById(ctx context.Context, id uuid.UUID, isMe bool) (interface{}, error) {
	if !isMe {
		if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidViewUser, constant.ViewUser, constant.ManageUser); err != nil {
			return nil, err
		}
	}

	var user *User
	var userErr error
	var permissionCodes []string
	var permCodeErr error
	var mapperUser interface{}

	user, userErr = s.userRepo.FindByID(ctx, id)
	if userErr != nil {
		return nil, userErr
	}

	if isMe {
		permissionCodes, permCodeErr = s.permissionService.GetUserPermissionCodeNames(ctx, id)
		if permCodeErr != nil {
			return nil, permCodeErr
		}
		mapperUser = MapUserToMeDto(*user, permissionCodes)
	} else {
		mapperUser = MapUserToDto(*user)
	}
	return mapperUser, nil
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, userInput *UserInput) error {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidUpdateUser, constant.ManageUser); err != nil {
		return err
	}

	var permissionGroups []permission.PermissionGroup
	var err error
	if len(userInput.Permissions) > 0 {
		permissionGroups, err = s.permissionService.GetPermissionGroupsByIDs(ctx, userInput.Permissions)
		if err != nil {
			return err
		}
	}

	userToUpdate := &User{
		Username:                 userInput.Username,
		Email:                    userInput.Email,
		AssignedPermissionGroups: permissionGroups,
		UpdatedAt:                time.Now(),
	}

	err = s.userRepo.Update(ctx, id, userToUpdate)
	if err != nil {
		return err
	}

	// Invalidate permission cache for the updated user
	_ = s.permissionService.InvalidateUserPermissionCache(ctx, id)

	return nil
}

func (s *service) DeleteUsers(ctx context.Context, ids []uuid.UUID) error {
	if err := authutil.RequirePermission(ctx, s.permissionService, ErrForbidDeleteUser, constant.ManageUser); err != nil {
		return err
	}

	return s.userRepo.WithTransaction(ctx, func(txRepo Repository) error {
		for _, id := range ids {
			err := txRepo.Delete(ctx, id)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
