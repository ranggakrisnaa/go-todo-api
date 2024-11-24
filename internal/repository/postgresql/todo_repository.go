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

func (r *TodoRepository) FindAll(ctx context.Context, offset, limit int) (*[]entity.Todo, error) {
	var todos []entity.Todo
	err := r.DB.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Preload("Tag").
		Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return &todos, nil
}

func (r *TodoRepository) CreateTodoTag(ctx context.Context, todoTag *entity.TodoTag) error {
	return r.DB.WithContext(ctx).Create(todoTag).Error
}

func (r *TodoRepository) FindTodoTagByTodoID(ctx context.Context, todoID uint) ([]entity.TodoTag, error) {
	var todoTags []entity.TodoTag
	if err := r.DB.WithContext(ctx).Where("todo_id = ?", todoID).Find(&todoTags).Error; err != nil {
		return nil, err
	}
	return todoTags, nil
}

func (r *TodoRepository) FindTodoTagByTagID(ctx context.Context, tagID uint) ([]entity.TodoTag, error) {
	var todoTags []entity.TodoTag
	if err := r.DB.WithContext(ctx).Where("tag_id = ?", tagID).Find(&todoTags).Error; err != nil {
		return nil, err
	}
	return todoTags, nil
}

func (r *TodoRepository) DeleteTodoTag(ctx context.Context, todoTags []entity.TodoTag) error {
	for _, todoTag := range todoTags {
		if err := r.DB.WithContext(ctx).Unscoped().Delete(&todoTag).Error; err != nil {
			return err
		}
	}
	return nil
}
