package postgresql

import (
	"context"
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

func (r *TodoRepository) FindAll(ctx context.Context) (*[]entity.Todo, error) {
	var todos []entity.Todo

	err := r.DB.WithContext(ctx).Find(&todos).Error
	if err != nil {
		return nil, err
	}

	return &todos, nil
}
