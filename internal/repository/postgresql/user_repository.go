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

func (r *UserRepository) Create(ctx context.Context, user *entity.User) (err error) {
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	return nil
}
