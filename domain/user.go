package domain

import (
	"go-todo-api/internal/entity"

	"github.com/google/uuid"
)

type UserResponse struct {
	UUID  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func UserToResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
	}
}
