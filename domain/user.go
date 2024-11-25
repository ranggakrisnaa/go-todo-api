package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	UUID      uuid.UUID `json:"uuid,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required"`
}

type LoginUserRequest struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email" validate:"required,max=255"`
	Password string    `json:"password" validate:"required,max=100"`
}

type GetUserId struct {
	ID uint `json:"id"`
}

type CurrentUserRequest struct {
	GetUserId
}

type LogoutUserRequest struct {
	GetUserId
}
