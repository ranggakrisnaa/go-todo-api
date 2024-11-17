package usecase

import (
	"context"
	"go-todo-api/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
}

type UserUseCase struct {
	userRepo UserRepository
}

func NewUserUseCase(u UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: u,
	}
}

func (u *UserUseCase) Create(ctx context.Context, user *entity.User) error {
	return nil
}
