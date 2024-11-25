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

	"github.com/gocraft/work"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	FindByID(ctx context.Context, id any) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, todo *entity.Todo) error
	FindAll(ctx context.Context, offset, limit int) (*[]entity.Todo, error)
	FindAllWithPagination(ctx context.Context, offset, limit int) (*[]entity.Todo, error)
	Count(ctx context.Context, query string, args ...any) (int64, error)
	CreateTodoTag(ctx context.Context, todoTag *entity.TodoTag) error
	FindTodoTagByTodoID(ctx context.Context, todoID uint) ([]entity.TodoTag, error)
	DeleteTodoTag(ctx context.Context, todoTags []entity.TodoTag) error
	FindTodoTagByTagID(ctx context.Context, tagID uint) ([]entity.TodoTag, error)
	FindUserById(ctx context.Context, id any) (*entity.User, error)
}

type TodoUsecase struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	JwtService *config.JwtConfig
	TodoRepo   TodoRepository
	Enqueuer   *work.Enqueuer
}

func NewTodoUseCase(t TodoRepository, db *gorm.DB, logger *logrus.Logger, jwtService *config.JwtConfig, enqueuer *work.Enqueuer) *TodoUsecase {
	return &TodoUsecase{
		DB:         db,
		Log:        logger,
		TodoRepo:   t,
		JwtService: jwtService,
		Enqueuer:   enqueuer,
	}
}

func (t *TodoUsecase) enqueueEmail(to string, todo *entity.Todo, status string) error {
	subject := fmt.Sprintf("Your new todo \"%s\" has been %s successfully.", todo.Title, status)
	body := fmt.Sprintf(`
    <html>
        <body>
            <h2>Your new todo "<strong>%s</strong>" has been  %s successfully.</h2>
            <p><strong>Status Completed:</strong> %v</p>
            <p><strong>Description:</strong> %s</p>
            <p><strong>Due Time:</strong> %s</p>
        </body>
    </html>
    `, todo.Title, status, todo.IsCompleted, todo.Description, todo.DueTime)

	_, err := t.Enqueuer.Enqueue("send_email", work.Q{
		"to":      to,
		"subject": subject,
		"body":    body,
	})
	if err != nil {
		t.Log.WithError(err).Error("Failed to enqueue email task")
		return err
	}

	t.Log.Info("Email task enqueued successfully")
	return nil
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

		for _, tagID := range request.TagID {
			todoTag, err := t.TodoRepo.FindTodoTagByTagID(tx.Statement.Context, uint(tagID))
			if err != nil {
				t.Log.WithError(err).Error("Failed to check existing todo_tag")
				return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
			}

			if todoTag != nil {
				continue
			}

			newTodoTag := &entity.TodoTag{
				TodoID: todo.ID,
				TagID:  uint(tagID),
			}

			if err := t.TodoRepo.CreateTodoTag(tx.Statement.Context, newTodoTag); err != nil {
				t.Log.WithError(err).Error("Failed to create todo_tag")
				tx.Rollback()
				return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
			}
		}

		user, _ := t.TodoRepo.FindUserById(ctx, todo.UserID)
		if err := t.enqueueEmail(user.Email, &todo, "created"); err != nil {
			t.Log.WithError(err).Error("Failed to enqueue email after creating todo")
		}

		todos = append(todos, converter.TodoToResponse(&todo))
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
			t.Log.WithError(err).Error("Failed to update todo")
			return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
		}

		if len(request.TagID) > 0 {
			existingTags, err := t.TodoRepo.FindTodoTagByTodoID(tx.Statement.Context, todo.ID)
			if err != nil {
				t.Log.WithError(err).Error("Failed to find todo_tags")
				tx.Rollback()
				return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
			}

			if len(existingTags) > 0 {
				if err := t.TodoRepo.DeleteTodoTag(tx.Statement.Context, existingTags); err != nil {
					t.Log.WithError(err).Error("Failed to delete todo_tags")
					tx.Rollback()
					return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
				}
			}

			for _, tagID := range request.TagID {
				todoTag, err := t.TodoRepo.FindTodoTagByTagID(tx.Statement.Context, uint(tagID))
				if err != nil {
					t.Log.WithError(err).Error("Failed to check existing todo_tag")
					return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
				}

				if todoTag != nil {
					continue
				}

				newTodoTag := &entity.TodoTag{
					TodoID: todo.ID,
					TagID:  uint(tagID),
				}

				if err := t.TodoRepo.CreateTodoTag(tx.Statement.Context, newTodoTag); err != nil {
					t.Log.WithError(err).Error("Failed to create todo_tag")
					tx.Rollback()
					return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
				}
			}
		}

		user, _ := t.TodoRepo.FindUserById(ctx, todo.UserID)
		if err := t.enqueueEmail(user.Email, todo, "updated"); err != nil {
			t.Log.WithError(err).Error("Failed to enqueue email after updated todo")
		}

		todos = append(todos, converter.TodoToResponse(todo))
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

	user, _ := t.TodoRepo.FindUserById(ctx, todo.UserID)
	if err := t.enqueueEmail(user.Email, todo, "deleted"); err != nil {
		t.Log.WithError(err).Error("Failed to enqueue email after deleted todo")
	}
	deletedTodos = append(deletedTodos, converter.TodoUUIDToResponse(todo))

	return deletedTodos, nil
}

func (t *TodoUsecase) FindAllTodo(ctx context.Context, page, size int) ([]*domain.TodoResponse, *domain.PaginationMeta, error) {
	var (
		todos         []entity.Todo
		todoResponses []*domain.TodoResponse
	)

	todosFromRepo, err := t.TodoRepo.FindAll(ctx, page, size)
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
