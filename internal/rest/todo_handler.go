package rest

import (
	"context"
	"go-todo-api/domain"
	"go-todo-api/internal/rest/middleware"
	"go-todo-api/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TodoUseCase interface {
	Create(ctx context.Context, requests []*domain.TodoCreateRequest) ([]*domain.TodoResponse, error)
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
	)

	if bulkInsert {
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

	for _, todo := range todos {
		if ok, err := util.IsRequestValid(&todo); !ok {
			t.Log.WithError(err).Error("Validation failed for todo")
			errors = append(errors, err)
			continue
		}

		todo.UserID = auth.ID
		response, err := t.UseCase.Create(c, []*domain.TodoCreateRequest{&todo})
		if err != nil {
			t.Log.WithError(err).Error("Error creating Todo")
			errors = append(errors, err)
			continue
		}

		responses = append(responses, response[0])
	}

	if len(errors) > 0 {
		c.JSON(http.StatusMultiStatus, gin.H{
			"status":  true,
			"message": "Some todos failed to process",
			"data":    responses,
			"errors":  errors,
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

}
