package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"nh-be/internal/app"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/internal/platform/email"
	"nh-be/internal/platform/session"
	"nh-be/internal/utils/crypto"
	"nh-be/internal/utils/stringutil"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service interface {
	SignUp(ctx context.Context, req SignUpDto) error
	VerifyEmail(ctx context.Context, token string) (string, error)
	Login(ctx context.Context, email, password string) (*UserResponseDto, string, string, error)
	GenerateProviderLoginURL(ctx context.Context, provider string) (string, string, string, error)
	ProviderCallback(ctx context.Context, provider string, code string, verifier string) (*UserResponseDto, string, string, error)
	Logout(ctx context.Context, sessionId string) error
	ChangePassword(ctx context.Context, id uuid.UUID, changePasswordDto ChangePasswordDto) error
	CreateVerificationToken(ctx context.Context, userId uuid.UUID, tokenType VerificationTokenType) (CreatedTokenDto, error)
}

type service struct {
	sessionStore      session.SessionStore
	authRepo          Repository
	userRepo          user.Repository
	permissionService permission.Service
	emailPublisher    email.EmailPublisher
	oauthProviders    map[string]*app.OAuthProviderConfig
}

func NewService(
	sessionStore session.SessionStore,
	authRepo Repository,
	userRepo user.Repository,
	permissionService permission.Service,
	emailPublisher email.EmailPublisher,
	oauthProviders map[string]*app.OAuthProviderConfig,
) Service {
	return &service{
		sessionStore:      sessionStore,
		authRepo:          authRepo,
		userRepo:          userRepo,
		permissionService: permissionService,
		emailPublisher:    emailPublisher,
		oauthProviders:    oauthProviders,
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
	existingToken, findErr := s.authRepo.FindVerificationTokenByUserId(ctx, userId)
	if findErr != nil && !errors.Is(findErr, ErrVerificationTokenNotFound) {
		return CreatedTokenDto{}, findErr
	}

	if existingToken != nil {
		deleteErr := s.authRepo.DeleteVerificationToken(ctx, existingToken)
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

	createdToken, createErr := s.authRepo.CreateVerificationToken(ctx, verificationToken)
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
	existingToken, findErr := s.authRepo.FindVerificationTokenByCodeHash(ctx, hashedToken)
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

	deleteErr := s.authRepo.DeleteVerificationToken(ctx, existingToken)
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

func (s *service) Login(ctx context.Context, email, password string) (*UserResponseDto, string, string, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}
	if u == nil {
		return nil, "", "", ErrInvalidCredentials
	}
	if u.Password == nil || !crypto.CheckPasswordHash(password, *u.Password) {
		return nil, "", "", ErrInvalidCredentials
	}

	permissions, err := s.permissionService.GetUserPermissionCodeNames(ctx, u.ID)
	if err != nil {
		return nil, "", "", err
	}

	sessionId, sessionIdErr := crypto.GenerateToken()
	if sessionIdErr != nil {
		return nil, "", "", sessionIdErr
	}

	sessionErr := s.sessionStore.CreateUserSession(ctx, sessionId, u.ID.String())
	if sessionErr != nil {
		return nil, "", "", sessionErr
	}

	csrfToken, genCSRFTokenErr := crypto.GenerateToken()
	if genCSRFTokenErr != nil {
		return nil, "", "", genCSRFTokenErr
	}

	mappedUser := MapUserDtoToLoginResponse(*u, permissions)

	return &mappedUser, sessionId, csrfToken, nil
}

func (s *service) GenerateProviderLoginURL(ctx context.Context, provider string) (string, string, string, error) {
	providerCfg, ok := s.oauthProviders[provider]
	if !ok {
		return "", "", "", errors.New("unsupported oauth provider")
	}

	// 1. Generate state (for CSRF protection)
	state, stateErr := crypto.GenerateRandomString(32)
	if stateErr != nil {
		return "", "", "", stateErr
	}

	// 2. Generate PKCE verifier and challenge
	verifier := oauth2.GenerateVerifier()

	var endpoint oauth2.Endpoint
	switch provider {
	case "google":
		endpoint = google.Endpoint
	default:
		return "", "", "", errors.New("unsupported oauth provider endpoint")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     providerCfg.ClientID,
		ClientSecret: providerCfg.ClientSecret,
		RedirectURL:  providerCfg.RedirectURL,
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		Endpoint: endpoint,
	}

	url := oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)

	return url, state, verifier, nil
}

func (s *service) ProviderCallback(ctx context.Context, provider string, code string, verifier string) (*UserResponseDto, string, string, error) {
	providerCfg, ok := s.oauthProviders[provider]
	if !ok {
		return nil, "", "", errors.New("unsupported oauth provider")
	}

	var endpoint oauth2.Endpoint
	var userInfoURL string

	switch provider {
	case "google":
		endpoint = google.Endpoint
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	default:
		return nil, "", "", errors.New("unsupported oauth provider endpoint")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     providerCfg.ClientID,
		ClientSecret: providerCfg.ClientSecret,
		RedirectURL:  providerCfg.RedirectURL,
		Endpoint:     endpoint,
	}

	token, err := oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, "", "", err
	}

	client := oauthConfig.Client(ctx, token)
	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", err
	}

	var userInfo map[string]any

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, "", "", err
	}

	existingUser, findUserErr := s.userRepo.FindByEmail(ctx, userInfo["email"].(string))
	if findUserErr != nil && !errors.Is(findUserErr, user.ErrUserNotFound) {
		return nil, "", "", findUserErr
	}

	if existingUser == nil {
		username := stringutil.ExtractUsernameFromEmail(userInfo["email"].(string))

		createdUser, createUserErr := s.userRepo.Create(ctx, &user.User{
			Email:      userInfo["email"].(string),
			Username:   username,
			IsVerified: true,
		})

		if createUserErr != nil {
			return nil, "", "", createUserErr
		}
		existingUser = &createdUser
	}

	permissions, err := s.permissionService.GetUserPermissionCodeNames(ctx, existingUser.ID)
	if err != nil {
		return nil, "", "", err
	}

	sessionId, sessionIdErr := crypto.GenerateToken()
	if sessionIdErr != nil {
		return nil, "", "", sessionIdErr
	}

	sessionErr := s.sessionStore.CreateUserSession(ctx, sessionId, existingUser.ID.String())
	if sessionErr != nil {
		return nil, "", "", sessionErr
	}

	csrfToken, genCSRFTokenErr := crypto.GenerateToken()
	if genCSRFTokenErr != nil {
		return nil, "", "", genCSRFTokenErr
	}

	mappedUser := MapUserDtoToLoginResponse(*existingUser, permissions)

	return &mappedUser, sessionId, csrfToken, nil
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
		return user.ErrUserNotFound
	}

	if changePasswordDto.NewPassword != changePasswordDto.ConfirmPassword {
		return ErrNewPasswordAndConfirmPasswordDoNotMatch
	}

	newHashedPassword, err := crypto.HashPassword(changePasswordDto.NewPassword)
	if err != nil {
		return err
	}

	if crypto.CheckPasswordHash(changePasswordDto.NewPassword, *oldPassword) {
		return ErrNewPasswordIsTheSameAsOldPassword
	}

	stringifiedPassword := string(newHashedPassword)

	err = s.userRepo.Update(ctx, id, &user.User{
		Password: &stringifiedPassword,
	})

	if err != nil {
		return err
	}

	return nil
}
