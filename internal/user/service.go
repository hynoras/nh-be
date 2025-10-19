package user

import (
	"context"
)

type Service interface {
  GetAllUsers(ctx context.Context) ([]User, error)
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