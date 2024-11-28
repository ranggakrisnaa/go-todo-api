package rest

import (
	"context"
	"fmt"
	"go-todo-api/domain"
	"go-todo-api/internal/util"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TagUsecase interface {
	Create(ctx context.Context, requests []*domain.TagCreateRequest) ([]*domain.TagResponse, error)
	Update(ctx context.Context, request *domain.TagUpdateRequest) (*domain.TagResponse, error)
	Delete(ctx context.Context, request *domain.TagDeleteRequest) (*domain.TagResponse, error)
	FindAllTag(ctx context.Context, offset, limit int) ([]*domain.TagResponse, *domain.PaginationMeta, error)
	FindTagById(ctx context.Context, request *domain.TagGetDataRequest) (*domain.TagResponse, error)
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
	r.PUT("v1/tags/:id", handler.Update)
	r.GET("v1/tags", handler.FindAllTag)
	r.GET("v1/tags/:id", handler.FindTagById)
	r.DELETE("v1/tags/:id", handler.Delete)
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
		Message:    "Tags created successfully",
		Data:       responses,
	})
}

func (t *TagHandler) Update(c *gin.Context) {
	var (
		tag           domain.TagUpdateRequest
		errValidation error
		ok            bool
	)

	tagId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		t.Log.WithError(err).Warn("Invalid parsing data")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&tag); err != nil {
		t.Log.WithError(err).Error("Error parsing request body (single mode)")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if ok, errValidation = util.IsRequestValid(&tag); !ok {
		t.Log.WithError(errValidation).Error("Error request body validation")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errValidation.Error()})
		return
	}

	tag.ID = uint(tagId)
	response, err := t.UseCase.Update(c, &tag)
	if err != nil {
		t.Log.WithError(err).Error("Error update tag")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[*domain.TagResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Tag updated successfully",
		Data:       response,
	})
}

func (t *TagHandler) Delete(c *gin.Context) {
	tagId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		t.Log.WithError(err).Warn("Invalid parsing data")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	response, err := t.UseCase.Delete(c, &domain.TagDeleteRequest{ID: uint(tagId)})
	if err != nil {
		t.Log.WithError(err).Error("Error update tag")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[*domain.TagResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Tag deleted successfully",
		Data:       response,
	})
}

func (t *TagHandler) FindAllTag(c *gin.Context) {
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

	responses, meta, err := t.UseCase.FindAllTag(c, page, size)
	if err != nil {
		t.Log.WithError(err).Error("Error find todo")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[[]*domain.TagResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Tags data retrieved successfully",
		Data:       responses,
		Meta:       meta,
	})
}

func (t *TagHandler) FindTagById(c *gin.Context) {
	tagId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		t.Log.WithError(err).Warn("Invalid parsing data")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	tag := &domain.TagGetDataRequest{ID: uint(tagId)}
	response, err := t.UseCase.FindTagById(c, tag)
	if err != nil {
		t.Log.WithError(err).Error("Error finding tag")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[*domain.TagResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "Tag data retrieved successfully",
		Data:       response,
	})
}
