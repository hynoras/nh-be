package user

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
  GetAllUsers(ctx context.Context) ([]User, error)
  GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
}

type service struct {
  userRepo Repository
}

func NewService(userRepo Repository) Service {
  return &service{userRepo: userRepo}
}

func (s *service) GetAllUsers(ctx context.Context) ([]User, error) {
	users, err := s.userRepo.FindAll(ctx)
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