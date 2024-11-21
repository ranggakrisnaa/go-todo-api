package postgresql

import (
	"context"
	"go-todo-api/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository[entity.User]
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository[entity.User](db),
		DB:             db,
	}
}

func (r *UserRepository) CountByEmailOrName(ctx context.Context, user *entity.User) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&entity.User{}).Where("email = ? OR name = ? ", user.Email, user.Name).Count(&count).Error
	return count, err
}

func (r *UserRepository) FindByEmailOrName(ctx context.Context, email string, name string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.WithContext(ctx).Where("email = ? OR name = ?", email, name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*entity.User, error) {
	var user entity.User
	err := r.DB.WithContext(ctx).Where("uuid = ? ", uuid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
