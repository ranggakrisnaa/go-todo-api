package converter

import (
	"go-todo-api/domain"
	"go-todo-api/internal/entity"
)

func TagToResponse(tag *entity.Tag) *domain.TagResponse {
	return &domain.TagResponse{
		UUID:      tag.UUID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}
