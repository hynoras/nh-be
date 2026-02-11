package auth

import (
	"context"
	"errors"

	"nh-be/utils"

	"nh-be/internal/permission"
	"nh-be/internal/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	SignUp(ctx context.Context, req SignUpDto) (*user.User, error)
	Login(ctx context.Context, email, password string) (*user.User, []string, error)
	Logout(c *gin.Context) error
	ChangePassword(ctx context.Context, id uuid.UUID, changePasswordDto ChangePasswordDto) error
}

type service struct {
	userRepo          user.Repository
	permissionService permission.Service
	authPublisher     AuthPublisher
}

func NewService(userRepo user.Repository, permissionService permission.Service, authPublisher AuthPublisher) Service {
	return &service{userRepo: userRepo, permissionService: permissionService, authPublisher: authPublisher}
}

func (s *service) SignUp(ctx context.Context, req SignUpDto) (*user.User, error) {
	u, err := s.userRepo.FindByEmail(ctx, req.Email)

	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}
	if u != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &user.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	// repoErr := s.userRepo.Create(ctx, user)
	// if repoErr != nil {
	// 	return nil, repoErr
	// }

	// Publish event to RabbitMQ
	err = s.authPublisher.PublishSendVerificationEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, email, password string) (*user.User, []string, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, []string{}, err
	}
	if u == nil {
		return nil, []string{}, errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, u.Password) {
		return nil, []string{}, errors.New("invalid credentials")
	}

	permissions, err := s.permissionService.GetUserPermissionCodeNames(ctx, u.ID)
	if err != nil {
		return nil, []string{}, err
	}

	return u, permissions, nil
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
		return errors.New("user not found")
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
