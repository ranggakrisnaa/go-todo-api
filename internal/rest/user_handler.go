package rest

import (
	"context"
	"go-todo-api/domain"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserUseCase interface {
	Create(ctx context.Context, user *domain.RegisterUserRequest) (*entity.User, error)
	Login(ctx context.Context, user *domain.LoginUserRequest) (*domain.UserResponse, error)
}

type UserHandler struct {
	Log     *logrus.Logger
	UseCase UserUseCase
}

func NewUserHandler(r *gin.Engine, usc UserUseCase, log *logrus.Logger) {
	handler := &UserHandler{
		UseCase: usc,
		Log:     log,
	}

	r.POST("v1/auth/register", handler.Register)
	r.POST("v1/auth/login", handler.Login)
}

func isRequestValid(u interface{}) (bool, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserHandler) Register(c *gin.Context) {
	var (
		user          domain.RegisterUserRequest
		errValidation error
		ok            bool
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		u.Log.WithError(err).Error("Error parsing request body")
		util.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if ok, errValidation = isRequestValid(&user); !ok {
		u.Log.WithError(errValidation).Error("Error request body validation")
		util.SendError(c, http.StatusBadRequest, errValidation.Error())
		return
	}

	response, err := u.UseCase.Create(c, &user)
	if err != nil {
		u.Log.WithError(err).Error("Error creating User")
		statusCode := util.GetStatusCode(err)
		util.SendError(c, statusCode, err.Error())
		return
	}

	util.SendSuccess(c, http.StatusCreated, "Success Created Data User", domain.UserToResponse(response))
}

func (u *UserHandler) Login(c *gin.Context) {
	var (
		user          domain.LoginUserRequest
		errValidation error
		ok            bool
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		u.Log.WithError(err).Error("Error parsing request body")
		util.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if ok, errValidation = isRequestValid(&user); !ok {
		u.Log.WithError(errValidation).Error("Error request body validation")
		util.SendError(c, http.StatusBadRequest, errValidation.Error())
		return
	}

	response, err := u.UseCase.Login(c, &user)
	if err != nil {
		u.Log.WithError(err).Error("Error login User")
		statusCode := util.GetStatusCode(err)
		util.SendError(c, statusCode, err.Error())
		return
	}

	util.SendSuccess(c, http.StatusCreated, "Success login User", response)
}
