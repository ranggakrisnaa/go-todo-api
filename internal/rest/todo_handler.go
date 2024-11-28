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

type TodoUsecase interface {
	Create(ctx context.Context, requests []*domain.TodoCreateRequest) ([]*domain.TodoResponse, error)
	Update(ctx context.Context, requests []*domain.TodoUpdateRequest) ([]*domain.TodoResponse, error)
	Delete(ctx context.Context, request *domain.TodoDeleteRequest) ([]*domain.TodoResponse, error)
	FindAllTodo(ctx context.Context, page, size int) ([]*domain.TodoResponse, *domain.PaginationMeta, error)
	FindTodoByID(ctx context.Context, request *domain.TodoGetDataRequest) (*domain.TodoResponse, error)
}

type TodoHandler struct {
	Log     *logrus.Logger
	UseCase TodoUsecase
}

func NewTodoHandler(r *gin.Engine, t TodoUsecase, log *logrus.Logger) {
	handler := &TodoHandler{
		UseCase: t,
		Log:     log,
	}

	requiredRole := middleware.NewRequiredRole()
	r.POST("v1/todos", requiredRole.RoleCheck(), handler.Create)
	r.GET("v1/todos", handler.FindAllTodo)
	r.GET("v1/todos/:id", handler.FindTodoById)
	r.PUT("v1/todos/:id", requiredRole.RoleCheck(), handler.Update)
	r.DELETE("v1/todos/:id", requiredRole.RoleCheck(), handler.Delete)
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
		todoIdParam = c.Param("id")
		todoIds     = c.QueryArray("ids")
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

	if len(todoIds) > 0 && len(todoIds) != len(todos) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "Mismatched number of IDs and Todo items"})
		return
	}

	if len(todoIds) == 0 && todoIdParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "No IDs provided for update"})
		return
	}

	if len(todoIds) == 0 && todoIdParam != "" {
		todoIds = make([]string, len(todos))
		for i := range todos {
			todoIds[i] = todoIdParam
		}
	}

	resultChan := make(chan *domain.TodoResponse, len(todos)*len(todoIds))
	errChan := make(chan error, len(todos)*len(todoIds))

	for i := range todos {
		wg.Add(1)
		go func(todo domain.TodoUpdateRequest, todoIdStr string) {
			defer wg.Done()

			if ok, err := util.IsRequestValid(&todo); !ok {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			todoId, err := strconv.ParseUint(todoIdStr, 10, 64)
			if err != nil {
				t.Log.WithError(err).Errorf("Invalid todo ID: %s", todoIdStr)
				errChan <- fmt.Errorf("invalid todo ID: %s", todoIdStr)
				return
			}

			todo.ID = uint(todoId)
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
		}(todos[i], todoIds[i])

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

func (t *TodoHandler) Delete(c *gin.Context) {
	var (
		responses   []*domain.TodoResponse
		errors      []error
		todoIdParam = c.Param("id")
		todoIds     = c.QueryArray("ids")
		wg          sync.WaitGroup
	)

	if len(todoIds) == 0 && todoIdParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "No IDs provided for delete"})
		return
	}

	if len(todoIds) == 0 && todoIdParam != "" {
		todoIds = append(todoIds, todoIdParam)
	}

	resultChan := make(chan *domain.TodoResponse, len(todoIds))
	errChan := make(chan error, len(todoIds))

	for _, todoIdStr := range todoIds {
		wg.Add(1)
		go func(todoIdStr string) {
			defer wg.Done()

			todoId, err := strconv.ParseUint(todoIdStr, 10, 64)
			if err != nil {
				t.Log.WithError(err).Errorf("Invalid todo ID: %s", todoIdStr)
				errChan <- fmt.Errorf("invalid todo ID: %s", todoIdStr)
				return
			}

			todo := &domain.TodoDeleteRequest{ID: uint(todoId)}
			response, err := t.UseCase.Delete(c, todo)
			if err != nil {
				errChan <- err
				return
			}

			if len(response) == 0 {
				errChan <- fmt.Errorf("no response returned for todo: %v", todoId)
				return
			}

			resultChan <- response[0]
		}(todoIdStr)
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
		Message:    "Todos deleted successfully",
		Data:       responses,
	})
}

func (t *TodoHandler) FindAllTodo(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		t.Log.WithError(err).Warn("Invalid parsing data")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		t.Log.WithError(err).Warn("Invalid parsing data")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	responses, meta, err := t.UseCase.FindAllTodo(c, page, size)
	if err != nil {
		t.Log.WithError(err).Error("Error find todo")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[[]*domain.TodoResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Todos data retrieved successfully",
		Data:       responses,
		Meta:       meta,
	})
}

func (t *TodoHandler) FindTodoById(c *gin.Context) {
	todoId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		t.Log.WithError(err).Warn("Invalid parsing data")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	todo := &domain.TodoGetDataRequest{ID: uint(todoId)}
	response, err := t.UseCase.FindTodoByID(c, todo)
	if err != nil {
		t.Log.WithError(err).Error("Error finding todo")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[*domain.TodoResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Todo data retrieved successfully",
		Data:       response,
	})
}
