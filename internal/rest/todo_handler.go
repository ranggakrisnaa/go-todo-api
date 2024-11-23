package rest

import (
	"context"
	"fmt"
	"go-todo-api/domain"
	"go-todo-api/internal/rest/middleware"
	"go-todo-api/internal/util"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TodoError struct {
	Message string `json:"message"`
}

func (e *TodoError) Error() string {
	return e.Message
}

type TodoUseCase interface {
	Create(ctx context.Context, requests []*domain.TodoCreateRequest) ([]*domain.TodoResponse, error)
	Update(ctx context.Context, requests []*domain.TodoUpdateRequest) ([]*domain.TodoResponse, error)
}

type TodoHandler struct {
	Log     *logrus.Logger
	UseCase TodoUseCase
}

func NewTodoHandler(r *gin.Engine, t TodoUseCase, log *logrus.Logger) {
	handler := &TodoHandler{
		UseCase: t,
		Log:     log,
	}

	requiredRole := middleware.NewRequiredRole()
	r.POST("v1/todos", requiredRole.RoleCheck(), handler.Create)
	r.PUT("v1/todos/:id", requiredRole.RoleCheck(), handler.Update)
}

func (t *TodoHandler) Create(c *gin.Context) {
	var (
		singleTodo domain.TodoCreateRequest
		todos      []domain.TodoCreateRequest
		auth       = middleware.GetUser(c)
		responses  []*domain.TodoResponse
		errors     []error
		bulkInsert = c.Query("bulk") != ""
		wg         sync.WaitGroup
		mu         sync.Mutex
	)

	if bulkInsert {
		if err := c.ShouldBindJSON(&todos); err != nil {
			t.Log.WithError(err).Error("Error parsing request body (bulk mode)")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&singleTodo); err != nil {
			t.Log.WithError(err).Error("Error parsing request body (single mode)")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
		todos = append(todos, singleTodo)
	}

	resultChan := make(chan *domain.TodoResponse, len(todos))
	errChan := make(chan error, len(todos))

	for _, todo := range todos {
		wg.Add(1)
		go func(todo domain.TodoCreateRequest) {
			defer wg.Done()

			if ok, err := util.IsRequestValid(&todo); !ok {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			todo.UserID = auth.ID
			response, err := t.UseCase.Create(c, []*domain.TodoCreateRequest{&todo})
			if err != nil {
				errChan <- err
				return
			}

			if len(response) == 0 {
				errChan <- fmt.Errorf("no response returned for todo: %v", todo)
				return
			}

			resultChan <- response[0]
		}(todo)

	}

	wg.Wait()
	close(resultChan)
	close(errChan)

	for response := range resultChan {
		responses = append(responses, response)
	}
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		var errorResponses []map[string]string
		for _, err := range errors {
			if todoErr, ok := err.(*TodoError); ok {
				errorResponses = append(errorResponses, map[string]string{"message": todoErr.Message})
			} else {
				errorResponses = append(errorResponses, map[string]string{"message": err.Error()})
			}
		}

		c.JSON(http.StatusMultiStatus, gin.H{
			"status":  true,
			"message": "Some todos failed to process",
			"data":    responses,
			"errors":  errorResponses,
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response[[]*domain.TodoResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Todos created successfully",
		Data:       responses,
	})
}

func (t *TodoHandler) Update(c *gin.Context) {
	var (
		singleTodo  domain.TodoUpdateRequest
		todos       []domain.TodoUpdateRequest
		auth        = middleware.GetUser(c)
		responses   []*domain.TodoResponse
		errors      []error
		bulkUpdate  = c.Query("bulk") != ""
		bookIdParam = c.Param("id")
		bookIds     = c.QueryArray("ids")
		wg          sync.WaitGroup
		mu          sync.Mutex
	)

	if bulkUpdate {
		if err := c.ShouldBindJSON(&todos); err != nil {
			t.Log.WithError(err).Error("Error parsing request body (bulk mode)")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "Invalid request body for bulk"})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&singleTodo); err != nil {
			t.Log.WithError(err).Error("Error parsing request body (single mode)")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "Invalid request body"})
			return
		}
		todos = append(todos, singleTodo)
	}

	if len(bookIds) > 0 && len(bookIds) != len(todos) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "Mismatched number of IDs and Todo items"})
		return
	}

	if len(bookIds) == 0 && bookIdParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "No IDs provided for update"})
		return
	}

	if len(bookIds) == 0 && bookIdParam != "" {
		bookIds = make([]string, len(todos))
		for i := range todos {
			bookIds[i] = bookIdParam
		}
	}

	resultChan := make(chan *domain.TodoResponse, len(todos)*len(bookIds))
	errChan := make(chan error, len(todos)*len(bookIds))

	for i := range todos {
		wg.Add(1)
		go func(todo domain.TodoUpdateRequest, bookIdStr string) {
			defer wg.Done()

			if ok, err := util.IsRequestValid(&todo); !ok {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			bookId, err := strconv.ParseUint(bookIdStr, 10, 64)
			if err != nil {
				t.Log.WithError(err).Errorf("Invalid book ID: %s", bookIdStr)
				errChan <- fmt.Errorf("invalid book ID: %s", bookIdStr)
				return
			}

			todo.ID = uint(bookId)
			todo.UserID = auth.ID

			response, err := t.UseCase.Update(c, []*domain.TodoUpdateRequest{&todo})
			if err != nil {
				errChan <- err
				return
			}

			if len(response) == 0 {
				errChan <- fmt.Errorf("no response returned for todo: %v", todo)
				return
			}

			resultChan <- response[0]
		}(todos[i], bookIds[i])

	}

	wg.Wait()
	close(resultChan)
	close(errChan)

	for response := range resultChan {
		responses = append(responses, response)
	}
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		var errorResponses []map[string]string
		for _, err := range errors {
			if todoErr, ok := err.(*TodoError); ok {
				errorResponses = append(errorResponses, map[string]string{"message": todoErr.Message})
			} else {
				errorResponses = append(errorResponses, map[string]string{"message": err.Error()})
			}
		}

		c.JSON(http.StatusMultiStatus, gin.H{
			"status":  true,
			"message": "Some todos failed to process",
			"data":    responses,
			"errors":  errorResponses,
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response[[]*domain.TodoResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Todos updated successfully",
		Data:       responses,
	})
}
