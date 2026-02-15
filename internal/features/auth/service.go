package auth

import (
	"context"
	"errors"
	"time"

	"nh-be/internal/email"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/internal/infra"
	"nh-be/internal/utils/crypto"
	"nh-be/internal/utils/stringutil"

	"github.com/google/uuid"
)

type Service interface {
	SignUp(ctx context.Context, req SignUpDto) error
	VerifyEmail(ctx context.Context, token string) (string, error)
	Login(ctx context.Context, email, password string) (*user.User, []string, string, error)
	Logout(ctx context.Context, sessionId string) error
	ChangePassword(ctx context.Context, id uuid.UUID, changePasswordDto ChangePasswordDto) error
	CreateVerificationToken(ctx context.Context, userId uuid.UUID, tokenType VerificationTokenType) (CreatedTokenDto, error)
}

type service struct {
	sessionStore      infra.SessionStore
	authRepo          Repository
	userRepo          user.Repository
	permissionService permission.Service
	emailPublisher    email.EmailPublisher
}

func NewService(
	sessionStore infra.SessionStore,
	authRepo Repository,
	userRepo user.Repository,
	permissionService permission.Service,
	emailPublisher email.EmailPublisher,
) Service {
	return &service{
		sessionStore:      sessionStore,
		authRepo:          authRepo,
		userRepo:          userRepo,
		permissionService: permissionService,
		emailPublisher:    emailPublisher,
	}
}

func (s *service) publishSendVerificationEmail(ctx context.Context, userId uuid.UUID, toEmail string, tokenType VerificationTokenType) error {
	createdToken, tokenErr := s.CreateVerificationToken(ctx, userId, tokenType)
	if tokenErr != nil {
		return tokenErr
	}

	sendVerificationEmailDto := email.MapToSendVerificationEmailDto(toEmail, createdToken.Token)

	err := s.emailPublisher.SendVerificationEmail(ctx, sendVerificationEmailDto)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateVerificationToken(
	ctx context.Context,
	userId uuid.UUID,
	tokenType VerificationTokenType,
) (CreatedTokenDto, error) {
	existingToken, findErr := s.authRepo.FindVerificationTokenByUserId(userId)
	if findErr != nil && !errors.Is(findErr, ErrVerificationTokenNotFound) {
		return CreatedTokenDto{}, findErr
	}

	if existingToken != nil {
		deleteErr := s.authRepo.DeleteVerificationToken(existingToken)
		if deleteErr != nil {
			return CreatedTokenDto{}, deleteErr
		}
	}

	generatedToken, genErr := crypto.GenerateToken()
	if genErr != nil {
		return CreatedTokenDto{}, genErr
	}

	hashedToken := crypto.HashToken(generatedToken)
	verificationToken := MapCreateDtoToVerificationToken(userId, hashedToken, tokenType)

	createdToken, createErr := s.authRepo.CreateVerificationToken(verificationToken)
	if createErr != nil {
		return CreatedTokenDto{}, createErr
	}
	mapToken := MapVerificationTokenToCreatedTokenDto(createdToken, generatedToken)

	return mapToken, nil
}

func (s *service) SignUp(ctx context.Context, req SignUpDto) error {
	u, err := s.userRepo.FindByEmail(ctx, req.Email)

	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return err
	}

	if u != nil {
		err := s.publishSendVerificationEmail(ctx, u.ID, u.Email, VerifyEmail)
		if err != nil {
			return err
		}
	} else {
		hashedPassword, err := crypto.HashPassword(req.Password)
		if err != nil {
			return err
		}

		username := stringutil.ExtractUsernameFromEmail(req.Email)
		userForm := MapSignUpDtoToUser(username, req.Email, hashedPassword)

		createdUser, createErr := s.userRepo.Create(ctx, &userForm)
		if createErr != nil {
			return createErr
		}

		mapUser := user.MapUserToCreatedUser(createdUser)

		sendEmailErr := s.publishSendVerificationEmail(ctx, mapUser.ID, mapUser.Email, VerifyEmail)
		if sendEmailErr != nil {
			return sendEmailErr
		}
	}

	return nil
}

func (s *service) VerifyEmail(ctx context.Context, token string) (string, error) {
	hashedToken := crypto.HashToken(token)
	existingToken, findErr := s.authRepo.FindVerificationTokenByCodeHash(hashedToken)
	if findErr != nil {
		return "", findErr
	}
	if existingToken.Type != VerifyEmail {
		return "", ErrInvalidVerificationToken
	}
	if existingToken.ExpireAt.Before(time.Now()) {
		return "", ErrVerificationTokenExpired
	}

	updateErr := s.userRepo.Update(ctx, existingToken.UserID, &user.User{
		IsVerified: true,
	})
	if updateErr != nil {
		return "", updateErr
	}

	deleteErr := s.authRepo.DeleteVerificationToken(existingToken)
	if deleteErr != nil {
		return "", deleteErr
	}

	sessionId, genErr := crypto.GenerateToken()
	if genErr != nil {
		return "", genErr
	}

	sessionErr := s.sessionStore.CreateUserSession(ctx, sessionId, existingToken.UserID.String())
	if sessionErr != nil {
		return "", sessionErr
	}
	return sessionId, nil
}

func (s *service) Login(ctx context.Context, email, password string) (*user.User, []string, string, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, []string{}, "", err
	}
	if u == nil {
		return nil, []string{}, "", errors.New("invalid credentials")
	}
	if !crypto.CheckPasswordHash(password, u.Password) {
		return nil, []string{}, "", errors.New("invalid credentials")
	}

	permissions, err := s.permissionService.GetUserPermissionCodeNames(ctx, u.ID)
	if err != nil {
		return nil, []string{}, "", err
	}

	sessionId, genErr := crypto.GenerateToken()
	if genErr != nil {
		return nil, []string{}, "", genErr
	}

	sessionErr := s.sessionStore.CreateUserSession(ctx, sessionId, u.ID.String())
	if sessionErr != nil {
		return nil, []string{}, "", sessionErr
	}

	return u, permissions, sessionId, nil
}

func (s *service) Logout(ctx context.Context, sessionId string) error {
	err := s.sessionStore.DeleteUserSession(ctx, sessionId)
	if err != nil {
		return err
	}

	return nil
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

	newHashedPassword, err := crypto.HashPassword(changePasswordDto.NewPassword)
	if err != nil {
		return err
	}

	if crypto.CheckPasswordHash(changePasswordDto.NewPassword, *oldPassword) {
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
