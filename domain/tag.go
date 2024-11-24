package domain

import (
	"time"

	"github.com/google/uuid"
)

type TagResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name,omitempty" validate:"required,max=255"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type TagCreateRequest struct {
	Name string `json:"name,omitempty" validate:"required,max=255"`
}

type TagUpdateRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name,omitempty" validate:"max=255"`
}

type TagDeleteRequest struct {
	ID uint `json:"id"`
}

type TagGetDataRequest struct {
	ID uint `json:"id"`
}
