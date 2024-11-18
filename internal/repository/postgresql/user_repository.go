package postgresql

import (
	"context"
	"go-todo-api/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) CountById(ctx context.Context, user *entity.User) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Count(&count).Error
	return count, err
}

func (r *UserRepository) FindByEmailOrName(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := r.DB.WithContext(ctx).Where("email = ? OR name = ? ", user.Email, user.Name).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
