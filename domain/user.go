package domain

import (
	"go-todo-api/internal/entity"

	"github.com/google/uuid"
)

type UserResponse struct {
	UUID  uuid.UUID `json:"uuid,omitempty"`
	Name  string    `json:"name,omitempty"`
	Email string    `json:"email,omitempty"`
	Token string    `json:"token,omitempty"`
}

type RegisterUserRequest struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name" validate:"required,max=100"`
	Email    string    `json:"email" validate:"required,email,max=255"`
	Password string    `json:"password" validate:"required"`
}

type LoginUserRequest struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email" validate:"required,max=255"`
	Password string    `json:"password" validate:"required,max=100"`
}

func UserToResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
	}
}
