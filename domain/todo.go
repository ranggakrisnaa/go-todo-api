package domain

import (
	"time"

	"github.com/google/uuid"
)

type TodoResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title,omitempty" validate:"required,max=255"`
	Description string    `json:"description,omitempty"  validate:"required"`
	IsCompleted bool      `json:"is_completed,omitempty" validate:"required"`
	DueTime     time.Time `json:"due_time,omitempty" validate:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type TodoCreateRequest struct {
	UUID        uuid.UUID `json:"uuid"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title" validate:"required,max=255"`
	TagID       []int     `json:"tag_id"`
	Description string    `json:"description"  validate:"required"`
	IsCompleted bool      `json:"is_completed"`
	DueTime     time.Time `json:"due_time" validate:"required"`
}

type TodoUpdateRequest struct {
	ID          uint      `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title,omitempty" validate:"max=255"`
	Description string    `json:"description,omitempty"`
	IsCompleted bool      `json:"is_completed,omitempty"`
	DueTime     time.Time `json:"due_time,omitempty"`
}

type TodoGetDataRequest struct {
	ID uint `json:"id"`
}

type TodoDeleteRequest struct {
	ID uint `json:"id"`
}
