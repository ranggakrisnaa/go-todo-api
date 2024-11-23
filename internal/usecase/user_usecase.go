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

type UserUsecase struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	UserRepo   UserRepository
	JwtService *config.JwtConfig
}

func NewUserUsecase(u UserRepository, db *gorm.DB, logger *logrus.Logger, jwtService *config.JwtConfig) *UserUsecase {
	return &UserUsecase{
		UserRepo:   u,
		Log:        logger,
		DB:         db,
		JwtService: jwtService,
	}
}

func (u *UserUsecase) Create(ctx context.Context, request *domain.RegisterUserRequest) (*domain.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()

	existingUser, _ := u.UserRepo.FindByEmailOrName(tx.Statement.Context, request.Email, request.Name)
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

	if err := u.UserRepo.Create(tx.Statement.Context, userPayload); err != nil {
		u.Log.WithError(err).Error("Failed to create user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.UserToResponse(userPayload), nil
}

func (u *UserUsecase) Login(ctx context.Context, request *domain.LoginUserRequest) (*domain.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()

	user, err := u.UserRepo.FindByEmailOrName(tx.Statement.Context, request.Email, "")
	if err != nil {
		u.Log.WithError(err).Error("Failed to found user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.Log.WithError(err).Error("Failed to compare user password with bcrype hash")
		return nil, util.NewCustomError(int(util.ErrUnauthorizedCode), err.Error())
	}

	token, err := u.JwtService.CreateToken(user)
	if err != nil {
		u.Log.WithError(err).Error("Failed to create jwt token")
		return nil, util.NewCustomError(int(util.ErrUnauthorizedCode), err.Error())
	}

	user.Token = token
	if err := u.UserRepo.Update(tx.Statement.Context, user); err != nil {
		u.Log.WithError(err).Error("Failed to update user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.UserToResponseWithToken(user, token), nil
}

func (u *UserUsecase) GetUserID(ctx context.Context, request *domain.GetUserId) (*entity.User, error) {
	user, err := u.UserRepo.FindByID(ctx, request.ID)
	if err != nil {
		u.Log.WithError(err).Error("Failed to find user by id")
		return nil, util.NewCustomError(int(util.ErrNotFoundCode), err.Error())
	}

	return user, nil
}

func (u *UserUsecase) Logout(ctx context.Context, request *domain.LogoutUserRequest) (bool, error) {
	tx := u.DB.WithContext(ctx).Begin()

	user, err := u.UserRepo.FindByID(tx.Statement.Context, request.ID)
	if err != nil {
		u.Log.WithError(err).Error("Failed to found user")
		return false, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	user.Token = ""
	if err := u.UserRepo.Update(tx.Statement.Context, user); err != nil {
		u.Log.WithError(err).Error("Failed to update token user")
		return false, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return false, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return true, nil
}

func (u *UserUsecase) Current(ctx context.Context, request *domain.CurrentUserRequest) (*domain.UserResponse, error) {
	user, err := u.UserRepo.FindByID(ctx, request.ID)
	if err != nil {
		u.Log.WithError(err).Error("Failed to found user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUsecase) Update(ctx context.Context, request *domain.UserUpdateRequest) (*domain.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()

	user, err := u.UserRepo.FindByID(tx.Statement.Context, request.ID)
	if err != nil {
		u.Log.WithError(err).Error("Failed to found user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if user.Name == request.Name || user.Email == request.Email {
		u.Log.Warn("Request name or email same before")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), "Request name or email same before")
	}

	if request.OldPassword != "" && request.NewPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
			u.Log.WithError(err).Error("Failed to compare user password with bcrype hash")
			return nil, util.NewCustomError(int(util.ErrUnauthorizedCode), err.Error())
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.NewPassword)); err == nil {
			u.Log.Warn("New password cannot be the same as the old password")
			return nil, util.NewCustomError(int(util.ErrBadRequestCode), "New password cannot be the same as the old password")
		}
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.NewPassword != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			u.Log.Warnf("Failed to generate bcrype hash : %+v", err)
			return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
		}
		user.Password = string(password)
	}

	if err := u.UserRepo.Update(tx.Statement.Context, user); err != nil {
		u.Log.WithError(err).Error("Failed to update token user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	return converter.UserToResponse(user), nil
}
