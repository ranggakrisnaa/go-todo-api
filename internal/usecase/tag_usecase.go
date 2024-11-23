package usecase

import (
	"context"
	"go-todo-api/domain"
	"go-todo-api/domain/converter"
	"go-todo-api/internal/config"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TagRepository interface {
	Create(ctx context.Context, todo *entity.Tag) error
}

type TagUseCase struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	JwtService *config.JwtConfig
	TagRepo    TagRepository
}

func NewTagUsecase(t TagRepository, db *gorm.DB, logger *logrus.Logger, jwtService *config.JwtConfig) *TagUseCase {
	return &TagUseCase{
		DB:         db,
		Log:        logger,
		TagRepo:    t,
		JwtService: jwtService,
	}
}

func (t *TagUseCase) Create(ctx context.Context, requests []*domain.TagCreateRequest) ([]*domain.TagResponse, error) {
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
