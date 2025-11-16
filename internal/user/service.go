package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
)

type Service interface {
  GetAllUsers(ctx context.Context, search string, role string, offset int, limit int) ([]User, int64, error)
  GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
  CreateUser(ctx context.Context, user *User) error
  UpdateUser(ctx context.Context, id uuid.UUID, user *User) error
  DeleteUsers(ctx context.Context, ids []uuid.UUID) error
}

type service struct {
  userRepo Repository
}

func NewService(userRepo Repository) Service {
  return &service{userRepo: userRepo}
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

func (s *service) CreateUser(ctx context.Context, user *User) error {
	// Check for duplicate username
	existingUser, err := s.userRepo.FindByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return ErrDuplicateUsername
	}

	// Check for duplicate email
	existingUser, err = s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return ErrDuplicateEmail
	}

	userDto := &User{
		Username: user.Username,
		Email: user.Email,
		Password: user.Password,
		Role: user.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err = s.userRepo.Create(ctx, userDto)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, user *User) error {
	userDto := &User{
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
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