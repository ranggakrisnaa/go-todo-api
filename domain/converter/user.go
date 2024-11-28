package converter

import (
	"go-todo-api/domain"
	"go-todo-api/internal/entity"
)

func UserToResponse(user *entity.User) *domain.UserResponse {
	return &domain.UserResponse{
		UUID:      user.UUID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToResponseWithToken(user *entity.User, token string) *domain.UserResponse {
	return &domain.UserResponse{
		UUID:  user.UUID,
		Token: token,
	}
}
