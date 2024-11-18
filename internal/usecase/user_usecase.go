package usecase

import (
	"context"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	CountById(ctx context.Context, user *entity.User) (int64, error)
	FindByEmailOrName(ctx context.Context, user *entity.User) (*entity.User, error)
}

type UserUseCase struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	userRepo UserRepository
}

func NewUserUseCase(u UserRepository, db *gorm.DB, logger *logrus.Logger) *UserUseCase {
	return &UserUseCase{
		userRepo: u,
		Log:      logger,
		DB:       db,
	}
}

func (u *UserUseCase) Create(ctx context.Context, user *entity.User) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if existingUser, _ := u.userRepo.FindByEmailOrName(tx.Statement.Context, user); existingUser != nil {
		u.Log.Warnf("Error checking for existing user")
		return util.NewCustomError(int(util.ErrInternalServerErrorCode), "User with email or name already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.WithError(err).Error("Failed to generate bcrype hash")
		return util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	user.Password = string(hashedPassword)

	if err = u.userRepo.Create(tx.Statement.Context, user); err != nil {
		u.Log.WithError(err).Error("Failed to create user")
		return util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	u.Log.Infof("User successfully created: %+v", user)
	return err
}
