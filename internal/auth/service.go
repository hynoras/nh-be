package auth

import (
	"context"
	"errors"

	"nh-be/utils"

	"nh-be/internal/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Service interface {
  Login(ctx context.Context, email, password string) (*user.User, error)
  Logout(c *gin.Context) error
}

type service struct {
  userRepo user.Repository
}

func NewService(userRepo user.Repository) Service {
  return &service{userRepo: userRepo}
}

func (s *service) Login(ctx context.Context, email, password string) (*user.User, error) {
  u, err := s.userRepo.FindByEmail(ctx, email)
  if err != nil {
    return nil, err
  }
  if u == nil {
    return nil, errors.New("invalid credentials")
  }
  if !utils.CheckPasswordHash(password, u.Password) {
    return nil, errors.New("invalid credentials")
  }
  return u, nil
}

func (s *service) Logout(c *gin.Context) error {
  sess := sessions.Default(c)
  sess.Clear()
  return sess.Save()
}