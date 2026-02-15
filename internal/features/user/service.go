package user

import (
	"context"
	"errors"
	"nh-be/constant"
	"nh-be/internal/features/permission"
	"nh-be/internal/utils/ctxutil"
	"slices"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	CheckExistingUser(ctx context.Context, userId uuid.UUID) (*User, error)
	CheckExistingUsers(ctx context.Context, userIds []uuid.UUID) ([]User, error)
	GetAllUsers(ctx context.Context, search string, page, pageSize int) ([]User, int64, error)
	GetUserById(ctx context.Context, id uuid.UUID, isMe bool) (*User, []string, error)
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

func (s *service) CreateUser(ctx context.Context, userInput *UserInput) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageUser) {
		return ErrForbidCreateUser
	}

	// Check for duplicate username
	existingUser, err := s.userRepo.FindByUsername(ctx, userInput.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return ErrDuplicateUsername
	}

	// Check for duplicate email
	existingUser, err = s.userRepo.FindByEmail(ctx, userInput.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (s *service) GetAllUsers(ctx context.Context, search string, page, pageSize int) ([]User, int64, error) {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	if !slices.Contains(userPerm, constant.ViewUser) && !slices.Contains(userPerm, constant.ManageUser) {
		return nil, 0, ErrForbidViewUsers
	}

	users, length, err := s.userRepo.FindAll(ctx, search, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return users, length, nil
}

func (s *service) GetUserById(ctx context.Context, id uuid.UUID, isMe bool) (*User, []string, error) {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, []string{}, err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return nil, []string{}, err
	}

	if isMe == false && !slices.Contains(userPerm, constant.ViewUser) && !slices.Contains(userPerm, constant.ManageUser) {
		return nil, []string{}, ErrForbidViewUser
	}

	var user *User
	var userErr error
	var permissionCodes []string
	var permCodeErr error

	user, userErr = s.userRepo.FindByID(ctx, id)
	if userErr != nil {
		return nil, []string{}, userErr
	}

	if isMe == true {
		permissionCodes, permCodeErr = s.permissionService.GetUserPermissionCodeNames(ctx, id)
		if permCodeErr != nil {
			return nil, []string{}, permCodeErr
		}
	}
	return user, permissionCodes, nil
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, userInput *UserInput) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageUser) {
		return ErrForbidUpdateUser
	}

	var permissionGroups []permission.PermissionGroup
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

	return nil
}

func (s *service) DeleteUsers(ctx context.Context, ids []uuid.UUID) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerm, err := s.permissionService.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	if !slices.Contains(userPerm, constant.ManageUser) {
		return ErrForbidDeleteUser
	}

	for _, id := range ids {
		err := s.userRepo.Delete(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}
