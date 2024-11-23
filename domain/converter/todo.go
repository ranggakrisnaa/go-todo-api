package converter

import (
	"go-todo-api/domain"
	"go-todo-api/internal/entity"
)

func TodoToResponse(todo *entity.Todo) *domain.TodoResponse {
	return &domain.TodoResponse{
		UUID:        todo.UUID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
		DueTime:     todo.DueTime,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

func TodoUUIDToResponse(todo *entity.Todo) *domain.TodoResponse {
	return &domain.TodoResponse{
		UUID: todo.UUID,
	}
}
