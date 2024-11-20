package usecase

import (
	"context"
	"go-todo-api/domain"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	CountByEmailOrName(ctx context.Context, user *entity.User) (int64, error)
	FindByEmailOrName(ctx context.Context, user *entity.User, email, name string) error
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

func (u *UserUseCase) Create(ctx context.Context, request *domain.RegisterUserRequest) (*entity.User, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	err := u.userRepo.FindByEmailOrName(tx.Statement.Context, user, request.Email, request.Name)
	if err != nil {
		u.Log.Warnf("Error checking for existing user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), "User with email or name already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.WithError(err).Error("Failed to generate bcrype hash")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	userPayload := &entity.User{
		Name:     user.Name,
		Email:    user.Email,
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

	u.Log.Infof("User successfully created: %+v", user)
	return user, nil
}

func (u *UserUseCase) Login(ctx context.Context, request *domain.LoginUserRequest) (*domain.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	err := u.userRepo.FindByEmailOrName(tx.Statement.Context, user, request.Email, "")
	if err != nil {
		u.Log.WithError(err).Error("Failed to found user")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.Log.WithError(err).Error("Failed to compare user password with bcrype hash")
		return nil, util.NewCustomError(int(util.ErrUnauthorizedCode), err.Error())
	}

	jwtKey := os.Getenv("JWT_KEY")
	expJwt := os.Getenv("JWT_EXP")
	expJwtConv, _ := strconv.Atoi(expJwt)
	if jwtKey == "" || expJwt == "" {
		jwtKey = "private"
		expJwtConv = 72

	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UUID,
		"exp": time.Now().Add(time.Hour * time.Duration(expJwtConv)).Unix(),
	}).SignedString([]byte(jwtKey))
	if err != nil {
		u.Log.WithError(err).Error("Failed to create jwt tokne")
		return nil, util.NewCustomError(int(util.ErrUnauthorizedCode), err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Failed to commit transaction")
		return nil, util.NewCustomError(int(util.ErrInternalServerErrorCode), err.Error())
	}

	response := &domain.UserResponse{
		UUID:  user.UUID,
		Token: token,
	}

	return response, nil
}
