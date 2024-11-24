package converter

import (
	"go-todo-api/domain"
	"go-todo-api/internal/entity"
)

func TodoToResponse(todo *entity.Todo) *domain.TodoResponse {
	var tagResponses []domain.TagResponse
	if len(todo.Tag) > 0 {
		for _, tag := range todo.Tag {
			tagResponses = append(tagResponses, domain.TagResponse{
				UUID:      tag.UUID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			})
		}
	}

	return &domain.TodoResponse{
		UUID:        todo.UUID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
		DueTime:     todo.DueTime,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
		Tags:        tagResponses,
	}
}

func TodoUUIDToResponse(todo *entity.Todo) *domain.TodoResponse {
	return &domain.TodoResponse{
		UUID: todo.UUID,
	}
}
