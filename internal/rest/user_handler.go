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

type UserUseCase interface {
	Create(ctx context.Context, user *domain.RegisterUserRequest) (*domain.UserResponse, error)
	Login(ctx context.Context, user *domain.LoginUserRequest) (*domain.UserResponse, error)
	Logout(ctx context.Context, request *domain.LogoutUserRequest) (bool, error)
	Current(ctx context.Context, request *domain.CurrentUserRequest) (*domain.UserResponse, error)
	Update(ctx context.Context, request *domain.UserUpdateRequest) (*domain.UserResponse, error)
}

type UserHandler struct {
	Log     *logrus.Logger
	UseCase UserUseCase
}

func NewUserHandler(r *gin.Engine, u UserUseCase, log *logrus.Logger, authMiddleware gin.HandlerFunc) {
	handler := &UserHandler{
		UseCase: u,
		Log:     log,
	}

	r.POST("v1/users", handler.Register)
	r.POST("v1/users/_login", handler.Login)
	r.Use(authMiddleware)
	r.DELETE("v1/users", handler.Logout)
	r.GET("v1/users/_current", handler.Current)
	r.PUT("v1/users/_current", handler.Update)
}

func (u *UserHandler) Register(c *gin.Context) {
	var (
		user          domain.RegisterUserRequest
		errValidation error
		ok            bool
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		u.Log.WithError(err).Error("Error parsing request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if ok, errValidation = util.IsRequestValid(&user); !ok {
		u.Log.WithError(errValidation).Error("Error request body validation")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errValidation.Error()})
		return
	}

	response, err := u.UseCase.Create(c, &user)
	if err != nil {
		u.Log.WithError(err).Error("Error creating User")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.Response[*domain.UserResponse]{
		Status:     true,
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Data:       response,
	})
}

func (u *UserHandler) Login(c *gin.Context) {
	var (
		user          domain.LoginUserRequest
		errValidation error
		ok            bool
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		u.Log.WithError(err).Error("Error parsing request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if ok, errValidation = util.IsRequestValid(&user); !ok {
		u.Log.WithError(errValidation).Error("Error request body validation")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errValidation.Error()})
		return
	}

	response, err := u.UseCase.Login(c, &user)
	if err != nil {
		u.Log.WithError(err).Error("Error login User")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[*domain.UserResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "User login successfully",
		Data:       response,
	})
}

func (u *UserHandler) Logout(c *gin.Context) {
	auth := middleware.GetUser(c)

	request := &domain.LogoutUserRequest{
		GetUserId: domain.GetUserId{
			ID: auth.ID,
		},
	}

	_, err := u.UseCase.Logout(c, request)
	if err != nil {
		u.Log.WithError(err).Error("Error logging out user")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[*domain.UserResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "User logged out successfully",
	})
}

func (u *UserHandler) Current(c *gin.Context) {
	auth := middleware.GetUser(c)

	request := &domain.CurrentUserRequest{
		GetUserId: domain.GetUserId{
			ID: auth.ID,
		},
	}

	response, err := u.UseCase.Current(c, request)
	if err != nil {
		u.Log.WithError(err).Error("Error get data user")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[*domain.UserResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "User data retrieved successfully",
		Data:       response,
	})
}

func (u *UserHandler) Update(c *gin.Context) {
	var (
		user          domain.UserUpdateRequest
		errValidation error
		ok            bool
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		u.Log.WithError(err).Error("Error parsing request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	if ok, errValidation = util.IsRequestValid(&user); !ok {
		u.Log.WithError(errValidation).Error("Error request body validation")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errValidation.Error()})
		return
	}

	auth := middleware.GetUser(c)
	user.ID = auth.ID

	response, err := u.UseCase.Update(c, &user)
	if err != nil {
		u.Log.WithError(err).Error("Error update user")
		c.AbortWithStatusJSON(util.GetStatusCode(err), gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response[*domain.UserResponse]{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "User updated successfully",
		Data:       response,
	})
}
