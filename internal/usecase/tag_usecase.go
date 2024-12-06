package usecase

import (
	"context"
	"go-todo-api/domain"
	"go-todo-api/domain/converter"
	"go-todo-api/internal/config"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"
	"math"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	FindByID(ctx context.Context, id any) (*entity.Tag, error)
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, tag *entity.Tag) error
	FindAllTag(ctx context.Context, offset, limit int) (*[]entity.Tag, error)
	Count(ctx context.Context, query string, args ...any) (int64, error)
}

type TagUsecase struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	JwtService *config.JwtConfig
	TagRepo    TagRepository
}

func NewTagUsecase(t TagRepository, db *gorm.DB, logger *logrus.Logger, jwtService *config.JwtConfig) *TagUsecase {
	return &TagUsecase{
		DB:         db,
		Log:        logger,
		TagRepo:    t,
		JwtService: jwtService,
	}
}

func (t *TagUsecase) Create(ctx context.Context, requests []*domain.TagCreateRequest) ([]*domain.TagResponse, error) {
	var (
		tx   = t.DB.WithContext(ctx).Begin()
		tags []*domain.TagResponse
	)

	for _, request := range requests {
		tag := entity.Tag{
			Name: request.Name,
		}

		if err := t.TagRepo.Create(tx.Statement.Context, &tag); err != nil {
			t.Log.WithError(err).Error("Failed to create user")
			return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
		}

		tags = append(tags, converter.TagToResponse(&tag))
	}

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return tags, nil
}

func (t *TagUsecase) Update(ctx context.Context, request *domain.TagUpdateRequest) (*domain.TagResponse, error) {
	tx := t.DB.WithContext(ctx).Begin()

	tag, err := t.TagRepo.FindByID(tx.Statement.Context, request.ID)
	if err != nil {
		t.Log.WithError(err).Error("Failed to found tag")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if request.Name != "" {
		tag.Name = request.Name
	}
	if err := t.TagRepo.Update(tx.Statement.Context, tag); err != nil {
		t.Log.WithError(err).Error("Failed to update todo")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.TagToResponse(tag), nil
}

func (t *TagUsecase) Delete(ctx context.Context, request *domain.TagDeleteRequest) (*domain.TagResponse, error) {
	tx := t.DB.WithContext(ctx).Begin()

	tag, err := t.TagRepo.FindByID(tx.Statement.Context, request.ID)
	if err != nil {
		t.Log.WithError(err).Error("Failed to found tag")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := t.TagRepo.Delete(tx.Statement.Context, tag); err != nil {
		t.Log.WithError(err).Error("Failed to update todo")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}
	return converter.TagToResponse(tag), nil
}

func (t *TagUsecase) FindAllTag(ctx context.Context, page, size int) ([]*domain.TagResponse, *domain.PaginationMeta, error) {
	var (
		tags          []entity.Tag
		tagsResponses []*domain.TagResponse
	)
	tagsFromRepo, err := t.TagRepo.FindAllTag(ctx, page, size)
	if err != nil {
		t.Log.WithError(err).Error("Failed to find todos")
		return nil, nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	tags = *tagsFromRepo
	for _, tag := range tags {
		tagsResponses = append(tagsResponses, converter.TagToResponse(&tag))
	}

	totalCount, err := t.TagRepo.Count(ctx, "1=1")
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

	return tagsResponses, meta, nil
}

func (t *TagUsecase) FindTagById(ctx context.Context, request *domain.TagGetDataRequest) (*domain.TagResponse, error) {
	tag, err := t.TagRepo.FindByID(ctx, request.ID)
	if err != nil {
		t.Log.WithError(err).Error("Failed to found tag")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.TagToResponse(tag), nil
}
