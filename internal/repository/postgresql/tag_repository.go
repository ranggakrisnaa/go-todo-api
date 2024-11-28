package postgresql

import (
	"context"
	"go-todo-api/internal/entity"

	"gorm.io/gorm"
)

type TagRepository struct {
	*BaseRepository[entity.Tag]
	DB *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		BaseRepository: NewBaseRepository[entity.Tag](db),
		DB:             db,
	}
}

func (r *TagRepository) FindAllTag(ctx context.Context, offset, limit int) (*[]entity.Tag, error) {
	var tags []entity.Tag
	err := r.DB.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return &tags, nil
}
