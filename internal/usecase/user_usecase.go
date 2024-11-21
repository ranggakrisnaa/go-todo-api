package usecase

import (
	"context"
	"go-todo-api/domain"
	"go-todo-api/domain/converter"
	"go-todo-api/internal/config"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	CountByEmailOrName(ctx context.Context, user *entity.User) (int64, error)
	FindByEmailOrName(ctx context.Context, email string, name string) (*entity.User, error)
	FindByUUID(ctx context.Context, uuid string) (*entity.User, error)
	FindByID(ctx context.Context, id any) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, user *entity.User) error
}

type UserUseCase struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	userRepo   UserRepository
	jwtService *config.JwtConfig
}

func NewUserUseCase(u UserRepository, db *gorm.DB, logger *logrus.Logger, jwtService *config.JwtConfig) *UserUseCase {
	return &UserUseCase{
		userRepo:   u,
		Log:        logger,
		DB:         db,
		jwtService: jwtService,
	}
}

func (u *UserUseCase) Create(ctx context.Context, request *domain.RegisterUserRequest) (*domain.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	existingUser, _ := u.userRepo.FindByEmailOrName(tx.Statement.Context, request.Email, request.Name)
	if existingUser != nil {
		u.Log.Warnf("Error checking for existing user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), "User with email or name already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.WithError(err).Error("Failed to generate bcrype hash")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	userPayload := &entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	if err := u.userRepo.Create(tx.Statement.Context, userPayload); err != nil {
		u.Log.WithError(err).Error("Failed to create user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	u.Log.Infof("User successfully created: %+v", userPayload)
	return converter.UserToResponse(userPayload), nil
}

func (u *UserUseCase) Login(ctx context.Context, request *domain.LoginUserRequest) (*domain.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user, err := u.userRepo.FindByEmailOrName(tx.Statement.Context, request.Email, "")
	if err != nil {
		u.Log.WithError(err).Error("Failed to found user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.Log.WithError(err).Error("Failed to compare user password with bcrype hash")
		return nil, util.NewCustomError(int(util.ErrUnauthorizedCode), err.Error())
	}

	token, err := u.jwtService.CreateToken(user)
	if err != nil {
		u.Log.WithError(err).Error("Failed to create jwt token")
		return nil, util.NewCustomError(int(util.ErrUnauthorizedCode), err.Error())
	}

	user.Token = token
	if err := u.userRepo.Update(tx.Statement.Context, user); err != nil {
		u.Log.WithError(err).Error("Failed to update user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.UserToResponseWithToken(user, token), nil
}

func (u *UserUseCase) GetUserID(ctx context.Context, request *domain.GetUserId) (*entity.User, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user, err := u.userRepo.FindByID(tx.Statement.Context, request.ID)
	if err != nil {
		u.Log.WithError(err).Error("Failed to find user by UUID")
		return nil, util.NewCustomError(int(util.ErrNotFoundCode), "User not found")
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	u.Log.Infof("User ID found: %d", user.ID)
	return user, nil
}

func (u *UserUseCase) Logout(ctx context.Context, request *domain.LogoutUserRequest) (bool, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user, err := u.userRepo.FindByID(tx.Statement.Context, request.ID)
	if err != nil {
		u.Log.WithError(err).Error("Failed to found user")
		return false, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	user.Token = ""
	if err := u.userRepo.Update(tx.Statement.Context, user); err != nil {
		u.Log.WithError(err).Error("Failed to update token user")
		return false, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return false, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return true, nil
}

func (u *UserUseCase) Current(ctx context.Context, request *domain.CurrentUserRequest) (*domain.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user, err := u.userRepo.FindByID(tx.Statement.Context, request.ID)
	if err != nil {
		u.Log.WithError(err).Error("Failed to found user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.UserToResponse(user), nil
}
