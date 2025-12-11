package user

import (
	"context"
	"errors"
	"nh-be/internal/permission"
	"nh-be/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
)

type Service interface {
	CheckExistingUser(ctx context.Context, userId uuid.UUID) (*User, error)
	CheckExistingUsers(ctx context.Context, userIds []uuid.UUID) ([]User, error)
	GetAllUsers(ctx context.Context, search string, role string, offset int, limit int) ([]User, int64, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
	CreateUser(ctx context.Context, user *CreateUserDto) error
	UpdateUser(ctx context.Context, id uuid.UUID, user *User) error
	DeleteUsers(ctx context.Context, ids []uuid.UUID) error
}

type service struct {
  userRepo Repository
  permissionRepo permission.Repository
  permissionService permission.Service
}

func NewService(userRepo Repository, permissionRepo permission.Repository, permissionService permission.Service) Service {
  return &service{userRepo: userRepo, permissionRepo: permissionRepo, permissionService: permissionService}
}

func (s *service) CheckExistingUser(ctx context.Context, userId uuid.UUID) (*User, error) {
		assignedUser, err := s.userRepo.FindByID(ctx, userId)
		if err != nil {
			return nil, err
		}
		if assignedUser == nil {
			return nil, errors.New("assigned users not found")
		}
	return assignedUser, nil
}

func (s *service) CheckExistingUsers(ctx context.Context, userIds []uuid.UUID) ([]User, error) {
	users, err := s.userRepo.FindByIDs(ctx, userIds)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("users not found")
	}
	return users, nil

}

func (s *service) GetAllUsers(ctx context.Context, search string, role string, offset int, limit int) ([]User, int64, error) {
	users, length, err := s.userRepo.FindAll(ctx, search, role, offset, limit)
    if err != nil {
    	return nil, 0, err
  	}
  return users, length, nil
}

func (s *service) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) CreateUser(ctx context.Context, dto *CreateUserDto) error {
	// Check for duplicate username
	existingUser, err := s.userRepo.FindByUsername(ctx, dto.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return ErrDuplicateUsername
	}

	// Check for duplicate email
	existingUser, err = s.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return ErrDuplicateEmail
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}
	
	// Parse and validate permission groups exist if provided
	var permissionGroups []permission.PermissionGroup
	if len(dto.Permissions) > 0 {
		ids, parseErr := utils.ParseStringsToUUIDs(dto.Permissions)
		if parseErr != nil {
			return errors.New("invalid permission group ID: " + parseErr.Error())
		}
		permissionGroups, err = s.permissionService.CheckExistingPermissionGroups(ctx, ids)
		if err != nil {
			return err
		}
	}
	
	return s.userRepo.WithTransaction(ctx, func(txRepo Repository) error {
		userDto := &User{
			Username: dto.Username,
			Email: dto.Email,
			Password: string(hashedPassword),
			Role: dto.Role,
			AssignedPermissionGroups: permissionGroups,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := txRepo.Create(ctx, userDto)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, user *User) error {
	// Validate permission groups exist if provided
	if len(user.AssignedPermissionGroups) > 0 {
		groupIDs := make([]uuid.UUID, len(user.AssignedPermissionGroups))
		for i, group := range user.AssignedPermissionGroups {
			groupIDs[i] = group.ID
		}
		validGroups, err := s.userRepo.FindPermissionGroupsByIDs(ctx, groupIDs)
		if err != nil {
			return err
		}
		if len(validGroups) != len(groupIDs) {
			return errors.New("one or more permission groups not found")
		}
	}
	
	userDto := &User{
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
		AssignedPermissionGroups: user.AssignedPermissionGroups,
		UpdatedAt: time.Now(),
	}
	err := s.userRepo.Update(ctx, id, userDto) 
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteUsers(ctx context.Context, ids []uuid.UUID) error {
	for _, id := range ids {
		err := s.userRepo.Delete(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}