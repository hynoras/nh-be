package auth

import (
	"nh-be/internal/features/user"
	"time"

	"github.com/google/uuid"
)

func MapUserDtoToLoginResponse(user user.User, permissions []string) UserResponseDto {
	return UserResponseDto{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Permissions: permissions,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func MapCreateDtoToVerificationToken(userId uuid.UUID, codeHash string, tokenType VerificationTokenType) *VerificationToken {
	return &VerificationToken{
		UserID:    userId,
		Type:      tokenType,
		CodeHash:  codeHash,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(15 * time.Minute),
	}
}

func MapSignUpDtoToUser(username, email, password string) user.User {
	return user.User{
		Username: username,
		Email:    email,
		Password: &password,
	}
}

func MapVerificationTokenToCreatedTokenDto(token VerificationToken, rawToken string) CreatedTokenDto {
	return CreatedTokenDto{
		UserID:    token.UserID,
		Type:      token.Type,
		Token:     rawToken,
		CreatedAt: token.CreatedAt,
		ExpireAt:  token.ExpireAt,
	}
}
