package user

import (
	"context"
	"nh-be/utils"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
  Create(ctx context.Context, u *User) error
  FindAll (ctx context.Context, search string, role string, offset int, limit int) ([]UserResponseDto, int64, error)
  FindByEmail(ctx context.Context, email string) (*User, error)
  FindByUsername(ctx context.Context, username string) (*User, error)
  FindByID(ctx context.Context, id uuid.UUID) (*User, error)
  FindPasswordById(ctx context.Context, id uuid.UUID) (*string, error)
  Update(ctx context.Context, id uuid.UUID, u *User) error
  Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
  db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
  return &repository{db: db}
}

func (r *repository) FindAll(ctx context.Context, search string, role string, offset int, limit int) ([]UserResponseDto, int64, error) {
  var users []UserResponseDto
  var length int64

  query := r.db.WithContext(ctx).Model(&User{}).
  Select("id", "username", "email", "role", "created_at").
  Where("LOWER(username) LIKE ?", "%"+strings.ToLower(search)+"%")

  if role != "" {
    query = query.Where("role = ?", role)
  }  
  query.Count(&length)

  result := query.Scopes(utils.Paginate(offset, limit)).Find(&users).Error
  return users, length, result
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

func (r *repository) FindPasswordById(ctx context.Context, id uuid.UUID) (*string, error) {
  var u User
  result := r.db.WithContext(ctx).Where("id = ?", id).Select("password").First(&u)
  if result.Error != nil {
    return nil, result.Error
  }
  return &u.Password, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
  var u User
  result := r.db.WithContext(ctx).Where("id = ?", id).First(&u)
  if result.Error != nil {
    return nil, result.Error
  }
  return &u, nil
}

func (r *repository) Create(ctx context.Context, u *User) error {
  return r.db.WithContext(ctx).Create(&u).Error
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, u *User) error {
  return r.db.WithContext(ctx).Where("id = ?", id).Updates(u).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
  return r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}