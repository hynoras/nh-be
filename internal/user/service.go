package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
  GetAllUsers(ctx context.Context, search string) ([]User, error)
  GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
  UpdateUser(ctx context.Context, id uuid.UUID, user *User) error
}

type service struct {
  userRepo Repository
}

func NewService(userRepo Repository) Service {
  return &service{userRepo: userRepo}
}

func (s *service) GetAllUsers(ctx context.Context, search string) ([]User, error) {
	users, err := s.userRepo.FindAll(ctx, search)
    if err != nil {
    	return nil, err
  	}
  return users, nil
}

func (s *service) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
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