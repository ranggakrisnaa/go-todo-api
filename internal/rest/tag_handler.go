package rest

import (
	"context"
	"fmt"
	"go-todo-api/domain"
	"go-todo-api/internal/util"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TagUsecase interface {
	Create(ctx context.Context, requests []*domain.TagCreateRequest) ([]*domain.TagResponse, error)
}
type TagHandler struct {
	Log     *logrus.Logger
	UseCase TagUsecase
}

func NewTagHandler(r *gin.Engine, t TagUsecase, log *logrus.Logger) {
	handler := &TagHandler{
		UseCase: t,
		Log:     log,
	}

	r.POST("v1/tags", handler.Create)
}

func (t *TagHandler) Create(c *gin.Context) {
	var (
		singleTag  domain.TagCreateRequest
		tags       []domain.TagCreateRequest
		responses  []*domain.TagResponse
		errors     []error
		bulkInsert = c.Query("bulk") != ""
		wg         sync.WaitGroup
		mu         sync.Mutex
	)

	if bulkInsert {
		if err := c.ShouldBindJSON(&tags); err != nil {
			t.Log.WithError(err).Error("Error parsing request body (bulk mode)")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&singleTag); err != nil {
			t.Log.WithError(err).Error("Error parsing request body (single mode)")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
		tags = append(tags, singleTag)
	}

	resultChan := make(chan *domain.TagResponse, len(tags))
	errChan := make(chan error, len(tags))

	for _, tag := range tags {
		wg.Add(1)
		go func(tag domain.TagCreateRequest) {
			defer wg.Done()

			if ok, err := util.IsRequestValid(&tag); !ok {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			response, err := t.UseCase.Create(c, []*domain.TagCreateRequest{&tag})
			if err != nil {
				errChan <- err
				return
			}

			if len(response) == 0 {
				errChan <- fmt.Errorf("no response returned for tag: %v", tag)
				return
			}

			resultChan <- response[0]
		}(tag)
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

	c.JSON(http.StatusOK, domain.Response[[]*domain.TagResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Todos created successfully",
		Data:       responses,
	})
}
