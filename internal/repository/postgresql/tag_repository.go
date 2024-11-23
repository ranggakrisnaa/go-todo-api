package postgresql

import (
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
