package postgresql

import (
	"go-todo-api/internal/entity"

	"gorm.io/gorm"
)

type TodoRepository struct {
	*BaseRepository[entity.Todo]
	DB *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{
		BaseRepository: NewBaseRepository[entity.Todo](db),
		DB:             db,
	}
}
