package postgresql

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Unscoped().Delete(entity).Error
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).
		Where("id = ?", id).
		Take(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Count(ctx context.Context, query string, args ...any) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(new(T)).
		Where(query, args...).
		Count(&count).Error
	return count, err
}
