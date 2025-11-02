package user

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
  Create(ctx context.Context, u *User) error
  FindAll (ctx context.Context) ([]User, error)
  FindByEmail(ctx context.Context, email string) (*User, error)
  FindByUsername(ctx context.Context, username string) (*User, error)
  FindByID(ctx context.Context, id uuid.UUID) (*User, error)
  Update(ctx context.Context, u *User) error
  Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
  db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
  return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, u *User) error {
  return r.db.WithContext(ctx).Create(u).Error
}

func (r *repository) FindAll(ctx context.Context) ([]User, error) {
  var users []User
  result := r.db.WithContext(ctx).Find(&users)
  return users, result.Error
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
  var u User
  result := r.db.WithContext(ctx).Where("email = ?", email).First(&u)
  if result.Error != nil {
    return nil, result.Error
  }
  return &u, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
  var u User
  result := r.db.WithContext(ctx).Where("username = ?", username).First(&u)
  if result.Error != nil {
    return nil, result.Error
  }
  return &u, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
  var u User
  result := r.db.WithContext(ctx).Where("id = ?", id).First(&u)
  if result.Error != nil {
    return nil, result.Error
  }
  return &u, nil
}

func (r *repository) Update(ctx context.Context, u *User) error {
  return r.db.WithContext(ctx).Save(u).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
  return r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}