package usecase

import (
	"context"
	"go-todo-api/domain"
	"go-todo-api/internal/config"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	FindByID(ctx context.Context, id any) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
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
	tx := t.DB.WithContext(ctx).Begin()

	var todos []*domain.TodoResponse

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
	tx := t.DB.WithContext(ctx).Begin()

	var todos []*domain.TodoResponse

	for _, request := range requests {
		todo, err := t.TodoRepo.FindByID(tx.Statement.Context, request.ID)
		if err != nil {

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
