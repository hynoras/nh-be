package auth

import (
	"nh-be/internal/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindVerificationTokenByUserId(userId uuid.UUID) (*VerificationToken, error)
	FindVerificationTokenByCodeHash(codeHash string) (*VerificationToken, error)
	CreateVerificationToken(token *VerificationToken) (VerificationToken, error)

	DeleteVerificationToken(token *VerificationToken) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindVerificationTokenByUserId(userId uuid.UUID) (*VerificationToken, error) {
	var token VerificationToken
	err := r.db.Model(&VerificationToken{}).Where("user_id = ?", userId).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *repository) FindVerificationTokenByCodeHash(codeHash string) (*VerificationToken, error) {
	var token VerificationToken
	err := r.db.Model(&VerificationToken{}).Where("code_hash = ?", codeHash).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrVerficationTokenNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *repository) CreateVerificationToken(token *VerificationToken) (VerificationToken, error) {
	err := r.db.Create(token).Error
	if err != nil {
		return VerificationToken{}, err
	}
	return *token, nil
}

func (r *repository) DeleteVerificationToken(token *VerificationToken) error {
	return r.db.Delete(token).Error
}
