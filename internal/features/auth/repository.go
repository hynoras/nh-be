package auth

import (
	"context"
	"nh-be/internal/constant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindVerificationTokenByUserId(ctx context.Context, userId uuid.UUID) (*VerificationToken, error)
	FindVerificationTokenByCodeHash(ctx context.Context, codeHash string) (*VerificationToken, error)
	CreateVerificationToken(ctx context.Context, token *VerificationToken) (VerificationToken, error)

	DeleteVerificationToken(ctx context.Context, token *VerificationToken) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindVerificationTokenByUserId(ctx context.Context, userId uuid.UUID) (*VerificationToken, error) {
	var token VerificationToken
	err := r.db.WithContext(ctx).Model(&VerificationToken{}).Where("user_id = ?", userId).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constant.ErrVerificationTokenNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *repository) FindVerificationTokenByCodeHash(ctx context.Context, codeHash string) (*VerificationToken, error) {
	var token VerificationToken
	err := r.db.WithContext(ctx).Model(&VerificationToken{}).Where("code_hash = ?", codeHash).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constant.ErrVerificationTokenNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *repository) CreateVerificationToken(ctx context.Context, token *VerificationToken) (VerificationToken, error) {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return VerificationToken{}, err
	}
	return *token, nil
}

func (r *repository) DeleteVerificationToken(ctx context.Context, token *VerificationToken) error {
	return r.db.WithContext(ctx).Delete(token).Error
}
