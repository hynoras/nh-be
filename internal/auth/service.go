package auth

import (
	"context"
	"errors"

	"nh-be/utils"

	"nh-be/internal/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
  Login(ctx context.Context, email, password string) (*user.User, error)
  Logout(c *gin.Context) error
  ChangePassword(ctx context.Context, id uuid.UUID, changePasswordDto ChangePasswordDto) error
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

func (s *service) ChangePassword(ctx context.Context, id uuid.UUID, changePasswordDto ChangePasswordDto) error {
  oldPassword, err := s.userRepo.FindPasswordById(ctx, id)
  if err != nil {
    return err
  }
  if oldPassword == nil {
    return  errors.New("user not found")
  }

  if changePasswordDto.NewPassword != changePasswordDto.ConfirmPassword {
    return errors.New("new password and confirm password do not match")
  }

  newHashedPassword, err := utils.HashPassword(changePasswordDto.NewPassword)
  if err != nil {
    return err
  }

  if utils.CheckPasswordHash(changePasswordDto.NewPassword, *oldPassword) {
    return errors.New("new password is the same as the old password")
  }

  err = s.userRepo.Update(ctx, id, &user.User{
    Password: newHashedPassword,
  })

  if err != nil {
    return err
  }

  return nil
}