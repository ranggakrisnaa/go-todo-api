package usecase

import (
	"context"
	"fmt"
	"go-todo-api/domain"
	"go-todo-api/domain/converter"
	"go-todo-api/internal/config"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"
	"math"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	FindByID(ctx context.Context, id any) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, todo *entity.Todo) error
	FindAll(ctx context.Context) (*[]entity.Todo, error)
	FindAllWithPagination(ctx context.Context, offset, limit int) (*[]entity.Todo, error)
	Count(ctx context.Context, query string, args ...any) (int64, error)
}

type TodoUsecase struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	JwtService *config.JwtConfig
	TodoRepo   TodoRepository
}

func NewTodoUseCase(t TodoRepository, db *gorm.DB, logger *logrus.Logger, jwtService *config.JwtConfig) *TodoUsecase {
	return &TodoUsecase{
		DB:         db,
		Log:        logger,
		TodoRepo:   t,
		JwtService: jwtService,
	}
}

func (t *TodoUsecase) Create(ctx context.Context, requests []*domain.TodoCreateRequest) ([]*domain.TodoResponse, error) {
	var (
		tx    = t.DB.WithContext(ctx).Begin()
		todos []*domain.TodoResponse
	)

	for _, request := range requests {
		todo := entity.Todo{
			Title:       request.Title,
			UserID:      request.UserID,
			Description: request.Description,
			IsCompleted: request.IsCompleted,
			DueTime:     request.DueTime,
		}

		if err := t.TodoRepo.Create(tx.Statement.Context, &todo); err != nil {
			t.Log.WithError(err).Error("Failed to create user")
			return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
		}

		todos = append(todos, &domain.TodoResponse{
			UUID:        todo.UUID,
			Title:       todo.Title,
			Description: todo.Description,
			IsCompleted: todo.IsCompleted,
			DueTime:     todo.DueTime,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		})
	}

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return todos, nil
}

func (t *TodoUsecase) Update(ctx context.Context, requests []*domain.TodoUpdateRequest) ([]*domain.TodoResponse, error) {
	var (
		tx    = t.DB.WithContext(ctx).Begin()
		todos []*domain.TodoResponse
	)

	for _, request := range requests {
		todo, err := t.TodoRepo.FindByID(tx.Statement.Context, request.ID)
		if err != nil {
			t.Log.WithError(err).Error("Failed to found todo")
			return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
		}

		if request.Title != "" {
			todo.Title = request.Title
		}

		if request.Description != "" {
			todo.Description = request.Description
		}
		todo.IsCompleted = request.IsCompleted
		todo.DueTime = request.DueTime

		if err := t.TodoRepo.Update(tx.Statement.Context, todo); err != nil {
			t.Log.WithError(err).Error("Failed to create user")
			return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
		}

		todos = append(todos, &domain.TodoResponse{
			UUID:        todo.UUID,
			Title:       todo.Title,
			Description: todo.Description,
			IsCompleted: todo.IsCompleted,
			DueTime:     todo.DueTime,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		})
	}

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return todos, nil
}

func (t *TodoUsecase) Delete(ctx context.Context, request *domain.TodoDeleteRequest) ([]*domain.TodoResponse, error) {
	var (
		tx           = t.DB.WithContext(ctx).Begin()
		deletedTodos []*domain.TodoResponse
	)

	todo, err := t.TodoRepo.FindByID(tx.Statement.Context, request.ID)
	if err != nil {
		t.Log.WithError(err).Error("Failed to found todo")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	err = t.TodoRepo.Delete(ctx, todo)
	if err != nil {
		t.Log.WithError(err).Error("Error deleting todo")
		return nil, fmt.Errorf("failed to delete todo: %w", err)
	}

	deletedTodos = append(deletedTodos, converter.TodoUUIDToResponse(todo))

	return deletedTodos, nil
}

func (t *TodoUsecase) FindAllTodo(ctx context.Context, page, size int) ([]*domain.TodoResponse, *domain.PaginationMeta, error) {
	var (
		todos         []entity.Todo
		todoResponses []*domain.TodoResponse
	)

	todosFromRepo, err := t.TodoRepo.FindAllWithPagination(ctx, page, size)
	if err != nil {
		t.Log.WithError(err).Error("Failed to find todos")
		return nil, nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	todos = *todosFromRepo

	for _, todo := range todos {
		todoResponses = append(todoResponses, converter.TodoToResponse(&todo))
	}

	totalCount, err := t.TodoRepo.Count(ctx, "1=1")
	if err != nil {
		t.Log.WithError(err).Error("Failed to count todos")
		return nil, nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(size)))

	meta := &domain.PaginationMeta{
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    size,
		TotalCount:  totalCount,
	}

	return todoResponses, meta, nil
}

func (t *TodoUsecase) FindTodoByID(ctx context.Context, request *domain.TodoGetDataRequest) (*domain.TodoResponse, error) {
	todo, err := t.TodoRepo.FindByID(ctx, request.ID)
	if err != nil {
		t.Log.WithError(err).Error("Failed to find todos")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.TodoToResponse(todo), nil
}

// func (t *TodoUsecase) FindTodoByID(ctx context.Context, requests []*domain.TodoUpdateRequest) ([]*domain.TodoResponse, error) {
// 	todo, err := t.TodoRepo.FindByID(tx.Statement.Context, request.ID)
// 	if err != nil {
// 		t.Log.WithError(err).Error("Failed to found todo")
// 		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
// 	}
// }
