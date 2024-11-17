package rest

import (
	"context"
	"go-todo-api/internal/entity"
	"go-todo-api/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	Create(ctx context.Context, user *entity.User) (err error)
}

type UserHandler struct {
	UseCase UserUseCase
}

func NewUserHandler(r *gin.Engine, usc UserUseCase) {
	handler := &UserHandler{
		UseCase: usc,
	}

	r.POST("v1/auth/register", handler.Create)
}

func (h *UserHandler) Create(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		util.SendError(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := h.UseCase.Create(c, &user)
	if err != nil {
		util.SendError(c, util.GetStatusCode(err), "failed", nil)
		return
	}

	util.SendSuccess(c, http.StatusCreated, "Success Created Data User", user)
}
